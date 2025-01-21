// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package db

import (
	"context"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/dex/encode"
	"decred.org/dcrdex/dex/lexi"
)

type DB struct {
	*lexi.DB
	log          dex.Logger
	scores       *lexi.Table
	scoredIdx    *lexi.Index
	bonds        *lexi.Table
	bonderIdx    *lexi.Index
	bondStampIdx *lexi.Index
}

func New(dir string, log dex.Logger) (*DB, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creating db dir: %w", err)
	}
	db, err := lexi.New(&lexi.Config{
		Path: filepath.Join(dir, "reputation.db"),
		Log:  log,
	})
	if err != nil {
		return nil, fmt.Errorf("error constructing new db: %w", err)
	}

	// Scores. Keyed on (scorer peer ID, scored peer ID).
	scoreTable, err := db.Table("scores")
	if err != nil {
		return nil, fmt.Errorf("error constructing reputation table: %w", err)
	}
	// Scored peer index with timestamp sorting.
	scoredIdx, err := scoreTable.AddIndex("scored-stamp", func(_, v encoding.BinaryMarshaler) ([]byte, error) {
		s, is := v.(*dbScore)
		if !is {
			return nil, fmt.Errorf("wrong type %T", v)
		}
		tB := make([]byte, 8)
		binary.BigEndian.PutUint64(tB, uint64(s.stamp.UnixMilli()))
		return append(s.scored[:], tB...), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error constructing reputation index: %w", err)
	}
	scoredIdx.UseDefaultIterationOptions(lexi.WithReverse())

	// Bonds. Keyed on coin ID
	bondsTable, err := db.Table("bonds")
	if err != nil {
		return nil, fmt.Errorf("error constructing bonds table: %w", err)
	}
	// Retrieve bonds by peer ID.
	bonderIdx, err := bondsTable.AddIndex("bonder", func(_, v encoding.BinaryMarshaler) ([]byte, error) {
		b, is := v.(*dbBond)
		if !is {
			return nil, fmt.Errorf("wrong type %T", v)
		}
		return b.PeerID[:], nil
	})
	if err != nil {
		return nil, fmt.Errorf("error constructing bonder index: %w", err)
	}
	// We'll periodically prune expired bonds.
	bondStampIdx, err := bondsTable.AddIndex("bond-stamp", func(_, v encoding.BinaryMarshaler) ([]byte, error) {
		b, is := v.(*dbBond)
		if !is {
			return nil, fmt.Errorf("wrong type %T", v)
		}
		return encode.Uint64Bytes(uint64(b.Expiration.Unix())), nil
	})
	bondStampIdx.UseDefaultIterationOptions(lexi.WithReverse())
	if err != nil {
		return nil, fmt.Errorf("error constructing bond stamp index: %w", err)
	}
	return &DB{
		DB:           db,
		scores:       scoreTable,
		scoredIdx:    scoredIdx,
		log:          log,
		bonds:        bondsTable,
		bonderIdx:    bonderIdx,
		bondStampIdx: bondStampIdx,
	}, nil
}

func (db *DB) NewOrderBook(baseID, quoteID uint32) (*OrderBook, error) {
	bSym := dex.BipIDSymbol(baseID)
	qSym := dex.BipIDSymbol(quoteID)
	if bSym == "" || qSym == "" {
		return nil, errors.New("could not find base or quote symbol")
	}
	prefix := fmt.Sprintf("%s-%s", bSym, qSym)
	orderBookTable, err := db.Table(fmt.Sprintf("%s-orderbook", prefix))
	if err != nil {
		return nil, fmt.Errorf("error constructing orderbook table: %w", err)
	}

	orderBookOrderIDIdx, err := orderBookTable.AddIndex(fmt.Sprintf("%s-orderbook-orderid", prefix), func(_, v encoding.BinaryMarshaler) ([]byte, error) {
		o, is := v.(*OrderUpdate)
		if !is {
			return nil, fmt.Errorf("wrong type %T", v)
		}
		oID := o.ID()
		return oID[:], nil
	})
	if err != nil {
		return nil, err
	}

	orderBookStampIdx, err := orderBookTable.AddIndex(fmt.Sprintf("%s-orderbook-stamp", prefix), func(_, v encoding.BinaryMarshaler) ([]byte, error) {
		o, is := v.(*OrderUpdate)
		if !is {
			return nil, fmt.Errorf("wrong type %T", v)
		}
		tB := make([]byte, 8)
		binary.BigEndian.PutUint64(tB, uint64(o.Stamp.UnixMilli()))
		return tB, nil
	})
	if err != nil {
		return nil, err
	}

	orderBookSellRateIdx, err := orderBookTable.AddIndex(fmt.Sprintf("%s-orderbook-sellrate", prefix), func(_, v encoding.BinaryMarshaler) ([]byte, error) {
		o, is := v.(*OrderUpdate)
		if !is {
			return nil, fmt.Errorf("wrong type %T", v)
		}
		srB := make([]byte, 1+8)
		if o.Sell {
			srB[0] = 1
		}
		binary.BigEndian.PutUint64(srB[1:], o.Rate)
		return srB, nil
	})
	if err != nil {
		return nil, err
	}

	return &OrderBook{
		orderBook:            orderBookTable,
		orderBookOrderIDIdx:  orderBookOrderIDIdx,
		orderBookStampIdx:    orderBookStampIdx,
		orderBookSellRateIdx: orderBookSellRateIdx,
	}, nil
}

func (db *DB) Connect(ctx context.Context) (*sync.WaitGroup, error) {
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case <-time.After(time.Hour):
				db.pruneOldBonds()
			case <-ctx.Done():
				return
			}
		}
	}()

	return &wg, nil
}

func (db *DB) pruneOldBonds() {
	if err := db.bondStampIdx.Iterate(nil, func(it *lexi.Iter) error {
		return it.Delete()
	}, lexi.WithSeek(encode.Uint64Bytes(uint64(time.Now().Unix()))), lexi.WithUpdate()); err != nil {
		db.log.Errorf("Error pruning bonds: %v", err)
	}
}
