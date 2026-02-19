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

	v1Gases = &dexeth.Gases{
		// First swap used 54284 gas Recommended Gases.Swap = 70569
		Swap: 70_569,
		// 4 additional swaps averaged 29713 gas each. Recommended Gases.SwapAdd = 38625
		// [54284 84005 113702 143424 173135]
		SwapAdd: 38_625,
		// First redeem used 44899 gas. Recommended Gases.Redeem = 58368
		Redeem: 58_368,
		// 4 additional redeems averaged 13348 gas each. Recommended Gases.RedeemAdd = 17351
		// [44899 58254 71586 84931 98290]
		RedeemAdd: 17_351,
		// Average of 6 refunds: 47130. Recommended Gases.Refund = 61269
		// [47975 47987 47987 47987 47987 42859]
		Refund: 61_269,

		// Gasless redeem (testnet, v0.7 EntryPoint):
		// Verification (n=1..5): [83000 94000 105000 116000 127000]
		GaslessRedeemVerification: 107_900,
		// VerificationAdd: avg diff = 11000. Recommended = 14300
		GaslessRedeemVerificationAdd: 14_300,
		// PreVerification (n=1..5): [70000 76000 82000 88000 94000]
		GaslessRedeemPreVerification: 91_000,
		// PreVerificationAdd: avg diff = 6000. Recommended = 7800
		GaslessRedeemPreVerificationAdd: 7_800,
		// Call gas uses the contract's hard minimums from validateUserOp.
		// The EntryPoint passes callGasLimit directly to the inner call
		// (Exec.call), so the full amount is available to redeemAA.
		// Raw bundler estimates (n=1..5): [120000 145000 170000 195000 220000]
		GaslessRedeemCall:    100_000, // MIN_CALL_GAS_BASE (75k) + MIN_CALL_GAS_PER_REDEMPTION (25k)
		GaslessRedeemCallAdd: 25_000,  // MIN_CALL_GAS_PER_REDEMPTION
	}

	VersionedGases = map[uint32]*dexeth.Gases{
		1: v1Gases,
	}

	ContractAddresses = map[uint32]map[dex.Network]common.Address{
		0: {
			dex.Simnet: common.HexToAddress(""), // Filled in by MaybeReadSimnetAddrs
		},
		1: {
			dex.Testnet: common.HexToAddress("0xbeA9D54f2bD1F54e9130D17d834C1172eD314A39"), // txid: 0x89fdba013c832fced6cbac89b1057a400f4eb955bf174d102fba3517e78dbbbf
			dex.Mainnet: common.HexToAddress("0xf4D3c25017928c563A074C6c6880dC6787E19bE0"), // txid: 0x0384462e9bdd13b1233eef1982d8c5ed23ce0a79694b4ec0ca6c95c28764535d
			dex.Simnet:  common.HexToAddress(""), // Filled in by MaybeReadSimnetAddrs
		},
	}

	MultiBalanceAddresses = map[dex.Network]common.Address{}

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
				Address: common.HexToAddress("0x036CbD53842c5426634e7929541eC2318f3dCF7e"),
				SwapContracts: map[uint32]*dexeth.SwapContract{
					1: {
						Gas: dexeth.Gases{
							// Testnet measurements (Base Sepolia, 2026-02-20):
							// Swaps (n=1..5):   [103963 133635 163321 193019 222694]
							// Redeems (n=1..5): [59881 73188 86508 99830 113140]
							// Refunds (n=1..6): [68102 68004 68004 68016 68016 57814]
							// Approvals: [55785 55785]
							// Transfers: [62135]
							Swap:      135_151,
							SwapAdd:   38_586,
							Redeem:    77_845,
							RedeemAdd: 17_308,
							Refund:    86_223,
							Approve:   72_520,
							Transfer:  80_775,
						},
					},
				},
			},
			dex.Simnet: {
				Address: common.Address{},
				SwapContracts: map[uint32]*dexeth.SwapContract{
					0: {},
					1: {
						Gas: dexeth.Gases{
							Swap:      114_515,
							SwapAdd:   34_672,
							Redeem:    58_272,
							RedeemAdd: 14_207,
							Refund:    61_911,
							Approve:   58_180,
							Transfer:  66_961,
						},
					},
				},
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
				Address: common.HexToAddress("0xb72fdb9f8190d8e1141e6a8e9c0732b0f4d93c09"),
				SwapContracts: map[uint32]*dexeth.SwapContract{
					1: {
						Gas: dexeth.Gases{
							Swap:      135_151,
							SwapAdd:   38_586,
							Redeem:    77_845,
							RedeemAdd: 17_308,
							Refund:    86_223,
							Approve:   67_693,
							Transfer:  82_180,
						},
					},
				},
			},
			dex.Simnet: {
				Address: common.Address{},
				SwapContracts: map[uint32]*dexeth.SwapContract{
					0: {},
					1: {
						Gas: dexeth.Gases{
							Swap:      123_743,
							SwapAdd:   34_453,
							Redeem:    72_237,
							RedeemAdd: 13_928,
							Refund:    81_252,
							Approve:   67_693,
							Transfer:  82_180,
						},
					},
				},
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
				Address: common.HexToAddress("0x78c8587b0b4d50b3a2110bc8188eef195cfa7f11"),
				SwapContracts: map[uint32]*dexeth.SwapContract{
					1: {
						Gas: dexeth.Gases{
							Swap:      135_151,
							SwapAdd:   38_586,
							Redeem:    77_845,
							RedeemAdd: 17_308,
							Refund:    86_223,
							Approve:   58_180,
							Transfer:  66_961,
						},
					},
				},
			},
			dex.Simnet: {
				Address:       common.Address{},
				SwapContracts: map[uint32]*dexeth.SwapContract{0: {}, 1: {}},
			},
		},
	}
)

// EntryPoints is a map of network to the ERC-4337 entrypoint address.
var EntryPoints = map[dex.Network]common.Address{
	// dex.Simnet:  common.Address{}, // populated by MaybeReadSimnetAddrs
	dex.Mainnet: dexeth.CanonicalEntryPointV07,
	dex.Testnet: dexeth.CanonicalEntryPointV07,
}

// MaybeReadSimnetAddrs attempts to read the info files generated by the eth
// simnet harness to populate swap contract and token addresses in
// ContractAddresses and Tokens.
func MaybeReadSimnetAddrs() {
	dexeth.MaybeReadSimnetAddrsDir("base", ContractAddresses, MultiBalanceAddresses, EntryPoints, Tokens[usdcTokenID].NetTokens[dex.Simnet], Tokens[usdtTokenID].NetTokens[dex.Simnet])
}
