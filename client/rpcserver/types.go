// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package rpcserver

import (
	"errors"
	"fmt"
	"strconv"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex/encode"
)

// An orderID is a 256 bit number encoded as a hex string.
const orderIdLen = 64

var (
	// errArgs is wrapped when arguments to the known command cannot be parsed.
	errArgs = errors.New("unable to parse arguments")
)

// RawParams is used for all server requests.
type RawParams struct {
	PWArgs []encode.PassBytes `json:"PWArgs"`
	Args   []string           `json:"args"`
}

// versionResponse holds a semver version JSON object.
type versionResponse struct {
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
	Patch uint32 `json:"patch"`
}

// String satisfies the Stringer interface.
func (vr versionResponse) String() string {
	return fmt.Sprintf("%d.%d.%d", vr.Major, vr.Minor, vr.Patch)
}

// getFeeResponse is used when responding to the getfee route.
type getFeeResponse struct {
	Fee uint64 `json:"fee"`
}

// tradeResponse is used when responding to the trade route.
type tradeResponse struct {
	OrderID string `json:"orderID"`
	Sig     string `json:"sig"`
	Stamp   uint64 `json:"stamp"`
}

// openWalletForm is information necessary to open a wallet.
type openWalletForm struct {
	AssetID uint32           `json:"assetID"`
	AppPass encode.PassBytes `json:"appPass"`
}

// newWalletForm is information necessary to create a new wallet.
type newWalletForm struct {
	AssetID    uint32           `json:"assetID"`
	Account    string           `json:"account"`
	ConfigText string           `json:"config"`
	WalletPass encode.PassBytes `json:"walletPass"`
	AppPass    encode.PassBytes `json:"appPass"`
}

// helpForm is information necessary to obtain help.
type helpForm struct {
	HelpWith         string `json:"helpwith"`
	IncludePasswords bool   `json:"includepasswords"`
}

// tradeForm is information necessary to trade.
type tradeForm struct {
	AppPass encode.PassBytes
	SrvForm *core.TradeForm
}

// cancelForm is information necessary to cancel a trade.
type cancelForm struct {
	AppPass encode.PassBytes `json:"appPass"`
	OrderID string           `json:"orderID"`
}

// checkNArgs checks that args and pwArgs are the correct length.
func checkNArgs(params *RawParams, nPWArgs, nArgs []int) error {
	// For want, one integer indicates an exact match, two are the min and max.
	check := func(have int, want []int) error {
		if len(want) == 1 {
			if want[0] != have {
				return fmt.Errorf("%w: wanted %d but got %d", errArgs, want[0], have)
			}
		} else {
			if have < want[0] || have > want[1] {
				return fmt.Errorf("%w: wanted between %d and %d but got %d", errArgs, want[0], want[1], have)
			}
		}
		return nil
	}
	if err := check(len(params.Args), nArgs); err != nil {
		return fmt.Errorf("arguments: %w", err)
	}
	if err := check(len(params.PWArgs), nPWArgs); err != nil {
		return fmt.Errorf("password arguments: %w", err)
	}
	return nil
}

func checkUIntArg(arg, name string, base, bitSize int) (uint64, error) {
	i, err := strconv.ParseUint(arg, base, bitSize)
	if err != nil {
		return i, fmt.Errorf("%w: cannot parse %s: %v", errArgs, name, err)
	}
	return i, nil
}

func checkBoolArg(arg, name string) (bool, error) {
	b, err := strconv.ParseBool(arg)
	if err != nil {
		return b, fmt.Errorf("%w: %s must be a boolean: %v", errArgs, name, err)
	}
	return b, nil
}

func parseHelpArgs(params *RawParams) (*helpForm, error) {
	if err := checkNArgs(params, []int{0}, []int{0, 2}); err != nil {
		return nil, err
	}
	var helpWith string
	if len(params.Args) > 0 {
		helpWith = params.Args[0]
	}
	var includePasswords bool
	if len(params.Args) > 1 {
		var err error
		includePasswords, err = checkBoolArg(params.Args[1], "includepasswords")
		if err != nil {
			return nil, err
		}
	}
	return &helpForm{
		HelpWith:         helpWith,
		IncludePasswords: includePasswords,
	}, nil
}

func parseInitArgs(params *RawParams) (encode.PassBytes, error) {
	if err := checkNArgs(params, []int{1}, []int{0}); err != nil {
		return nil, err
	}
	return params.PWArgs[0], nil
}

func parseLoginArgs(params *RawParams) (encode.PassBytes, error) {
	if err := checkNArgs(params, []int{1}, []int{0}); err != nil {
		return nil, err
	}
	return params.PWArgs[0], nil
}

