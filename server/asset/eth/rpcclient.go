// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Check that rpcclient satisfies the ethFetcher interface.
var _ ethFetcher = (*rpcclient)(nil)

type rpcclient struct {
	ec *ethclient.Client
}

const gweiFactor = 1e9

var gweiFactorBig = big.NewInt(gweiFactor)

// connect connects to an ipc socket. It then wraps ethclient's client and
// bundles commands in a form we can easil use.
func (c *rpcclient) connect(ctx context.Context, IPC string) error {
	client, err := rpc.DialIPC(ctx, IPC)
	if err != nil {
		return fmt.Errorf("unable to dial rpc: %v", err)
	}
	c.ec = ethclient.NewClient(client)
	return nil
}

// shutdown shuts down the client.
func (c *rpcclient) shutdown() {
	if c.ec != nil {
		c.ec.Close()
	}
}

// bestBlockHash gets the best blocks hash at the time of calling. Due to the
// speed of Ethereum blocks, this changes often.
func (c *rpcclient) bestBlockHash(ctx context.Context) (common.Hash, error) {
	header, err := c.bestHeader(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return header.Hash(), nil
}

// bestHeader gets the best header at the time of calling.
func (c *rpcclient) bestHeader(ctx context.Context) (*types.Header, error) {
	bn, err := c.ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	header, err := c.ec.HeaderByNumber(ctx, big.NewInt(int64(bn)))
	if err != nil {
		return nil, err
	}
	return header, nil
}

// block gets the block identified by hash.
func (c *rpcclient) block(ctx context.Context, hash common.Hash) (*types.Block, error) {
	block, err := c.ec.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// suggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (c *rpcclient) suggestGasPrice(ctx context.Context) (sgp *big.Int, err error) {
	// NOTE: geth will panic if this is called too soon after creating the
	// node.
	defer func() {
		if r := recover(); r != nil {
			sgp = nil
			err = fmt.Errorf("geth panic: %v", r)
		}
	}()
	sgp, err = c.ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return sgp, nil
}

// syncProgress return the current sync progress. Returns no error and nil when not syncing.
func (c *rpcclient) syncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	return c.ec.SyncProgress(ctx)
}

// blockNumber returns the current block number.
func (c *rpcclient) blockNumber(ctx context.Context) (uint64, error) {
	return c.ec.BlockNumber(ctx)
}
