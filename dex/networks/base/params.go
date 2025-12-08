// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package base

import (
	"decred.org/dcrdex/dex"
	dexeth "decred.org/dcrdex/dex/networks/eth"
	"github.com/ethereum/go-ethereum/common"
)

const (
	BaseBipID = 8453 // weth.base
)

var (
	UnitInfo = dex.UnitInfo{
		AtomicUnit: "gwei",
		Conventional: dex.Denomination{
			Unit:             "WETH",
			ConversionFactor: 1e9,
		},
		Alternatives: []dex.Denomination{
			{
				Unit:             "Szabos",
				ConversionFactor: 1e6,
			},
			{
				Unit:             "Finneys",
				ConversionFactor: 1e3,
			},
		},
		FeeRateDenom: "gas",
	}

	ContractAddresses = map[uint32]map[dex.Network]common.Address{
		1: {
			dex.Testnet: common.HexToAddress("0x4F8024dd716E4ec9AC8E0CD971616F73e3f6625d"), // txid: https://base-sepolia.blockscout.com/tx/0xa8137a8f84b28149f27045a4d6c15cceb784f3886f8a476790002e4da4728979
			dex.Mainnet: common.HexToAddress("0x4F8024dd716E4ec9AC8E0CD971616F73e3f6625d"), // txid: https://base-sepolia.blockscout.com/tx/0xa8137a8f84b28149f27045a4d6c15cceb784f3886f8a476790002e4da4728979
			dex.Simnet:  common.HexToAddress(""),                                           // Filled in by MaybeReadSimnetAddrs
		},
	}

	VersionedGases = map[uint32]*dexeth.Gases{}

	usdcTokenID, _ = dex.BipSymbolID("usdc.base")
	usdtTokenID, _ = dex.BipSymbolID("usdt.base")
	wbtcTokenID, _ = dex.BipSymbolID("wbtc.base")

	Tokens = map[uint32]*dexeth.Token{
		usdcTokenID: TokenUSDC,
		usdtTokenID: TokenUSDT,
		wbtcTokenID: TokenWBTC,
	}

	TokenUSDC = &dexeth.Token{
		EVMFactor: new(int64),
		Token: &dex.Token{
			ParentID: BaseBipID,
			Name:     "USDC",
			UnitInfo: dex.UnitInfo{
				AtomicUnit: "µUSD",
				Conventional: dex.Denomination{
					Unit:             "USDC",
					ConversionFactor: 1e6,
				},
				Alternatives: []dex.Denomination{
					{
						Unit:             "cents",
						ConversionFactor: 1e2,
					},
				},
				FeeRateDenom: "gas",
			},
		},
		NetTokens: map[dex.Network]*dexeth.NetToken{
			dex.Mainnet: {
				Address:       common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"),
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
			dex.Testnet: {
				Address:       common.HexToAddress("0x036CbD53842c5426634e7929541eC2318f3dCF7e"),
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
			dex.Simnet: {
				Address:       common.Address{},
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
		},
	}

	TokenUSDT = &dexeth.Token{
		EVMFactor: new(int64),
		Token: &dex.Token{
			ParentID: BaseBipID,
			Name:     "Tether",
			UnitInfo: dex.UnitInfo{
				AtomicUnit: "microUSD",
				Conventional: dex.Denomination{
					Unit:             "USDT",
					ConversionFactor: 1e6,
				},
			},
		},
		NetTokens: map[dex.Network]*dexeth.NetToken{
			dex.Mainnet: {
				Address:       common.HexToAddress("0xfde4C96c8593536E31F229EA8f37b2ADa2699bb2"),
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
			dex.Testnet: {
				// Is USDCT tether?
				Address:       common.HexToAddress("0xb72fdb9f8190d8e1141e6a8e9c0732b0f4d93c09"),
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
			dex.Simnet: {
				Address:       common.Address{},
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
		},
	}

	TokenWBTC = &dexeth.Token{
		EVMFactor: new(int64),
		Token: &dex.Token{
			ParentID: BaseBipID,
			Name:     "Wrapped Bitcoin",
			UnitInfo: dex.UnitInfo{
				AtomicUnit: "Sats",
				Conventional: dex.Denomination{
					Unit:             "WBTC",
					ConversionFactor: 1e8,
				},
				Alternatives: []dex.Denomination{
					{
						Unit:             "mWBTC",
						ConversionFactor: 1e5,
					},
					{
						Unit:             "µWBTC",
						ConversionFactor: 1e2,
					},
				},
				FeeRateDenom: "gas",
			},
		},
		NetTokens: map[dex.Network]*dexeth.NetToken{
			dex.Mainnet: {
				Address:       common.HexToAddress("0x0555E30da8f98308EdB960aa94C0Db47230d2B9c"),
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
			dex.Testnet: {
				Address:       common.HexToAddress("0x78c8587b0b4d50b3a2110bc8188eef195cfa7f11"),
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
			dex.Simnet: {
				Address:       common.Address{},
				SwapContracts: map[uint32]*dexeth.SwapContract{},
			},
		},
	}
)
