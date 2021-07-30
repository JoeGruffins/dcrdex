// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package reentryattack

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ReentryAttackABI is the input ABI used to generate the binding from.
const ReentryAttackABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"allYourBase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"areBelongToUs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"es\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"sh\",\"type\":\"bytes32\"}],\"name\":\"setUsUpTheBomb\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ReentryAttackFuncSigs maps the 4-byte function signature to its string representation.
var ReentryAttackFuncSigs = map[string]string{
	"8f110770": "allYourBase()",
	"627599ee": "areBelongToUs()",
	"5b5630ac": "setUsUpTheBomb(address,bytes32)",
}

// ReentryAttackBin is the compiled bytecode used for deploying new contracts.
var ReentryAttackBin = "0x608060405234801561001057600080fd5b5061025d806100206000396000f3fe6080604052600436106100345760003560e01c80635b5630ac1461005c578063627599ee146100ab5780638f110770146100c0575b600054673782dace9d9000006001600160a01b0390911631101561005a5761005a6100d1565b005b34801561006857600080fd5b5061005a6100773660046101b4565b60008054336001600160a01b031991821617909155600180549091166001600160a01b039390931692909217909155600255565b3480156100b757600080fd5b5061005a610177565b3480156100cc57600080fd5b5061005a5b6001546002546040516001600160a01b039092169166038d7ea4c68000916100ff9160240190815260200190565b60408051601f198184030181529181526020820180516001600160e01b0316633924fddb60e11b1790525161013491906101ec565b60006040518083038160008787f1925050503d8060008114610172576040519150601f19603f3d011682016040523d82523d6000602084013e505050565b505050565b600080546040516001600160a01b03909116914780156108fc02929091818181858888f193505050501580156101b1573d6000803e3d6000fd5b50565b600080604083850312156101c757600080fd5b82356001600160a01b03811681146101de57600080fd5b946020939093013593505050565b6000825160005b8181101561020d57602081860181015185830152016101f3565b8181111561021c576000828501525b50919091019291505056fea264697066735822122024cc2ae8c80886278602b7b7d9e3a414030648c38b5ed62dac79a007ae02edc664736f6c63430008060033"

