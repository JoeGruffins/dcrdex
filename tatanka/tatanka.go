// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package tatanka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/dex/feerates"
	"decred.org/dcrdex/dex/fiatrates"
	"decred.org/dcrdex/dex/msgjson"
	"decred.org/dcrdex/server/comms"
	"decred.org/dcrdex/tatanka/chain"
	"decred.org/dcrdex/tatanka/db"
	"decred.org/dcrdex/tatanka/mj"
	"decred.org/dcrdex/tatanka/tanka"
	"decred.org/dcrdex/tatanka/tcp"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	version = 0

	// tatankaUniqueID is the unique ID used to register a Tatanka as a fiat
	// rate listener.
	tatankaUniqueID = "Tatanka"
)

// remoteTatanka is a remote tatanka node. A remote tatanka node can either
// be outgoing (whitelist loop) or incoming via handleInboundTatankaConnect.
type remoteTatanka struct {
	*peer
	cfg atomic.Value // mj.TatankaConfig
}

type Topic struct {
	subjects    map[tanka.Subject]map[tanka.PeerID]struct{}
	subscribers map[tanka.PeerID]struct{}
}

func (topic *Topic) unsubUser(peerID tanka.PeerID) {
	if _, found := topic.subscribers[peerID]; !found {
		return
	}
	for subID, subs := range topic.subjects {
		delete(subs, peerID)
		if len(subs) == 0 {
			delete(topic.subjects, subID)
		}
	}
	delete(topic.subscribers, peerID)
}

// BootNode represents a configured boot node. Tatanka is whitelist only, and
// node operators are responsible for keeping their whitelist up to date.
type BootNode struct {
	// Protocol is one of ("ws", "wss"), though other tatanka comms protocols
	// may be implemented later. Or we may end up using e.g. go-libp2p.
	Protocol string    `json:"protocol"`
	PeerID   dex.Bytes `json:"peer_id"`
	// Config can take different forms depending on the comms protocol, but is
	// probably a tcp.RemoteNodeConfig.
	Config json.RawMessage `json:"config"`
}

// Draft note: Why do we need a parsedBootNode? Why not just make BootNode
// have a tanka.PeerID instead of a dex.Bytes? Also, instead of having a
// TatankaCredentials in client/conn, can we use BootNode as well?

// parsedBootNode is the unexported version of BootNode, but with a PeerID
// instead of []byte.
type parsedBootNode struct {
	peerID   tanka.PeerID
	cfg      json.RawMessage
	protocol string
}

// Tatanka is a server node on Tatanka Mesh. Tatanka implements two APIs, one
// for fellow tatanka nodes, and one for clients. The primary roles of a
// tatanka node are
//  1. Maintain reputation information about client nodes.
//  2. Distribute broadcasts and relay tankagrams.
//  3. Provide some basic oracle services.
type Tatanka struct {
	net             dex.Network
	log             dex.Logger
	tcpSrv          *tcp.Server
	dataDir         string
	ctx             context.Context
	wg              *sync.WaitGroup
	whitelist       map[tanka.PeerID]*parsedBootNode
	db              *db.DB
	nets            atomic.Value // []uint32
	specialHandlers map[string]func(tanka.Sender, *msgjson.Message) *msgjson.Error
	clientHandlers  map[string]interface{} // clientRequestHandler | clientNotificationHandler
	tatankaHandlers map[string]interface{} // tatankaRequestHandler | tatankaNotificationHandler
	routes          []string
	httpReqHandlers map[string]comms.HTTPHandler
	httpRoutes      []string
	maxClients      int
	// bondTier  atomic.Uint64

	priv *secp256k1.PrivateKey
	id   tanka.PeerID

	chainMtx sync.RWMutex
	chains   map[uint32]chain.Chain

	relayMtx     sync.Mutex
	recentRelays map[[32]byte]time.Time

	clientMtx sync.RWMutex
	clients   map[tanka.PeerID]*client
	topics    map[tanka.Topic]*Topic

	clientJobs    chan *clientJob
	remoteClients map[tanka.PeerID]map[tanka.PeerID]struct{}

	tatankasMtx sync.RWMutex
	tatankas    map[tanka.PeerID]*remoteTatanka

	fiatRateOracle *fiatrates.Oracle
	fiatRateChan   chan map[string]*fiatrates.FiatRateInfo

	feeRatesOracle       *feerates.Oracle
	feeRateEstimatesChan chan map[uint32]*feerates.Estimate
}

