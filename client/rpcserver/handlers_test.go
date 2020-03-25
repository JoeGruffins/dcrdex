// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package rpcserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"decred.org/dcrdex/dex/msgjson"
)

func verifyResponse(payload *msgjson.ResponsePayload, res interface{}, wantErrCode int) error {
	if wantErrCode != -1 {
		if payload.Error.Code != wantErrCode {
			return errors.New("wrong error code")
		}
	} else {
		if payload.Error != nil {
			return fmt.Errorf("unexpected error: %v", payload.Error)
		}
	}
	if err := json.Unmarshal(payload.Result, res); err != nil {
		return errors.New("unable to unmarshal res")
	}
	return nil
}

type Dummy struct {
	Status string
}

func TestCreateResponse(t *testing.T) {
	tests := []struct {
		name        string
		res         interface{}
		resErr      *msgjson.Error
		wantErrCode int
	}{{
		name:        "ok",
		res:         Dummy{"ok"},
		resErr:      nil,
		wantErrCode: -1,
	}, {
		name:        "parse error",
		res:         "",
		resErr:      msgjson.NewError(msgjson.RPCParseError, "failed to encode response"),
		wantErrCode: msgjson.RPCParseError,
	}}

	for _, test := range tests {
		payload := createResponse(test.name, &test.res, test.resErr)
		if err := verifyResponse(payload, &test.res, test.wantErrCode); err != nil {
			t.Fatal(err)
		}

	}
}

func TestHelpMsgs(t *testing.T) {
	// routes and helpMsgs must have the same keys.
	if len(routes) != len(helpMsgs) {
		t.Fatal("routes and helpMsgs have different number of routes")
	}
	for k := range routes {
		if _, exists := helpMsgs[k]; !exists {
			t.Fatalf("%v exists in routes but not in helpMsgs", k)
		}
	}
}

func TestListCommands(t *testing.T) {
	res := ListCommands()
	if res == "" {
		t.Fatal("unable to parse helpMsgs")
	}
	want := ""
	for _, r := range sortHelpKeys() {
		want += r
		want += " "
		want += helpMsgs[r][0]
		want += "\n"
	}
	if res != want[:len(want)-1] {
		t.Fatalf("wanted %s but got %s", want, res)
	}
}

func TestCommandUsage(t *testing.T) {
	for r, msg := range helpMsgs {
		res, err := CommandUsage(r)
		if err != nil {
			t.Fatalf("unexpected error for command %s", r)
		}
		want := r
		want += " "
		want += msg[0]
		want += "\n\n"
		want += msg[1]
		if res != want {
			t.Fatalf("wanted %s but got %s for usage of %s", want, res, r)
		}
	}
	if _, err := CommandUsage("never make this command"); !errors.Is(err, ErrUnknownCmd) {
		t.Fatal("expected error for bogus command")
	}
}

func TestHandleHelp(t *testing.T) {
	tests := []struct {
		name        string
		arg         interface{}
		wantErrCode int
	}{{
		name:        "ok no arg",
		arg:         "",
		wantErrCode: -1,
	}, {
		name:        "ok arg",
		arg:         "version",
		wantErrCode: -1,
	}, {
		name:        "unknown route",
		arg:         "versio",
		wantErrCode: msgjson.RPCUnknownRoute,
	}, {
		name:        "argument wrong type",
		arg:         2,
		wantErrCode: msgjson.RPCParseError,
	}}
	for _, test := range tests {
		msg := new(msgjson.Message)
		reqPayload, err := json.Marshal(test.arg)
		if err != nil {
			t.Fatal(err)
		}
		msg.Payload = reqPayload
		payload := handleHelp(nil, msg)
		res := ""
		if err := verifyResponse(payload, &res, test.wantErrCode); err != nil {
			t.Fatal(err)
		}
	}
}

func TestHandleVersion(t *testing.T) {
	msg := new(msgjson.Message)
	payload := handleVersion(nil, msg)
	res := ""
	if err := verifyResponse(payload, &res, -1); err != nil {
		t.Fatal(err)
	}
}

func TestHandlePreRegister(t *testing.T) {
	tests := []struct {
		name           string
		arg            interface{}
		preRegisterFee uint64
		preRegisterErr error
		wantErrCode    int
	}{{
		name:           "ok",
		arg:            "dex",
		preRegisterFee: 5,
		wantErrCode:    -1,
	}, {
		name:        "argument wrong type",
		arg:         2,
		wantErrCode: msgjson.RPCParseError,
	}, {
		name:           "core.PreRegister error",
		arg:            "dex",
		preRegisterFee: 5,
		preRegisterErr: errors.New("error"),
		wantErrCode:    msgjson.RPCErrorUnspecified,
	}}
	for _, test := range tests {
		msg := new(msgjson.Message)
		reqPayload, err := json.Marshal(test.arg)
		if err != nil {
			t.Fatal(err)
		}
		msg.Payload = reqPayload
		tc := &TCore{preRegisterFee: test.preRegisterFee, preRegisterErr: test.preRegisterErr}
		r := &RPCServer{core: tc}
		payload := handlePreRegister(r, msg)
		res := new(preRegisterResponse)
		if err := verifyResponse(payload, &res, test.wantErrCode); err != nil {
			t.Fatal(err)
		}
		if test.wantErrCode == -1 && res.Fee != test.preRegisterFee {
			t.Fatalf("wanted registration fee %d but got %d for test %s", test.preRegisterFee, res.Fee, test.name)
		}
	}
}
