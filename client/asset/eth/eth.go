// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package eth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/server/asset/eth"
	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
)

func init() {
	asset.Register(BipID, &Driver{})
}

const (
	// BipID is the BIP-0044 asset ID.
	BipID              = 60
	defaultGasFee      = 8.2e10
	defaultGasFeeLimit = 2e11
)

var (
	defaultAppDir = dcrutil.AppDataDir("dexethclient", false)
	// blockTicker is the delay between calls to check for new blocks.
	blockTicker = time.Second
	configOpts  = []*asset.ConfigOption{
		// TODO: Use this limit.
		{
			Key:         "gasfeelimit",
			DisplayName: "Gas Fee Limit",
			Description: "This is the highest network fee rate you are willing to " +
				"pay on swap transactions. If gasfeelimit is lower than a market's " +
				"maxfeerate, you will not be able to trade on that market with this " +
				"wallet.  Units: wei",
			DefaultValue: defaultGasFeeLimit,
		},
		{
			Key:          "appdir",
			DisplayName:  "DCR Dex Ethereum directory location.",
			Description:  "Location of the ethereum client data. This SHOULD NOT be a directory used by other ethereum applications. The default is recommended.",
			DefaultValue: defaultAppDir,
		},
		{
			Key:          "nodelistenaddr",
			DisplayName:  "Node Listening Address.",
			Description:  "The address the node will listen on when connecting to other ethereum nodes.",
			DefaultValue: "0.0.0.0:30303",
		},
	}
	// WalletInfo defines some general information about a Ethereum wallet.
	WalletInfo = &asset.WalletInfo{
		Name:              "Ethereum",
		Units:             "wei",
		DefaultConfigPath: defaultAppDir, // Incorrect if changed by user?
		ConfigOpts:        configOpts,
	}
	notImplementedErr   = errors.New("not implemented")
	mainnetContractAddr = common.HexToAddress("")
)

// Check that Driver implements asset.Driver.
var _ asset.Driver = (*Driver)(nil)

// Driver implements asset.Driver.
type Driver struct{}

// Setup creates the ETH exchange wallet. Start the wallet with its Run method.
func (d *Driver) Setup(cfg *asset.WalletConfig, logger dex.Logger, network dex.Network) (asset.Wallet, error) {
	return NewWallet(cfg, logger, network)
}

// DecodeCoinID creates a human-readable representation of a coin ID for Ethereum.
func (d *Driver) DecodeCoinID(coinID []byte) (string, error) {
	return eth.CoinIDToString(coinID)
}

// Info returns basic information about the wallet and asset.
func (d *Driver) Info() *asset.WalletInfo {
	return WalletInfo
}

// rawWallet is a return type from the eth client.
type rawWallet struct {
	URL      string             `json:"url"`
	Status   string             `json:"status"`
	Failure  string             `json:"failure,omitempty"`
	Accounts []accounts.Account `json:"accounts,omitempty"`
}

// ethFetcher represents a blockchain information fetcher. In practice, it is
// satisfied by rpcclient. For testing, it can be satisfied by a stub.
type ethFetcher interface {
	accounts() []*accounts.Account
	addPeer(ctx context.Context, peer string) error
	balance(ctx context.Context, acct *accounts.Account) (*big.Int, error)
	bestBlockHash(ctx context.Context) (common.Hash, error)
	block(ctx context.Context, hash common.Hash) (*types.Block, error)
	blockNumber(ctx context.Context) (uint64, error)
	connect(ctx context.Context, node *node.Node, contractAddr common.Address) error
	importAccount(pw string, privKeyB []byte) (*accounts.Account, error)
	listWallets(ctx context.Context) ([]rawWallet, error)
	lock(ctx context.Context, acct *accounts.Account) error
	nodeInfo(ctx context.Context) (*p2p.NodeInfo, error)
	pendingTransactions(ctx context.Context) ([]*types.Transaction, error)
	transactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	sendTransaction(ctx context.Context, tx map[string]string) (common.Hash, error)
	shutdown()
	syncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
	unlock(ctx context.Context, pw string, acct *accounts.Account) error
}

