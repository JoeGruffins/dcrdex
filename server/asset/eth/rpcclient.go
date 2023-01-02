// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

//go:build lgpl

package eth

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"decred.org/dcrdex/dex"
	dexeth "decred.org/dcrdex/dex/networks/eth"
	swapv0 "decred.org/dcrdex/dex/networks/eth/contracts/v0"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Check that rpcclient satisfies the ethFetcher interface.
var (
	_ ethFetcher = (*rpcclient)(nil)

	bigZero = new(big.Int)
)

type ContextCaller interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type rpcclient struct {
	net dex.Network
	ep  *endpoint
	// ec wraps a *rpc.Client with some useful calls.
	ec *ethclient.Client
	// caller is a client for raw calls not implemented by *ethclient.Client.
	caller ContextCaller
	// swapContract is the current ETH swapContract.
	swapContract swapContract

	// tokens are tokeners for loaded tokens. tokens is not protected by a
	// mutex, as it is expected that the caller will connect and place calls to
	// loadToken sequentially in the same thread during initialization.
	tokens map[uint32]*tokener
}

func newRPCClient(net dex.Network, ep *endpoint) *rpcclient {
	return &rpcclient{
		net:    net,
		tokens: make(map[uint32]*tokener),
		ep:     ep,
	}
}

// connect connects to an ipc socket. It then wraps ethclient's client and
// bundles commands in a form we can easily use.
func (c *rpcclient) connect(ctx context.Context) error {
	var client *rpc.Client
	var err error
	addr := c.ep.addr
	if strings.HasSuffix(addr, ".ipc") {
		client, err = rpc.DialIPC(ctx, addr)
		if err != nil {
			return err
		}
	} else {
		var wsURL *url.URL
		wsURL, err = url.Parse(addr)
		if err != nil {
			return fmt.Errorf("Failed to parse url %q", addr)
		}
		wsURL.Scheme = "ws"
		var authFn func(h http.Header) error
		authFn, err = dexeth.JWTHTTPAuthFn(c.ep.jwt)
		if err != nil {
			return fmt.Errorf("unable to create auth function: %v", err)
		}
		client, err = rpc.DialOptions(ctx, wsURL.String(), rpc.WithHTTPAuth(authFn))
		if err != nil {
			return err
		}
	}

	c.ec = ethclient.NewClient(client)

	netAddrs, found := dexeth.ContractAddresses[ethContractVersion]
	if !found {
		return fmt.Errorf("no contract address for eth version %d", ethContractVersion)
	}
	contractAddr, found := netAddrs[c.net]
	if !found {
		return fmt.Errorf("no contract address for eth version %d on %s", ethContractVersion, c.net)
	}

	es, err := swapv0.NewETHSwap(contractAddr, c.ec)
	if err != nil {
		return fmt.Errorf("unable to find swap contract: %v", err)
	}
	c.swapContract = &swapSourceV0{es}
	c.caller = client
	return nil
}

// shutdown shuts down the client.
func (c *rpcclient) shutdown() {
	if c.ec != nil {
		c.ec.Close()
	}
}

func (c *rpcclient) loadToken(ctx context.Context, assetID uint32) error {
	tkn, err := newTokener(ctx, assetID, c.net, c.ec)
	if err != nil {
		return fmt.Errorf("error constructing ERC20Swap: %w", err)
	}

	c.tokens[assetID] = tkn
	return nil
}

func (c *rpcclient) withTokener(assetID uint32, f func(*tokener) error) error {
	tkn, found := c.tokens[assetID]
	if !found {
		return fmt.Errorf("no swap source for asset %d", assetID)
	}
	return f(tkn)
}

// bestHeader gets the best header at the time of calling.
func (c *rpcclient) bestHeader(ctx context.Context) (*types.Header, error) {
	bn, err := c.ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return c.ec.HeaderByNumber(ctx, big.NewInt(int64(bn)))
}

// headerByHeight gets the best header at height.
func (c *rpcclient) headerByHeight(ctx context.Context, height uint64) (*types.Header, error) {
	return c.ec.HeaderByNumber(ctx, big.NewInt(int64(height)))
}

// suggestGasTipCap retrieves the currently suggested priority fee to allow a
// timely execution of a transaction.
func (c *rpcclient) suggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return c.ec.SuggestGasTipCap(ctx)
}

// syncProgress return the current sync progress. Returns no error and nil when not syncing.
func (c *rpcclient) syncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	return c.ec.SyncProgress(ctx)
}

// blockNumber gets the chain length at the time of calling.
func (c *rpcclient) blockNumber(ctx context.Context) (uint64, error) {
	return c.ec.BlockNumber(ctx)
}

// swap gets a swap keyed by secretHash in the contract.
func (c *rpcclient) swap(ctx context.Context, assetID uint32, secretHash [32]byte) (state *dexeth.SwapState, err error) {
	if assetID == BipID {
		return c.swapContract.Swap(ctx, secretHash)
	}
	return state, c.withTokener(assetID, func(tkn *tokener) error {
		state, err = tkn.Swap(ctx, secretHash)
		return err
	})
}

// transaction gets the transaction that hashes to hash from the chain or
// mempool. Errors if tx does not exist.
func (c *rpcclient) transaction(ctx context.Context, hash common.Hash) (tx *types.Transaction, isMempool bool, err error) {
	return c.ec.TransactionByHash(ctx, hash)
}

// accountBalance gets the account balance, including the effects of known
// unmined transactions.
func (c *rpcclient) accountBalance(ctx context.Context, assetID uint32, addr common.Address) (*big.Int, error) {
	if assetID == BipID {
		return c.ec.BalanceAt(ctx, addr, nil)
	}

	bal := new(big.Int)
	return bal, c.withTokener(assetID, func(tkn *tokener) error {
		var err error
		bal, err = tkn.balanceOf(ctx, addr)
		return err
	})
}