// Config is the configuration of the Tatanka.
type Config struct {
	Net        dex.Network
	DataDir    string
	Logger     dex.Logger
	RPC        comms.RPCConfig
	ConfigPath string
	MaxClients int

	feeRatesOracleCfg feerates.Config

	WhiteList []BootNode

	FiatOracleConfig fiatrates.Config
}

func New(cfg *Config) (*Tatanka, error) {
	chainCfg, err := loadConfig(cfg.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config %w", err)
	}

	chains := make(map[uint32]chain.Chain)
	nets := make([]uint32, 0, len(chainCfg.Chains))
	for _, c := range chainCfg.Chains {
		chainID, found := dex.BipSymbolID(c.Symbol)
		if !found {
			return nil, fmt.Errorf("no chain ID found for symbol %s", c.Symbol)
		}
		chains[chainID], err = chain.New(chainID, c.Config, cfg.Logger.SubLogger(c.Symbol), cfg.Net)
		if err != nil {
			return nil, fmt.Errorf("error creating chain backend: %w", err)
		}
		nets = append(nets, chainID)
	}

	db, err := db.New(filepath.Join(cfg.DataDir, "db"), cfg.Logger.SubLogger("DB"))
	if err != nil {
		return nil, fmt.Errorf("db.New error: %w", err)
	}

	keyPath := filepath.Join(cfg.DataDir, "priv.key")
	keyB, err := os.ReadFile(keyPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error reading key file")
		}
		cfg.Logger.Infof("No key file found. Generating new identity.")
		priv, err := secp256k1.GeneratePrivateKey()
		if err != nil {
			return nil, fmt.Errorf("GeneratePrivateKey error: %w", err)
		}
		keyB = priv.Serialize()
		if err = os.WriteFile(keyPath, keyB, 0600); err != nil {
			return nil, fmt.Errorf("error writing newly-generated key to %q: %v", keyPath, err)
		}
	}
	priv := secp256k1.PrivKeyFromBytes(keyB)
	var peerID tanka.PeerID
	copy(peerID[:], priv.PubKey().SerializeCompressed())

	whitelist := make(map[tanka.PeerID]*parsedBootNode, len(cfg.WhiteList))
	for _, n := range cfg.WhiteList {
		if len(n.PeerID) != tanka.PeerIDLength {
			return nil, fmt.Errorf("invalid peer ID length %d for %s boot node with configuration %q", len(n.PeerID), n.Protocol, n.Config)
		}
		var peerID tanka.PeerID
		copy(peerID[:], n.PeerID)
		whitelist[peerID] = &parsedBootNode{
			peerID:   peerID,
			cfg:      n.Config,
			protocol: n.Protocol,
		}
	}

	if cfg.MaxClients < 1 {
		return nil, fmt.Errorf("max clients must be greater than 0")
	}

	t := &Tatanka{
		net:             cfg.Net,
		dataDir:         cfg.DataDir,
		log:             cfg.Logger,
		whitelist:       whitelist,
		db:              db,
		priv:            priv,
		id:              peerID,
		chains:          chains,
		tatankas:        make(map[tanka.PeerID]*remoteTatanka),
		clients:         make(map[tanka.PeerID]*client),
		remoteClients:   make(map[tanka.PeerID]map[tanka.PeerID]struct{}),
		topics:          make(map[tanka.Topic]*Topic),
		recentRelays:    make(map[[32]byte]time.Time),
		clientJobs:      make(chan *clientJob, 128),
		clientHandlers:  make(map[string]interface{}),
		tatankaHandlers: make(map[string]interface{}),
		httpReqHandlers: make(map[string]comms.HTTPHandler),
		maxClients:      cfg.MaxClients,
	}

	if !cfg.FiatOracleConfig.AllFiatSourceDisabled() {
		var tickers string
		upperCaser := cases.Upper(language.AmericanEnglish)
		for _, c := range chainCfg.Chains {
			tickers += upperCaser.String(c.Symbol) + ","
		}
		tickers = strings.Trim(tickers, ",")

		t.fiatRateOracle, err = fiatrates.NewFiatOracle(cfg.FiatOracleConfig, tickers, t.log)
		if err != nil {
			return nil, fmt.Errorf("error initializing fiat oracle: %w", err)
		}

		// Register tatanka as a listener
		t.fiatRateChan = make(chan map[string]*fiatrates.FiatRateInfo)
		t.fiatRateOracle.AddFiatRateListener(tatankaUniqueID, t.fiatRateChan)
	}

	t.nets.Store(nets)
	t.prepareHandlers()
	t.tcpSrv, err = tcp.NewServer(&cfg.RPC, &tcpCore{t}, cfg.Logger.SubLogger("TCP"))
	if err != nil {
		return nil, fmt.Errorf("error starting TPC server:: %v", err)
	}

	if t.net == dex.Mainnet || t.net == dex.Testnet {
		t.feeRateEstimatesChan = make(chan map[uint32]*feerates.Estimate)
		t.feeRatesOracle, err = feerates.NewOracle(t.net, cfg.feeRatesOracleCfg, nets, t.feeRateEstimatesChan)
		if err != nil {
			return nil, fmt.Errorf("feerates.NewOracle error: %w", err)
		}
	}

	return t, nil
}