// Check that ExchangeWallet satisfies the asset.Wallet interface.
var _ asset.Wallet = (*ExchangeWallet)(nil)

// ExchangeWallet is a wallet backend for Ethereum. The backend is how the DEX
// client app communicates with the Ethereum blockchain and wallet. ExchangeWallet
// satisfies the dex.Wallet interface.
type ExchangeWallet struct {
	// 64-bit atomic variables first. See
	// https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	tipAtConnect int64

	// syncingStarted and initBlockNum are atomics and are used to determine
	// whether the node has started syncing.
	initBlockNum   uint64
	syncingStarted uint32

	ctx       context.Context // the asset subsystem starts with Connect(ctx)
	node      ethFetcher
	log       dex.Logger
	tipChange func(error)

	internalNode *node.Node

	tipMtx     sync.RWMutex
	currentTip *types.Block

	acctMtx sync.RWMutex
	acct    *accounts.Account
}

// Info returns basic information about the wallet and asset.
func (*ExchangeWallet) Info() *asset.WalletInfo {
	return WalletInfo
}

// NewWallet is the exported constructor by which the DEX will import the
// exchange wallet. It starts an internal light node.
func NewWallet(assetCFG *asset.WalletConfig, logger dex.Logger, network dex.Network) (*ExchangeWallet, error) {
	cfg, err := loadConfig(assetCFG.Settings, network)
	if err != nil {
		return nil, err
	}
	nodeCFG := &nodeConfig{
		net:        network,
		listenAddr: cfg.NodeListenAddr,
		appDir:     cfg.AppDir,
	}
	node, err := runNode(nodeCFG)
	if err != nil {
		return nil, err
	}
	return &ExchangeWallet{
		log:          logger,
		tipChange:    assetCFG.TipChange,
		internalNode: node,
		acct:         new(accounts.Account),
	}, nil
}

func (eth *ExchangeWallet) shutdown() {
	eth.node.shutdown()
	eth.internalNode.Close()
	eth.internalNode.Wait()
}

// Connect connects to the node RPC server. A dex.Connector.
func (eth *ExchangeWallet) Connect(ctx context.Context) (*sync.WaitGroup, error) {
	c := rpcclient{}
	if err := c.connect(ctx, eth.internalNode, mainnetContractAddr); err != nil {
		return nil, err
	}
	eth.node = &c
	eth.ctx = ctx

	// Initialize the best block.
	bestHash, err := eth.node.bestBlockHash(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting best block hash from geth: %w", err)
	}
	block, err := eth.node.block(ctx, bestHash)
	if err != nil {
		return nil, fmt.Errorf("error getting best block from geth: %w", err)
	}
	eth.tipMtx.Lock()
	eth.currentTip = block
	eth.tipMtx.Unlock()
	atomic.StoreInt64(&eth.tipAtConnect, int64(eth.currentTip.NumberU64()))
	eth.log.Infof("Connected to geth.")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		eth.monitorBlocks(ctx)
		eth.shutdown()
	}()
	return &wg, nil
}

// OwnsAddress indicates if an address belongs to the wallet.
//
// In Ethereum, an address is an account.
//
// TODO: Consider adding multiple accounts.
func (eth *ExchangeWallet) OwnsAddress(address string) (bool, error) {
	eth.acctMtx.RLock()
	defer eth.acctMtx.RUnlock()
	return eth.acct.Address.String() == address, nil
}

