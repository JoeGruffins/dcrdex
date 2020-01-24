// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

// Package rpcserver provides a JSON RPC to communicate with the client core.
package rpcserver

import (
	"context"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex/msgjson"
	"github.com/decred/dcrd/certgen"
	"github.com/decred/slog"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	// rpcTimeoutSeconds is the number of seconds a connection to the
	// RPC server is allowed to stay open without authenticating before it
	// is closed.
	rpcTimeoutSeconds = 10

	// RPC version
	rpcSemverMajor = 0
	rpcSemverMinor = 0
	rpcSemverPatch = 0
)

var log slog.Logger

// ClientCore is satisfied by core.Core.
type ClientCore interface {
	ListMarkets() []*core.MarketInfo
	Register(*core.Registration) error
	Login(dex, pw string) error
	Sync(dex string, base, quote uint32) (chan *core.BookUpdate, error)
	Book(dex string, base, quote uint32) *core.OrderBook
	Unsync(dex string, base, quote uint32)
	Balance(uint32) (uint64, error)
}

// marketSyncer is used to synchronize market subscriptions. The marketSyncer
// manages a map of clients who are subscribed to the market, and distributes
// order book updates when received.
type marketSyncer struct {
	mtx     sync.Mutex
	core    ClientCore
	dex     string
	base    uint32
	quote   uint32
	clients map[int32]*wsClient
}

// newMarketSyncer is the constructor for a marketSyncer.
func newMarketSyncer(ctx context.Context, wg *sync.WaitGroup, core ClientCore, dex string, base, quote uint32) (*marketSyncer, error) {
	m := &marketSyncer{
		core:    core,
		dex:     dex,
		base:    base,
		quote:   quote,
		clients: make(map[int32]*wsClient),
	}

	// Get an updates channel, and begin syncing the book.
	updates, err := core.Sync(dex, base, quote)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Debugf("monitoring market %d-%d @ %s", base, quote, dex)
	out:
		for {
			select {
			case update := <-updates:
				// Distribute the book to the subscribed clients.
				log.Tracef("order book update received for " + update.Market)
			case <-ctx.Done():
				break out
			}
		}
	}()
	return m, nil
}

// add adds a client to the client map, and returns a fresh orderbook.
func (m *marketSyncer) add(cl *wsClient) *core.OrderBook {
	m.mtx.Lock()
	m.clients[cl.cid] = cl
	m.mtx.Unlock()
	return m.core.Book(m.dex, m.base, m.quote)
}

// remove removes a client from the client map. If this is the last client,
// the market will be "unsynced".
func (m *marketSyncer) remove(cl *wsClient) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.clients, cl.cid)
	if len(m.clients) == 0 {
		m.core.Unsync(m.dex, m.base, m.quote)
	}
}

// RPCServer is a single-client http and websocket server enabling a JSON
// interface to the DEX client.
type RPCServer struct {
	ctx      context.Context
	core     ClientCore
	listener net.Listener
	srv      *http.Server
	authsha  [32]byte
	mtx      sync.RWMutex
	syncers  map[string]*marketSyncer
	clients  map[int32]*wsClient
	wg       sync.WaitGroup
}

// genCertPair generates a key/cert pair to the paths provided.
func genCertPair(certFile, keyFile string) error {
	log.Infof("Generating TLS certificates...")

	org := "dcrdex autogenerated cert"
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, err := certgen.NewTLSCertPair(elliptic.P521(), org,
		validUntil, nil)
	if err != nil {
		return err
	}

	// Write cert and key files.
	if err = ioutil.WriteFile(certFile, cert, 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(keyFile, key, 0600); err != nil {
		os.Remove(certFile)
		return err
	}

	log.Infof("Done generating TLS certificates")
	return nil
}

// handler is the type for functions that handle RPC requests.
type handler func(*RPCServer, *msgjson.Message) *msgjson.ResponsePayload

// routes maps routes to their handler.
var routes = map[string]handler{
	"help":    handleHelp,
	"version": handleVersion,
}

// writeJSON marshals the provided interface and writes the bytes to the
// ResponseWriter. The response code is assumed to be StatusOK.
func writeJSON(w http.ResponseWriter, thing interface{}) {
	writeJSONWithStatus(w, thing, http.StatusOK)
}

// writeJSONWitStatus marshals the provided interface and writes the bytes to the
// ResponseWriter with the specified response code.
func writeJSONWithStatus(w http.ResponseWriter, thing interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(thing); err != nil {
		log.Errorf("JSON encode error: %v", err)
	}
}

// handleJSON handles all https json requests.
func (s *RPCServer) handleJSON(w http.ResponseWriter, r *http.Request) {
	// All http routes are available over websocket too, so do not support
	// persistent http connections. Inform the user and close the connection
	// when response handling is completed.
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "application/json")
	r.Close = true

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}
	req := new(msgjson.Message)
	err = json.Unmarshal(body, req)
	if err != nil {
		http.Error(w, "JSON decode error", http.StatusUnprocessableEntity)
		return
	}
	if req.Type != msgjson.Request {
		http.Error(w, "Responses not accepted", http.StatusMethodNotAllowed)
		return
	}
	s.parseHTTPRequest(w, req)
}

// Config holds variables neede to create a new RPC Server.
type Config struct {
	Core                        ClientCore
	Addr, User, Pass, Cert, Key string
}