func (t *Tatanka) hasFeeRatesOracle() bool {
	return t.feeRatesOracle != nil
}

func (t *Tatanka) prepareHandlers() {
	t.specialHandlers = map[string]func(tanka.Sender, *msgjson.Message) *msgjson.Error{
		// tatanka messages
		mj.RouteTatankaConnect: t.handleInboundTatankaConnect,
		// client messages
		mj.RouteConnect:  t.handleClientConnect,
		mj.RoutePostBond: t.handlePostBond,
	}
	// Tatanka routes
	registerTatankaHandler := func(route string, handler interface{}) {
		if _, is := handler.(tatankaRequestHandler); !is {
			if _, is := handler.(tatankaNotificationHandler); !is {
				panic("unknown handler type for " + route)
			}
		}
		t.tatankaHandlers[route] = handler
	}
	for route, handler := range map[string]interface{}{
		mj.RouteTatankaConfig:    t.handleTatankaConfig,
		mj.RouteRelayBroadcast:   t.handleRelayBroadcast,
		mj.RouteNewClient:        t.handleNewRemoteClientNotification,
		mj.RouteClientDisconnect: t.handleRemoteClientDisconnect,
		mj.RouteRelayTankagram:   t.handleRelayedTankagram,
		mj.RoutePathInquiry:      t.handlePathInquiry,
		mj.RouteShareScore:       t.handleShareScore,
	} {
		registerTatankaHandler(route, handler)
	}
	// Client routes
	registerClientHandler := func(route string, handler interface{}) {
		if _, is := handler.(clientRequestHandler); !is {
			if _, is := handler.(clientNotificationHandler); !is {
				panic("unknown handler type for " + route)
			}
		}
		t.clientHandlers[route] = handler
	}
	for route, handler := range map[string]interface{}{
		mj.RouteSubscribe:           t.handleSubscription,
		mj.RouteUpdateSubscriptions: t.handleUpdateSubscriptions,
		// mj.RouteUnsubscribe: t.handleUnsubscribe,
		mj.RouteBroadcast: t.handleBroadcast,
		mj.RouteTankagram: t.handleTankagram,
		mj.RouteSetScore:  t.handleSetScore,
	} {
		registerClientHandler(route, handler)
	}
	for route := range t.tatankaHandlers {
		t.routes = append(t.routes, route)
	}
	for route := range t.specialHandlers {
		t.routes = append(t.routes, route)
	}
	for route := range t.clientHandlers {
		t.routes = append(t.routes, route)
	}

	// HTTP Requests
	registerHTTPRequestHandler := func(route string, handler comms.HTTPHandler) {
		t.httpReqHandlers[route] = handler
		t.httpRoutes = append(t.httpRoutes, route)
	}
	for route, handler := range map[string]comms.HTTPHandler{
		mj.RouteNodeInfo: t.handleNodeInfo,
	} {
		registerHTTPRequestHandler(route, handler)
	}
}

func (t *Tatanka) assets() []uint32 {
	return t.nets.Load().([]uint32)
}

func (t *Tatanka) fiatOracleEnabled() bool {
	return t.fiatRateOracle != nil
}

func (t *Tatanka) tatankaNodes() []*remoteTatanka {
	t.tatankasMtx.RLock()
	defer t.tatankasMtx.RUnlock()
	nodes := make([]*remoteTatanka, 0, len(t.tatankas))
	for _, n := range t.tatankas {
		nodes = append(nodes, n)
	}
	return nodes
}