// Balance returns the total available funds in the account.
//
// NOTE: Ethereum balance does not return Immature or Locked values.
//
// TODO: Ethereum balances can easily go over the max value of a uint64.
// asset.Balance must be changed in a way to accomodate this.
func (eth *ExchangeWallet) Balance() (*asset.Balance, error) {
	eth.acctMtx.RLock()
	defer eth.acctMtx.RUnlock()
	bigbal, err := eth.node.balance(eth.ctx, eth.acct)
	if err != nil {
		return nil, err
	}
	bal := &asset.Balance{
		Available: bigbal.Uint64(),
		// Immature: , How to know?
		// Locked: , Not lockable?
	}
	return bal, nil
}

// MaxOrder generates information about the maximum order size and associated
// fees that the wallet can support for the given DEX configuration. The fees are an
// estimate based on current network conditions, and will be <= the fees
// associated with nfo.MaxFeeRate. For quote assets, the caller will have to
// calculate lotSize based on a rate conversion from the base asset's lot size.
func (*ExchangeWallet) MaxOrder(lotSize uint64, feeSuggestion uint64, nfo *dex.Asset) (*asset.SwapEstimate, error) {
	return nil, notImplementedErr
}

// PreSwap gets order estimates based on the available funds and the wallet
// configuration.
func (*ExchangeWallet) PreSwap(req *asset.PreSwapForm) (*asset.PreSwap, error) {
	return nil, notImplementedErr
}

// PreRedeem generates an estimate of the range of redemption fees that could
// be assessed.
func (*ExchangeWallet) PreRedeem(req *asset.PreRedeemForm) (*asset.PreRedeem, error) {
	return nil, notImplementedErr
}

// FundOrder selects coins for use in an order. The coins will be locked, and
// will not be returned in subsequent calls to FundOrder or calculated in calls
// to Available, unless they are unlocked with ReturnCoins.
// The returned []dex.Bytes contains the redeem scripts for the selected coins.
// Equal number of coins and redeemed scripts must be returned. A nil or empty
// dex.Bytes should be appended to the redeem scripts collection for coins with
// no redeem script.
func (*ExchangeWallet) FundOrder(ord *asset.Order) (asset.Coins, []dex.Bytes, error) {
	return nil, nil, notImplementedErr
}

// ReturnCoins unlocks coins. This would be necessary in the case of a
// canceled order.
func (*ExchangeWallet) ReturnCoins(unspents asset.Coins) error {
	return notImplementedErr
}

// FundingCoins gets funding coins for the coin IDs. The coins are locked. This
// method might be called to reinitialize an order from data stored externally.
// This method will only return funding coins, e.g. unspent transaction outputs.
func (*ExchangeWallet) FundingCoins(ids []dex.Bytes) (asset.Coins, error) {
	return nil, notImplementedErr
}

// Swap sends the swaps in a single transaction. The Receipts returned can be
// used to refund a failed transaction. The Input coins are manually unlocked
// because they're not auto-unlocked by the wallet and therefore inaccurately
// included as part of the locked balance despite being spent.
func (*ExchangeWallet) Swap(swaps *asset.Swaps) ([]asset.Receipt, asset.Coin, uint64, error) {
	return nil, nil, 0, notImplementedErr
}

// Redeem sends the redemption transaction, which may contain more than one
// redemption.
func (*ExchangeWallet) Redeem(form *asset.RedeemForm) ([]dex.Bytes, asset.Coin, uint64, error) {
	return nil, nil, 0, notImplementedErr
}

// SignMessage signs the message with the private key associated with the
// specified funding Coin. A slice of pubkeys required to spend the Coin and a
// signature for each pubkey are returned.
func (*ExchangeWallet) SignMessage(coin asset.Coin, msg dex.Bytes) (pubkeys, sigs []dex.Bytes, err error) {
	return nil, nil, notImplementedErr
}

// AuditContract retrieves information about a swap contract on the
// blockchain. This would be used to verify the counter-party's contract
// during a swap.
func (*ExchangeWallet) AuditContract(coinID, contract, txData dex.Bytes) (*asset.AuditInfo, error) {
	return nil, notImplementedErr
}