// DeployReentryAttack deploys a new Ethereum contract, binding an instance of ReentryAttack to it.
func DeployReentryAttack(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ReentryAttack, error) {
	parsed, err := abi.JSON(strings.NewReader(ReentryAttackABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ReentryAttackBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReentryAttack{ReentryAttackCaller: ReentryAttackCaller{contract: contract}, ReentryAttackTransactor: ReentryAttackTransactor{contract: contract}, ReentryAttackFilterer: ReentryAttackFilterer{contract: contract}}, nil
}

// ReentryAttack is an auto generated Go binding around an Ethereum contract.
type ReentryAttack struct {
	ReentryAttackCaller     // Read-only binding to the contract
	ReentryAttackTransactor // Write-only binding to the contract
	ReentryAttackFilterer   // Log filterer for contract events
}

// ReentryAttackCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReentryAttackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentryAttackTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReentryAttackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentryAttackFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReentryAttackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentryAttackSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReentryAttackSession struct {
	Contract     *ReentryAttack    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReentryAttackCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReentryAttackCallerSession struct {
	Contract *ReentryAttackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ReentryAttackTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReentryAttackTransactorSession struct {
	Contract     *ReentryAttackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ReentryAttackRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReentryAttackRaw struct {
	Contract *ReentryAttack // Generic contract binding to access the raw methods on
}

// ReentryAttackCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReentryAttackCallerRaw struct {
	Contract *ReentryAttackCaller // Generic read-only contract binding to access the raw methods on
}

// ReentryAttackTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReentryAttackTransactorRaw struct {
	Contract *ReentryAttackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReentryAttack creates a new instance of ReentryAttack, bound to a specific deployed contract.
func NewReentryAttack(address common.Address, backend bind.ContractBackend) (*ReentryAttack, error) {
	contract, err := bindReentryAttack(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReentryAttack{ReentryAttackCaller: ReentryAttackCaller{contract: contract}, ReentryAttackTransactor: ReentryAttackTransactor{contract: contract}, ReentryAttackFilterer: ReentryAttackFilterer{contract: contract}}, nil
}

// NewReentryAttackCaller creates a new read-only instance of ReentryAttack, bound to a specific deployed contract.
func NewReentryAttackCaller(address common.Address, caller bind.ContractCaller) (*ReentryAttackCaller, error) {
	contract, err := bindReentryAttack(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReentryAttackCaller{contract: contract}, nil
}

// NewReentryAttackTransactor creates a new write-only instance of ReentryAttack, bound to a specific deployed contract.
func NewReentryAttackTransactor(address common.Address, transactor bind.ContractTransactor) (*ReentryAttackTransactor, error) {
	contract, err := bindReentryAttack(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReentryAttackTransactor{contract: contract}, nil
}

// NewReentryAttackFilterer creates a new log filterer instance of ReentryAttack, bound to a specific deployed contract.
func NewReentryAttackFilterer(address common.Address, filterer bind.ContractFilterer) (*ReentryAttackFilterer, error) {
	contract, err := bindReentryAttack(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReentryAttackFilterer{contract: contract}, nil
}

// bindReentryAttack binds a generic wrapper to an already deployed contract.
func bindReentryAttack(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReentryAttackABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReentryAttack *ReentryAttackRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReentryAttack.Contract.ReentryAttackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReentryAttack *ReentryAttackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentryAttack.Contract.ReentryAttackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReentryAttack *ReentryAttackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReentryAttack.Contract.ReentryAttackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReentryAttack *ReentryAttackCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReentryAttack.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReentryAttack *ReentryAttackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentryAttack.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReentryAttack *ReentryAttackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReentryAttack.Contract.contract.Transact(opts, method, params...)
}

// AllYourBase is a paid mutator transaction binding the contract method 0x8f110770.
//
// Solidity: function allYourBase() returns()
func (_ReentryAttack *ReentryAttackTransactor) AllYourBase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentryAttack.contract.Transact(opts, "allYourBase")
}

// AllYourBase is a paid mutator transaction binding the contract method 0x8f110770.
//
// Solidity: function allYourBase() returns()
func (_ReentryAttack *ReentryAttackSession) AllYourBase() (*types.Transaction, error) {
	return _ReentryAttack.Contract.AllYourBase(&_ReentryAttack.TransactOpts)
}

// AllYourBase is a paid mutator transaction binding the contract method 0x8f110770.
//
// Solidity: function allYourBase() returns()
func (_ReentryAttack *ReentryAttackTransactorSession) AllYourBase() (*types.Transaction, error) {
	return _ReentryAttack.Contract.AllYourBase(&_ReentryAttack.TransactOpts)
}

// AreBelongToUs is a paid mutator transaction binding the contract method 0x627599ee.
//
// Solidity: function areBelongToUs() returns()
func (_ReentryAttack *ReentryAttackTransactor) AreBelongToUs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentryAttack.contract.Transact(opts, "areBelongToUs")
}

// AreBelongToUs is a paid mutator transaction binding the contract method 0x627599ee.
//
// Solidity: function areBelongToUs() returns()
func (_ReentryAttack *ReentryAttackSession) AreBelongToUs() (*types.Transaction, error) {
	return _ReentryAttack.Contract.AreBelongToUs(&_ReentryAttack.TransactOpts)
}

// AreBelongToUs is a paid mutator transaction binding the contract method 0x627599ee.
//
// Solidity: function areBelongToUs() returns()
func (_ReentryAttack *ReentryAttackTransactorSession) AreBelongToUs() (*types.Transaction, error) {
	return _ReentryAttack.Contract.AreBelongToUs(&_ReentryAttack.TransactOpts)
}

// SetUsUpTheBomb is a paid mutator transaction binding the contract method 0x5b5630ac.
//
// Solidity: function setUsUpTheBomb(address es, bytes32 sh) returns()
func (_ReentryAttack *ReentryAttackTransactor) SetUsUpTheBomb(opts *bind.TransactOpts, es common.Address, sh [32]byte) (*types.Transaction, error) {
	return _ReentryAttack.contract.Transact(opts, "setUsUpTheBomb", es, sh)
}

// SetUsUpTheBomb is a paid mutator transaction binding the contract method 0x5b5630ac.
//
// Solidity: function setUsUpTheBomb(address es, bytes32 sh) returns()
func (_ReentryAttack *ReentryAttackSession) SetUsUpTheBomb(es common.Address, sh [32]byte) (*types.Transaction, error) {
	return _ReentryAttack.Contract.SetUsUpTheBomb(&_ReentryAttack.TransactOpts, es, sh)
}

// SetUsUpTheBomb is a paid mutator transaction binding the contract method 0x5b5630ac.
//
// Solidity: function setUsUpTheBomb(address es, bytes32 sh) returns()
func (_ReentryAttack *ReentryAttackTransactorSession) SetUsUpTheBomb(es common.Address, sh [32]byte) (*types.Transaction, error) {
	return _ReentryAttack.Contract.SetUsUpTheBomb(&_ReentryAttack.TransactOpts, es, sh)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ReentryAttack *ReentryAttackTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ReentryAttack.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ReentryAttack *ReentryAttackSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ReentryAttack.Contract.Fallback(&_ReentryAttack.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ReentryAttack *ReentryAttackTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ReentryAttack.Contract.Fallback(&_ReentryAttack.TransactOpts, calldata)
}