func (t *Tatanka) tatankaNode(peerID tanka.PeerID) *remoteTatanka {
	t.tatankasMtx.RLock()
	defer t.tatankasMtx.RUnlock()
	return t.tatankas[peerID]
}

func (t *Tatanka) clientNode(peerID tanka.PeerID) *client {
	t.clientMtx.RLock()
	defer t.clientMtx.RUnlock()
	return t.clients[peerID]
}

func (t *Tatanka) Connect(ctx context.Context) (_ *sync.WaitGroup, err error) {
	t.ctx = ctx
	var wg sync.WaitGroup
	t.wg = &wg

	t.log.Infof("Starting Tatanka node with peer ID %s", t.id)

	// Start WebSocket server
	cm := dex.NewConnectionMaster(t.tcpSrv)
	if err := cm.ConnectOnce(ctx); err != nil {
		return nil, fmt.Errorf("error connecting TCP server: %v", err)
	}

	wg.Add(1)
	go func() {
		cm.Wait()
		wg.Done()
	}()

	// Start a ticker to clean up the recent relays map.
	wg.Add(1)
	go func() {
		defer wg.Done()
		tick := time.NewTicker(tanka.EpochLength * 4) // 1 minute
		for {
			select {
			case <-tick.C:
				t.relayMtx.Lock()
				for bid, stamp := range t.recentRelays {
					if time.Since(stamp) > time.Minute {
						delete(t.recentRelays, bid)
					}
				}
				t.relayMtx.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.runRemoteClientsLoop(ctx)
	}()

	var success bool
	defer func() {
		if !success {
			cm.Disconnect()
			cm.Wait()
		}
	}()

	t.chainMtx.RLock()
	for assetID, c := range t.chains {
		feeRater, is := c.(chain.FeeRater)
		if !is {
			continue
		}
		wg.Add(1)
		go func(assetID uint32, feeRater chain.FeeRater) {
			defer wg.Done()
			t.monitorChainFees(ctx, assetID, feeRater)
		}(assetID, feeRater)
	}
	t.chainMtx.RUnlock()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.runWhitelistLoop(ctx)
	}()

	if t.fiatOracleEnabled() {
		wg.Add(2)
		go func() {
			defer wg.Done()
			t.fiatRateOracle.Run(t.ctx)
		}()

		go func() {
			defer wg.Done()
			t.broadcastRates()
		}()
	}

	if t.hasFeeRatesOracle() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.feeRatesOracle.Run(t.ctx, t.log)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			t.broadcastFeeRateEstimates()
		}()
	}

	success = true
	return &wg, nil
}

// runWhitelistLoop attempts to connect to the whitelist, and then periodically
// tries again.
func (t *Tatanka) runWhitelistLoop(ctx context.Context) {
	connectWhitelist := func() {
		for proto, n := range t.whitelist {
			t.tatankasMtx.RLock()
			_, exists := t.tatankas[n.peerID]
			t.tatankasMtx.RUnlock()
			if exists {
				continue
			}

			p, err := t.db.Peer(n.peerID)
			if err != nil {
				t.log.Errorf("error getting peer info for boot node at %q (proto %q): %v", string(n.cfg), proto, err)
				continue
			}

			bondTier := p.BondTier()
			// TODO: Check Tatanka Node reputation too
			// if calcTier(rep, bondTier) <= 0 {
			// 	t.log.Errorf("not attempting to contact banned boot node at %q (proto %q)", string(n.cfg), proto)
			// }

			handleDisconnect := func() {
				// TODO: schedule a reconnect?
				t.tatankasMtx.Lock()
				delete(t.tatankas, p.ID)
				t.tatankasMtx.Unlock()
			}

			handleMessage := func(cl tanka.Sender, msg *msgjson.Message) {
				t.handleTatankaMessage(cl, msg)
			}

			var cl tanka.Sender
			switch n.protocol {
			case "ws", "wss":
				cl, err = t.tcpSrv.ConnectBootNode(ctx, n.cfg, handleMessage, handleDisconnect)
			default:
				t.log.Errorf("unknown boot node network protocol: %s", proto)
				continue
			}
			if err != nil {
				t.log.Errorf("error connecting boot node with proto = %s, config = %s", proto, string(n.cfg))
				continue
			}

			t.log.Infof("Connected to boot node with peer ID %s, config %s", n.peerID, string(n.cfg))

			cl.SetPeerID(p.ID)
			pp := &peer{Peer: p, Sender: cl, rrs: make(map[tanka.PeerID]*tanka.Reputation)}
			tt := &remoteTatanka{peer: pp}
			t.tatankasMtx.Lock()
			t.tatankas[p.ID] = tt
			t.tatankasMtx.Unlock()

			cfgMsg := mj.MustRequest(mj.RouteTatankaConnect, t.generateConfig(bondTier))
			if err := t.request(cl, cfgMsg, func(msg *msgjson.Message) {
				// Nothing to do. The only non-error result is payload = true.
			}); err != nil {
				t.log.Errorf("Error sending connect message: %w", err)
				cl.Disconnect()
			}
		}
	}

	for {
		connectWhitelist()

		select {
		case <-time.After(time.Minute * 5):
		case <-ctx.Done():
			return
		}
	}
}