func parseNewWalletArgs(params *RawParams) (*newWalletForm, error) {
	if err := checkNArgs(params, []int{2}, []int{2, 3}); err != nil {
		return nil, err
	}
	assetID, err := checkUIntArg(params.Args[0], "assetID", 10, 32)
	if err != nil {
		return nil, err
	}
	req := &newWalletForm{
		AppPass:    params.PWArgs[0],
		WalletPass: params.PWArgs[1],
		AssetID:    uint32(assetID),
		Account:    params.Args[1],
	}
	if len(params.Args) > 2 {
		req.ConfigText = params.Args[2]
	}
	return req, nil
}

func parseOpenWalletArgs(params *RawParams) (*openWalletForm, error) {
	if err := checkNArgs(params, []int{1}, []int{1}); err != nil {
		return nil, err
	}
	assetID, err := checkUIntArg(params.Args[0], "assetID", 10, 32)
	if err != nil {
		return nil, err
	}
	req := &openWalletForm{AppPass: params.PWArgs[0], AssetID: uint32(assetID)}
	return req, nil
}

func parseCloseWalletArgs(params *RawParams) (uint32, error) {
	if err := checkNArgs(params, []int{0}, []int{1}); err != nil {
		return 0, err
	}
	assetID, err := checkUIntArg(params.Args[0], "assetID", 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(assetID), nil
}

func parseGetFeeArgs(params *RawParams) (host, cert string, err error) {
	if err := checkNArgs(params, []int{0}, []int{1, 2}); err != nil {
		return "", "", err
	}
	if len(params.Args) == 1 {
		return params.Args[0], "", nil
	}
	return params.Args[0], params.Args[1], nil
}

func parseRegisterArgs(params *RawParams) (*core.RegisterForm, error) {
	if err := checkNArgs(params, []int{1}, []int{2, 3}); err != nil {
		return nil, err
	}
	fee, err := checkUIntArg(params.Args[1], "fee", 10, 64)
	if err != nil {
		return nil, err
	}
	cert := ""
	if len(params.Args) > 2 {
		cert = params.Args[2]
	}
	req := &core.RegisterForm{
		AppPass: params.PWArgs[0],
		Addr:    params.Args[0],
		Fee:     fee,
		Cert:    cert,
	}
	return req, nil
}

func parseTradeArgs(params *RawParams) (*tradeForm, error) {
	if err := checkNArgs(params, []int{1}, []int{8}); err != nil {
		return nil, err
	}
	isLimit, err := checkBoolArg(params.Args[1], "isLimit")
	if err != nil {
		return nil, err
	}
	sell, err := checkBoolArg(params.Args[2], "sell")
	if err != nil {
		return nil, err
	}
	base, err := checkUIntArg(params.Args[3], "base", 10, 32)
	if err != nil {
		return nil, err
	}
	quote, err := checkUIntArg(params.Args[4], "quote", 10, 32)
	if err != nil {
		return nil, err
	}
	qty, err := checkUIntArg(params.Args[5], "qty", 10, 64)
	if err != nil {
		return nil, err
	}
	rate, err := checkUIntArg(params.Args[6], "rate", 10, 64)
	if err != nil {
		return nil, err
	}
	tifnow, err := checkBoolArg(params.Args[7], "immediate")
	if err != nil {
		return nil, err
	}
	req := &tradeForm{
		AppPass: params.PWArgs[0],
		SrvForm: &core.TradeForm{
			Host:    params.Args[0],
			IsLimit: isLimit,
			Sell:    sell,
			Base:    uint32(base),
			Quote:   uint32(quote),
			Qty:     qty,
			Rate:    rate,
			TifNow:  tifnow,
		},
	}
	return req, nil
}

func parseCancelArgs(params *RawParams) (*cancelForm, error) {
	if err := checkNArgs(params, []int{1}, []int{1}); err != nil {
		return nil, err
	}
	id := params.Args[0]
	if len(id) != orderIdLen {
		return nil, fmt.Errorf("%w: orderID has incorrect length", errArgs)
	}
	// We can only check up to 64 bits at a time. Break the 256 bit hex
	// into 4 sections and check each.
	div := orderIdLen / 4
	for i := 0; i < 4; i++ {
		if _, err := checkUIntArg(id[i*div:(i+1)*div], "orderID", 16, 64); err != nil {
			return nil, fmt.Errorf("%w: invalid order id hex", errArgs)
		}
	}
	return &cancelForm{AppPass: params.PWArgs[0], OrderID: id}, nil
}