// LocktimeExpired returns true if the specified contract's locktime has
// expired, making it possible to issue a Refund.
func (*ExchangeWallet) LocktimeExpired(contract dex.Bytes) (bool, time.Time, error) {
	return false, time.Time{}, notImplementedErr
}

// FindRedemption watches for the input that spends the specified contract
// coin, and returns the spending input and the contract's secret key when it
// finds a spender.
//
// This method blocks until the redemption is found, an error occurs or the
// provided context is canceled.
func (*ExchangeWallet) FindRedemption(ctx context.Context, coinID dex.Bytes) (redemptionCoin, secret dex.Bytes, err error) {
	return nil, nil, notImplementedErr
}

// Refund refunds a contract. This can only be used after the time lock has
// expired.
// NOTE: The contract cannot be retrieved from the unspent coin info as the
// wallet does not store it, even though it was known when the init transaction
// was created. The client should store this information for persistence across
// sessions.
func (*ExchangeWallet) Refund(coinID, contract dex.Bytes) (dex.Bytes, error) {
	return nil, notImplementedErr
}

// Address returns an address for the exchange wallet.
func (eth *ExchangeWallet) Address() (string, error) {
	eth.acctMtx.RLock()
	defer eth.acctMtx.RUnlock()
	return eth.acct.Address.String(), nil
}

// Unlock unlocks the exchange wallet.
func (eth *ExchangeWallet) Unlock(pw string) error {
	eth.acctMtx.RLock()
	defer eth.acctMtx.RUnlock()
	return eth.node.unlock(eth.ctx, pw, eth.acct)
}

// Lock locks the exchange wallet.
func (eth *ExchangeWallet) Lock() error {
	eth.acctMtx.RLock()
	defer eth.acctMtx.RUnlock()
	return eth.node.lock(eth.ctx, eth.acct)
}

// Locked will be true if the wallet is currently locked.
func (eth *ExchangeWallet) Locked() bool {
	wallets, err := eth.node.listWallets(eth.ctx)
	if err != nil {
		eth.log.Errorf("list wallets error: %v", err)
		return false
	}
	eth.acctMtx.RLock()
	defer eth.acctMtx.RUnlock()
	var wallet rawWallet
	findWallet := func() bool {
		for _, w := range wallets {
			for _, a := range w.Accounts {
				if bytes.Equal(a.Address[:], eth.acct.Address[:]) {
					wallet = w
					return true
				}
			}
		}
		return false
	}
	if !findWallet() {
		eth.log.Error("unable to find wallet for account: %v", eth.acct)
		return false
	}
	return wallet.Status != "Unlocked"
}

// PayFee sends the dex registration fee. Transaction fees are in addition to
// the registration fee, and the fee rate is taken from the DEX configuration.
//
// NOTE: PayFee is not intended to be used with Ethereum at this time.
func (*ExchangeWallet) PayFee(address string, regFee uint64) (asset.Coin, error) {
	return nil, notImplementedErr
}

// sendToAddr sends funds from acct to addr.
func (eth *ExchangeWallet) sendToAddr(addr common.Address, amt, gasFee *big.Int) (common.Hash, error) {
	tx := map[string]string{
		"from":     fmt.Sprintf("0x%x", eth.acct.Address),
		"to":       fmt.Sprintf("0x%x", addr),
		"value":    fmt.Sprintf("0x%x", amt),
		"gasPrice": fmt.Sprintf("0x%x", gasFee),
	}
	return eth.node.sendTransaction(eth.ctx, tx)
}