// monitorChainFees monitors chains for new fee rates, and will distribute them
// as part of the not-yet-implemented oracle services the mesh provides.
func (t *Tatanka) monitorChainFees(ctx context.Context, assetID uint32, c chain.FeeRater) {
	feeC := c.FeeChannel()
	for {
		select {
		case feeRate := <-feeC:
			// TODO: Distribute the fee rate to other Tatanka nodes, then to
			// clients. Should fee rates be averaged across tatankas somehow?
			fmt.Printf("new fee rate from %s: %d\n", dex.BipIDSymbol(assetID), feeRate)
		case <-ctx.Done():
			return
		}
	}
}

// sendResult sends the response to a request and logs errors.
func (t *Tatanka) sendResult(cl tanka.Sender, msgID uint64, result interface{}) {
	resp := mj.MustResponse(msgID, result, nil)
	if err := t.send(cl, resp); err != nil {
		peerID := cl.PeerID()
		t.log.Errorf("error sending result to %q: %v", dex.Bytes(peerID[:]), err)
	}
}

// batchSend must be called with the clientMtx >= RLocked.
func (t *Tatanka) batchSend(peers map[tanka.PeerID]struct{}, msg *msgjson.Message) {
	mj.SignMessage(t.priv, msg)
	msgB, err := json.Marshal(msg)
	if err != nil {
		t.log.Errorf("error marshaling batch send message: %v", err)
		return
	}
	disconnects := make(map[tanka.PeerID]struct{})
	t.clientMtx.RLock()
	for peerID := range peers {
		if c, found := t.clients[peerID]; found {
			if err := c.SendRaw(msgB); err != nil {
				t.log.Tracef("Disconnecting client %s after SendRaw error: %v", peerID, err)
				disconnects[peerID] = struct{}{}
			}
		} else {
			t.log.Error("found a subscriber ID without a client")
		}
	}
	t.clientMtx.RUnlock()
	if len(disconnects) > 0 {
		for peerID := range disconnects {
			t.clientDisconnected(peerID)
		}
	}
}

// send signs and sends the message, returning any errors.
func (t *Tatanka) send(s tanka.Sender, msg *msgjson.Message) error {
	mj.SignMessage(t.priv, msg)
	err := s.Send(msg)
	if err != nil {
		t.clientDisconnected(s.PeerID())
	}
	return err
}

// request signs and sends the request, returning any errors.
func (t *Tatanka) request(s tanka.Sender, msg *msgjson.Message, respHandler func(*msgjson.Message)) error {
	mj.SignMessage(t.priv, msg)
	err := s.Request(msg, respHandler)
	if err != nil {
		t.clientDisconnected(s.PeerID())
	}
	return err
}

func (t *Tatanka) generateConfig(bondTier uint64) *mj.TatankaConfig {
	return &mj.TatankaConfig{
		ID:       t.id,
		Version:  version,
		Chains:   t.assets(),
		BondTier: bondTier,
	}
}

func calcTier(r *tanka.Reputation, bondTier uint64) int64 {
	return int64(bondTier) + r.Score/tanka.TierIncrement
}

// ChainConfig is how the chain configuration is specified in the Tatanka
// configuration file.
type ChainConfig struct {
	Symbol string          `json:"symbol"`
	Config json.RawMessage `json:"config"`
}

// ConfigFile represents the JSON Tatanka configuration file.
type ConfigFile struct {
	Chains []ChainConfig `json:"chains"`
}