// SetLogger sets the logger for the RPCServer package.
func SetLogger(logger slog.Logger) {
	log = logger
}

// New is the constructor for an RPCServer.
func New(cfg *Config) (*RPCServer, error) {

	// Find or create the key pair.
	keyExists := fileExists(cfg.Key)
	certExists := fileExists(cfg.Cert)
	if certExists == !keyExists {
		return nil, fmt.Errorf("missing cert pair file")
	}
	if !keyExists && !certExists {
		err := genCertPair(cfg.Cert, cfg.Key)
		if err != nil {
			return nil, err
		}
	}
	keypair, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return nil, err
	}

	// Prepare the TLS configuration.
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{keypair},
		MinVersion:   tls.VersionTLS12,
	}

	// Create listener.
	listener, err := tls.Listen("tcp", cfg.Addr, &tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("Can't listen on %s. web server quitting: %v", cfg.Addr, err)
	}

	// Create an HTTP router.
	mux := chi.NewRouter()
	httpServer := &http.Server{
		Handler:      mux,
		ReadTimeout:  rpcTimeoutSeconds * time.Second, // slow requests should not hold connections opened
		WriteTimeout: rpcTimeoutSeconds * time.Second, // hung responses must die
	}

	// Make the server.
	s := &RPCServer{
		core:     cfg.Core,
		listener: listener,
		srv:      httpServer,
		syncers:  make(map[string]*marketSyncer),
		clients:  make(map[int32]*wsClient),
	}

	// Create authsha to verify requests against.
	if cfg.User != "" && cfg.Pass != "" {
		login := cfg.User + ":" + cfg.Pass
		auth := "Basic " +
			base64.StdEncoding.EncodeToString([]byte(login))
		s.authsha = sha256.Sum256([]byte(auth))
	}

	// Middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RealIP)
	mux.Use(s.authMiddleware)

	// Websocket endpoint
	mux.Get("/ws", s.handleWS)

	// https endpoint
	mux.Post("/", s.handleJSON)

	return s, nil
}

// Run starts the web server. Satisfies the dex.Runner interface. ctx
// passed to newMarketSyncer when making new market syncers.
func (s *RPCServer) Run(ctx context.Context) {
	s.ctx = ctx
	// Close the listener on context cancellation.
	go func() {
		<-ctx.Done()

		if err := s.srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners:
			err = fmt.Errorf("HTTP server Shutdown: %v", err)
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	log.Infof("RPC server listening on %s", s.listener.Addr())
	if err := s.srv.Serve(s.listener); err != http.ErrServerClosed {
		log.Warnf("unexpected (http.Server).Serve error: %v", err)
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, cl := range s.clients {
		cl.Disconnect()
	}

	// Wait for market syncers to finish.
	s.wg.Wait()
	log.Infof("RPC server off")
}

// watchMarket watches the specified market for changes in orders. A fresh order
// book and a quit function are returned on success. The quit function should be
// called to unsubscribe the client from the market.
func (s *RPCServer) watchMarket(cl *wsClient, dex string, base, quote uint32) (book *core.OrderBook, quit func(), err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	mktID := marketID(base, quote)
	syncer, found := s.syncers[mktID]
	if !found {
		syncer, err = newMarketSyncer(s.ctx, &s.wg, s.core, dex, base, quote)
		if err != nil {
			return
		}
		s.syncers[mktID] = syncer
	}
	book = syncer.add(cl)
	return book, func() {
		syncer.remove(cl)
	}, nil
}

// handleRequest sends the request to the correct handler function if able.
func (s *RPCServer) handleRequest(req *msgjson.Message) *msgjson.ResponsePayload {
	payload := new(msgjson.ResponsePayload)
	if req.Route == "" {
		log.Debugf("received empty request")
		payload.Error = msgjson.NewError(msgjson.RPCUnknownRoute, "no route was supplied")
		return payload
	}

	// Find the correct handler for this route.
	h, ok := routes[req.Route]

	if !ok {
		log.Debugf("could not find route: %v", req.Route)
		payload.Error = msgjson.NewError(msgjson.RPCUnknownRoute, "unknown route")
		return payload
	}

	return h(s, req)
}

// parseHTTPRequest parses the msgjson message in the request body and creates a
// response message.
func (s *RPCServer) parseHTTPRequest(w http.ResponseWriter, req *msgjson.Message) {
	payload := s.handleRequest(req)
	resp, err := msgjson.NewResponse(req.ID, payload.Result, payload.Error)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
	writeJSON(w, resp)
}

// authMiddleware checks incoming requests for authentication.
func (s *RPCServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header["Authorization"]
		if len(auth) == 0 || s.authsha != sha256.Sum256([]byte(auth[0])) {
			log.Warnf("authentication failure from ip: %s with auth: %s", r.RemoteAddr, auth)
			w.Header().Add("WWW-Authenticate", `Basic realm="dex RPC"`)
			http.Error(w, "401 Unauthorized.", http.StatusUnauthorized)
			return
		}
		log.Tracef("authenticated user with ip: %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// filesExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Create a unique ID for a market.
func marketID(base, quote uint32) string {
	return strconv.Itoa(int(base)) + "_" + strconv.Itoa(int(quote))
}