// Withdraw withdraws funds to the specified address. Fees are subtracted from
// the value.
//
// TODO: Value could be larger than a uint64. Deal with it...
// TODO: Return the asset.Coin.
func (eth *ExchangeWallet) Withdraw(addr string, value uint64) (asset.Coin, error) {
	eth.acctMtx.RLock()
	defer eth.acctMtx.RUnlock()
	_, err := eth.sendToAddr(common.HexToAddress(addr),
		big.NewInt(0).SetUint64(value), big.NewInt(0).SetUint64(defaultGasFee))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateSecret checks that the secret satisfies the contract.
func (*ExchangeWallet) ValidateSecret(secret, secretHash []byte) bool {
	h := sha256.Sum256(secret)
	return bytes.Equal(h[:], secretHash)
}

// Confirmations gets the number of confirmations for the specified coin ID.
func (*ExchangeWallet) Confirmations(ctx context.Context, id dex.Bytes) (confs uint32, spent bool, err error) {
	return 0, false, notImplementedErr
}

// SyncStatus is information about the blockchain sync status.
func (eth *ExchangeWallet) SyncStatus() (bool, float32, error) {
	// node.SyncProgress will return nil both before syncing has begun and
	// after it has finished. In order to discern when syncing has begun,
	// wait for at least one change in block count, then defer to
	// node.SyncProgress for the lifetime of the eth backend.
	syncingStarted := atomic.LoadUint32(&eth.syncingStarted)
	if syncingStarted != 1 {
		ibn := atomic.LoadUint64(&eth.initBlockNum)
		bn, err := eth.node.blockNumber(eth.ctx)
		if err != nil {
			return false, 0, err
		}
		if ibn == 0 {
			atomic.StoreUint64(&eth.initBlockNum, bn)
			return false, 0, nil
		}
		if bn == ibn {
			return false, 0, nil
		}
		atomic.StoreUint32(&eth.syncingStarted, 1)
	}
	sp, err := eth.node.syncProgress(eth.ctx)
	if err != nil {
		return false, 0, err
	}
	if sp == nil {
		return true, 1, nil
	}
	ratio := float32(sp.CurrentBlock) / float32(sp.HighestBlock)
	return false, ratio, nil
}

// RefundAddress extracts and returns the refund address from a contract.
func (eth *ExchangeWallet) RefundAddress(contract dex.Bytes) (string, error) {
	return "", notImplementedErr
}

// monitorBlocks pings for new blocks and runs the tipChange callback function
// when the block changes. New blocks are also scanned for potential contract
// redeems.
func (eth *ExchangeWallet) monitorBlocks(ctx context.Context) {
	ticker := time.NewTicker(blockTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			eth.checkForNewBlocks()
		case <-ctx.Done():
			return
		}
	}
}

// checkForNewBlocks checks for new blocks. When a tip change is detected, the
// tipChange callback function is invoked and a goroutine is started to check
// if any contracts in the findRedemptionQueue are redeemed in the new blocks.
func (eth *ExchangeWallet) checkForNewBlocks() {
	ctx, cancel := context.WithTimeout(eth.ctx, 2*time.Second)
	defer cancel()
	bestHash, err := eth.node.bestBlockHash(ctx)
	if err != nil {
		go eth.tipChange(fmt.Errorf("failed to get best hash: %w", err))
		return
	}
	// This method is called frequently. Don't hold write lock
	// unless tip has changed.
	eth.tipMtx.RLock()
	currentTipHash := eth.currentTip.Hash()
	eth.tipMtx.RUnlock()
	sameTip := bytes.Equal(currentTipHash[:], bestHash[:])
	if sameTip {
		return
	}

	newTip, err := eth.node.block(ctx, bestHash)
	if err != nil {
		go eth.tipChange(fmt.Errorf("failed to get best block: %w", err))
		return
	}

	eth.tipMtx.Lock()
	defer eth.tipMtx.Unlock()

	prevTip := eth.currentTip
	eth.currentTip = newTip
	eth.log.Debugf("tip change: %d (%s) => %d (%s)", prevTip.NumberU64(),
		prevTip.Hash(), newTip.NumberU64(), newTip.Hash())
	go eth.tipChange(nil)
}