func loadConfig(configPath string) (*ConfigFile, error) {
	var cfg ConfigFile
	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("OpenFile error: %w", err)
	}
	return &cfg, json.Unmarshal(b, &cfg)
}

// tcpCore implements tcp.TankaCore.
type tcpCore struct {
	*Tatanka
}

func (t *tcpCore) Routes() []string {
	return t.routes
}

func (t *tcpCore) HandleMessage(cl tanka.Sender, msg *msgjson.Message) *msgjson.Error {
	if t.log.Level() == dex.LevelTrace {
		t.log.Tracef("Tatanka node handling message. route = %s, payload = %s", msg.Route, string(msg.Payload))
	}

	if _, found := t.clientHandlers[msg.Route]; found {
		return t.handleClientMessage(cl, msg)
	}
	if _, found := t.tatankaHandlers[msg.Route]; found {
		return t.handleTatankaMessage(cl, msg)
	}
	if handle, found := t.specialHandlers[msg.Route]; found {
		return handle(cl, msg)
	}

	fmt.Println("no route found", msg.Route)

	return msgjson.NewError(mj.ErrBadRequest, "route %q not known", msg.Route)
}

func (t *tcpCore) HTTPRoutes() []string {
	return t.httpRoutes
}

func (t *tcpCore) HandleRequest(route string, thing any) (any, error) {
	handler, found := t.httpReqHandlers[route]
	if !found {
		return nil, fmt.Errorf("http route %v not found", route)
	}

	return handler(thing)
}

// clientDisconnected handle a client disconnect, removing the client from the
// clients map and unsubscribing from all topics.
func (t *Tatanka) clientDisconnected(peerID tanka.PeerID) {
	unsubs := make(map[tanka.Topic]*Topic)

	t.clientMtx.Lock()
	delete(t.clients, peerID)
	for n, topic := range t.topics {
		if _, found := topic.subscribers[peerID]; found {
			unsubs[n] = topic
			delete(topic.subscribers, peerID)
			for _, subs := range topic.subjects {
				delete(subs, peerID)
			}
		}
	}
	t.clientMtx.Unlock()

	if len(unsubs) == 0 {
		return
	}

	stamp := time.Now()
	for n, topic := range unsubs {
		note := mj.MustNotification(mj.RouteBroadcast, &mj.Broadcast{
			Topic:       n,
			PeerID:      peerID,
			MessageType: mj.MessageTypeUnsubTopic,
			Stamp:       stamp,
		})
		t.batchSend(topic.subscribers, note)
	}

	note := mj.MustNotification(mj.RouteClientDisconnect, &mj.Disconnect{ID: peerID})
	mj.SignMessage(t.priv, note)
	for _, tt := range t.tatankaNodes() {
		tt.Send(note)
	}
}

// broadcastRates sends market rates to all fiat rate subscribers once new rates
// are received from the fiat oracle.
func (t *Tatanka) broadcastRates() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case rates, ok := <-t.fiatRateChan:
			if !ok {
				t.log.Debug("Tatanka stopped listening for fiat rates.")
				return
			}

			t.clientMtx.RLock()
			topic := t.topics[mj.TopicFiatRate]
			t.clientMtx.RUnlock()

			// TODO: Is this a race condition on topic.subscribers?
			if topic != nil && len(topic.subscribers) > 0 {
				t.batchSend(topic.subscribers, mj.MustNotification(mj.RouteRates, &mj.RateMessage{
					Topic: mj.TopicFiatRate,
					Rates: rates,
				}))
			}
		}
	}
}

// broadcastFeeRateEstimates broadcasts fee rate estimates to all subscribed clients once
// we receive data from t.feeRateEstimatesChan.
func (t *Tatanka) broadcastFeeRateEstimates() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case feeRateEstimates, ok := <-t.feeRateEstimatesChan:
			if !ok {
				t.log.Debug("Tatanka stopped listening for fee rate estimates.")
				return
			}

			t.clientMtx.RLock()
			topic := t.topics[mj.TopicFeeRateEstimate]
			t.clientMtx.RUnlock()
			if topic == nil || len(topic.subscribers) == 0 {
				continue
			}

			note := mj.MustNotification(mj.RouteFeeRateEstimate, &mj.FeeRateEstimateMessage{
				Topic:            mj.TopicFeeRateEstimate,
				FeeRateEstimates: feeRateEstimates,
			})
			t.batchSend(topic.subscribers, note)
		}
	}
}
