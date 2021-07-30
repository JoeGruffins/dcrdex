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

// ETHSwapSwap is an auto generated low-level Go binding around an user-defined struct.
type ETHSwapSwap struct {
	InitBlockNumber      *big.Int
	RefundBlockTimestamp *big.Int
	SecretHash           [32]byte
	Secret               [32]byte
	Initiator            common.Address
	Participant          common.Address
	Value                *big.Int
	State                uint8
}

// ETHSwapABI is the input ABI used to generate the binding from.
const ETHSwapABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"refundTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"secret\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"swap\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"initBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"refundBlockTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"secret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"enumETHSwap.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structETHSwap.Swap\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swaps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"initBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"refundBlockTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"secret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"enumETHSwap.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ETHSwapFuncSigs maps the 4-byte function signature to its string representation.
var ETHSwapFuncSigs = map[string]string{
	"ae052147": "initiate(uint256,bytes32,address)",
	"b31597ad": "redeem(bytes32,bytes32)",
	"7249fbb6": "refund(bytes32)",
	"76467cbd": "swap(bytes32)",
	"eb84e7f2": "swaps(bytes32)",
}

// ETHSwapBin is the compiled bytecode used for deploying new contracts.
var ETHSwapBin = "0x608060405234801561001057600080fd5b5061079c806100206000396000f3fe60806040526004361061004a5760003560e01c80637249fbb61461004f57806376467cbd14610071578063ae052147146100a7578063b31597ad146100ba578063eb84e7f2146100da575b600080fd5b34801561005b57600080fd5b5061006f61006a366004610595565b61015c565b005b34801561007d57600080fd5b5061009161008c366004610595565b610278565b60405161009e919061068b565b60405180910390f35b61006f6100b53660046105e9565b610363565b3480156100c657600080fd5b5061006f6100d53660046105c7565b610421565b3480156100e657600080fd5b506101486100f5366004610595565b6000602081905290815260409020805460018201546002830154600384015460048501546005860154600687015460079097015495969495939492936001600160a01b0392831693919092169160ff1688565b60405161009e9897969594939291906106fb565b32331461016857600080fd5b8033600160008381526020819052604090206007015460ff16600381111561019257610192610750565b1461019c57600080fd5b6000828152602081905260409020600401546001600160a01b038281169116146101c557600080fd5b600082815260208190526040902060010154428111156101e457600080fd5b6000848152602081905260408082206006015490513391908381818185875af1925050503d8060008114610234576040519150601f19603f3d011682016040523d82523d6000602084013e610239565b606091505b509091505060018115151461024d57600080fd5b600085815260208190526040902060070180546003919060ff19166001835b02179055505050505050565b6102bd6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290529060e082015290565b6000828152602081815260409182902082516101008101845281548152600182015492810192909252600281015492820192909252600380830154606083015260048301546001600160a01b03908116608084015260058401541660a0830152600683015460c0830152600783015491929160e084019160ff9091169081111561034957610349610750565b600381111561035a5761035a610750565b90525092915050565b826000341161037157600080fd5b6000811161037e57600080fd5b32331461038a57600080fd5b826000808281526020819052604090206007015460ff1660038111156103b2576103b2610750565b146103bc57600080fd5b6000848152602081905260409020438155600180820187905560028201869055600482018054336001600160a01b0319918216179091556005830180549091166001600160a01b0387161790553460068301556007909101805460ff1916828061026c565b32331461042d57600080fd5b808233600160008481526020819052604090206007015460ff16600381111561045857610458610750565b1461046257600080fd5b6000838152602081905260409020600501546001600160a01b0382811691161461048b57600080fd5b826002836040516020016104a191815260200190565b60408051601f19818403018152908290526104bb91610650565b602060405180830381855afa1580156104d8573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906104fb91906105ae565b1461050557600080fd5b60008481526020819052604080822060078101805460ff191660021790556006015490513391908381818185875af1925050503d8060008114610564576040519150601f19603f3d011682016040523d82523d6000602084013e610569565b606091505b509091505060018115151461057d57600080fd5b50505060009182525060208190526040902060030155565b6000602082840312156105a757600080fd5b5035919050565b6000602082840312156105c057600080fd5b5051919050565b600080604083850312156105da57600080fd5b50508035926020909101359150565b6000806000606084860312156105fe57600080fd5b833592506020840135915060408401356001600160a01b038116811461062357600080fd5b809150509250925092565b6004811061064c57634e487b7160e01b600052602160045260246000fd5b9052565b6000825160005b818110156106715760208186018101518583015201610657565b81811115610680576000828501525b509190910192915050565b60006101008201905082518252602083015160208301526040830151604083015260608301516060830152608083015160018060a01b0380821660808501528060a08601511660a0850152505060c083015160c083015260e08301516106f460e084018261062e565b5092915050565b8881526020810188905260408101879052606081018690526001600160a01b038581166080830152841660a082015260c08101839052610100810161074360e083018461062e565b9998505050505050505050565b634e487b7160e01b600052602160045260246000fdfea264697066735822122038c592bc3922997de5d495c92bcc2721f4e4f2103a2933a3df70fb2738bf792b64736f6c63430008060033"

// DeployETHSwap deploys a new Ethereum contract, binding an instance of ETHSwap to it.
func DeployETHSwap(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ETHSwap, error) {
	parsed, err := abi.JSON(strings.NewReader(ETHSwapABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ETHSwapBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ETHSwap{ETHSwapCaller: ETHSwapCaller{contract: contract}, ETHSwapTransactor: ETHSwapTransactor{contract: contract}, ETHSwapFilterer: ETHSwapFilterer{contract: contract}}, nil
}

// ETHSwap is an auto generated Go binding around an Ethereum contract.
type ETHSwap struct {
	ETHSwapCaller     // Read-only binding to the contract
	ETHSwapTransactor // Write-only binding to the contract
	ETHSwapFilterer   // Log filterer for contract events
}

// ETHSwapCaller is an auto generated read-only Go binding around an Ethereum contract.
type ETHSwapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHSwapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ETHSwapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHSwapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ETHSwapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHSwapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ETHSwapSession struct {
	Contract     *ETHSwap          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ETHSwapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ETHSwapCallerSession struct {
	Contract *ETHSwapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ETHSwapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ETHSwapTransactorSession struct {
	Contract     *ETHSwapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ETHSwapRaw is an auto generated low-level Go binding around an Ethereum contract.
type ETHSwapRaw struct {
	Contract *ETHSwap // Generic contract binding to access the raw methods on
}

// ETHSwapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ETHSwapCallerRaw struct {
	Contract *ETHSwapCaller // Generic read-only contract binding to access the raw methods on
}

// ETHSwapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ETHSwapTransactorRaw struct {
	Contract *ETHSwapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewETHSwap creates a new instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwap(address common.Address, backend bind.ContractBackend) (*ETHSwap, error) {
	contract, err := bindETHSwap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ETHSwap{ETHSwapCaller: ETHSwapCaller{contract: contract}, ETHSwapTransactor: ETHSwapTransactor{contract: contract}, ETHSwapFilterer: ETHSwapFilterer{contract: contract}}, nil
}

// NewETHSwapCaller creates a new read-only instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwapCaller(address common.Address, caller bind.ContractCaller) (*ETHSwapCaller, error) {
	contract, err := bindETHSwap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ETHSwapCaller{contract: contract}, nil
}

// NewETHSwapTransactor creates a new write-only instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwapTransactor(address common.Address, transactor bind.ContractTransactor) (*ETHSwapTransactor, error) {
	contract, err := bindETHSwap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ETHSwapTransactor{contract: contract}, nil
}

// NewETHSwapFilterer creates a new log filterer instance of ETHSwap, bound to a specific deployed contract.
func NewETHSwapFilterer(address common.Address, filterer bind.ContractFilterer) (*ETHSwapFilterer, error) {
	contract, err := bindETHSwap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ETHSwapFilterer{contract: contract}, nil
}

// bindETHSwap binds a generic wrapper to an already deployed contract.
func bindETHSwap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ETHSwapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ETHSwap *ETHSwapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ETHSwap.Contract.ETHSwapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ETHSwap *ETHSwapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHSwap.Contract.ETHSwapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ETHSwap *ETHSwapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ETHSwap.Contract.ETHSwapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ETHSwap *ETHSwapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ETHSwap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ETHSwap *ETHSwapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHSwap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ETHSwap *ETHSwapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ETHSwap.Contract.contract.Transact(opts, method, params...)
}

// Swap is a free data retrieval call binding the contract method 0x76467cbd.
//
// Solidity: function swap(bytes32 secretHash) view returns((uint256,uint256,bytes32,bytes32,address,address,uint256,uint8))
func (_ETHSwap *ETHSwapCaller) Swap(opts *bind.CallOpts, secretHash [32]byte) (ETHSwapSwap, error) {
	var out []interface{}
	err := _ETHSwap.contract.Call(opts, &out, "swap", secretHash)

	if err != nil {
		return *new(ETHSwapSwap), err
	}

	out0 := *abi.ConvertType(out[0], new(ETHSwapSwap)).(*ETHSwapSwap)

	return out0, err

}

// Swap is a free data retrieval call binding the contract method 0x76467cbd.
//
// Solidity: function swap(bytes32 secretHash) view returns((uint256,uint256,bytes32,bytes32,address,address,uint256,uint8))
func (_ETHSwap *ETHSwapSession) Swap(secretHash [32]byte) (ETHSwapSwap, error) {
	return _ETHSwap.Contract.Swap(&_ETHSwap.CallOpts, secretHash)
}

// Swap is a free data retrieval call binding the contract method 0x76467cbd.
//
// Solidity: function swap(bytes32 secretHash) view returns((uint256,uint256,bytes32,bytes32,address,address,uint256,uint8))
func (_ETHSwap *ETHSwapCallerSession) Swap(secretHash [32]byte) (ETHSwapSwap, error) {
	return _ETHSwap.Contract.Swap(&_ETHSwap.CallOpts, secretHash)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(uint256 initBlockNumber, uint256 refundBlockTimestamp, bytes32 secretHash, bytes32 secret, address initiator, address participant, uint256 value, uint8 state)
func (_ETHSwap *ETHSwapCaller) Swaps(opts *bind.CallOpts, arg0 [32]byte) (struct {
	InitBlockNumber      *big.Int
	RefundBlockTimestamp *big.Int
	SecretHash           [32]byte
	Secret               [32]byte
	Initiator            common.Address
	Participant          common.Address
	Value                *big.Int
	State                uint8
}, error) {
	var out []interface{}
	err := _ETHSwap.contract.Call(opts, &out, "swaps", arg0)

	outstruct := new(struct {
		InitBlockNumber      *big.Int
		RefundBlockTimestamp *big.Int
		SecretHash           [32]byte
		Secret               [32]byte
		Initiator            common.Address
		Participant          common.Address
		Value                *big.Int
		State                uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.InitBlockNumber = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RefundBlockTimestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.SecretHash = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Secret = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Initiator = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Participant = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.Value = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[7], new(uint8)).(*uint8)

	return *outstruct, err

}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(uint256 initBlockNumber, uint256 refundBlockTimestamp, bytes32 secretHash, bytes32 secret, address initiator, address participant, uint256 value, uint8 state)
func (_ETHSwap *ETHSwapSession) Swaps(arg0 [32]byte) (struct {
	InitBlockNumber      *big.Int
	RefundBlockTimestamp *big.Int
	SecretHash           [32]byte
	Secret               [32]byte
	Initiator            common.Address
	Participant          common.Address
	Value                *big.Int
	State                uint8
}, error) {
	return _ETHSwap.Contract.Swaps(&_ETHSwap.CallOpts, arg0)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 ) view returns(uint256 initBlockNumber, uint256 refundBlockTimestamp, bytes32 secretHash, bytes32 secret, address initiator, address participant, uint256 value, uint8 state)
func (_ETHSwap *ETHSwapCallerSession) Swaps(arg0 [32]byte) (struct {
	InitBlockNumber      *big.Int
	RefundBlockTimestamp *big.Int
	SecretHash           [32]byte
	Secret               [32]byte
	Initiator            common.Address
	Participant          common.Address
	Value                *big.Int
	State                uint8
}, error) {
	return _ETHSwap.Contract.Swaps(&_ETHSwap.CallOpts, arg0)
}

// Initiate is a paid mutator transaction binding the contract method 0xae052147.
//
// Solidity: function initiate(uint256 refundTimestamp, bytes32 secretHash, address participant) payable returns()
func (_ETHSwap *ETHSwapTransactor) Initiate(opts *bind.TransactOpts, refundTimestamp *big.Int, secretHash [32]byte, participant common.Address) (*types.Transaction, error) {
	return _ETHSwap.contract.Transact(opts, "initiate", refundTimestamp, secretHash, participant)
}

// Initiate is a paid mutator transaction binding the contract method 0xae052147.
//
// Solidity: function initiate(uint256 refundTimestamp, bytes32 secretHash, address participant) payable returns()
func (_ETHSwap *ETHSwapSession) Initiate(refundTimestamp *big.Int, secretHash [32]byte, participant common.Address) (*types.Transaction, error) {
	return _ETHSwap.Contract.Initiate(&_ETHSwap.TransactOpts, refundTimestamp, secretHash, participant)
}

// Initiate is a paid mutator transaction binding the contract method 0xae052147.
//
// Solidity: function initiate(uint256 refundTimestamp, bytes32 secretHash, address participant) payable returns()
func (_ETHSwap *ETHSwapTransactorSession) Initiate(refundTimestamp *big.Int, secretHash [32]byte, participant common.Address) (*types.Transaction, error) {
	return _ETHSwap.Contract.Initiate(&_ETHSwap.TransactOpts, refundTimestamp, secretHash, participant)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 secret, bytes32 secretHash) returns()
func (_ETHSwap *ETHSwapTransactor) Redeem(opts *bind.TransactOpts, secret [32]byte, secretHash [32]byte) (*types.Transaction, error) {
	return _ETHSwap.contract.Transact(opts, "redeem", secret, secretHash)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 secret, bytes32 secretHash) returns()
func (_ETHSwap *ETHSwapSession) Redeem(secret [32]byte, secretHash [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Redeem(&_ETHSwap.TransactOpts, secret, secretHash)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(bytes32 secret, bytes32 secretHash) returns()
func (_ETHSwap *ETHSwapTransactorSession) Redeem(secret [32]byte, secretHash [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Redeem(&_ETHSwap.TransactOpts, secret, secretHash)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 secretHash) returns()
func (_ETHSwap *ETHSwapTransactor) Refund(opts *bind.TransactOpts, secretHash [32]byte) (*types.Transaction, error) {
	return _ETHSwap.contract.Transact(opts, "refund", secretHash)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 secretHash) returns()
func (_ETHSwap *ETHSwapSession) Refund(secretHash [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Refund(&_ETHSwap.TransactOpts, secretHash)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 secretHash) returns()
func (_ETHSwap *ETHSwapTransactorSession) Refund(secretHash [32]byte) (*types.Transaction, error) {
	return _ETHSwap.Contract.Refund(&_ETHSwap.TransactOpts, secretHash)
}

// ReentryAttackABI is the input ABI used to generate the binding from.
const ReentryAttackABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"allYourBase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"areBelongToUs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"es\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"sh\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"refundTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"}],\"name\":\"setUsUpTheBomb\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// ReentryAttackFuncSigs maps the 4-byte function signature to its string representation.
var ReentryAttackFuncSigs = map[string]string{
	"8f110770": "allYourBase()",
	"627599ee": "areBelongToUs()",
	"8da5cb5b": "owner()",
	"b9ce28a4": "setUsUpTheBomb(address,bytes32,uint256,address)",
}

// ReentryAttackBin is the compiled bytecode used for deploying new contracts.
var ReentryAttackBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163317905561029b806100326000396000f3fe60806040526004361061003f5760003560e01c8063627599ee146100595780638da5cb5b1461006e5780638f110770146100aa578063b9ce28a4146100bf575b674563918244f40000471015610057576100576100d2565b005b34801561006557600080fd5b5061005761013b565b34801561007a57600080fd5b5060005461008e906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b3480156100b657600080fd5b506100576100d2565b6100576100cd36600461021f565b610178565b600254600154604051633924fddb60e11b81526001600160a01b0390921691637249fbb6916101079160040190815260200190565b600060405180830381600087803b15801561012157600080fd5b505af1158015610135573d6000803e3d6000fd5b50505050565b600080546040516001600160a01b03909116914780156108fc02929091818181858888f19350505050158015610175573d6000803e3d6000fd5b50565b600280546001600160a01b0319166001600160a01b03868116918217909255600185905560405163ae05214760e01b8152600481018590526024810186905291831660448301529063ae0521479034906064016000604051808303818588803b1580156101e457600080fd5b505af11580156101f8573d6000803e3d6000fd5b505050505050505050565b80356001600160a01b038116811461021a57600080fd5b919050565b6000806000806080858703121561023557600080fd5b61023e85610203565b9350602085013592506040850135915061025a60608601610203565b90509295919450925056fea2646970667358221220bf3973109f94de0b8b314e876e2a046f49d5fec689cecf584031c718bbf0834c64736f6c63430008060033"

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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReentryAttack *ReentryAttackCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReentryAttack.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReentryAttack *ReentryAttackSession) Owner() (common.Address, error) {
	return _ReentryAttack.Contract.Owner(&_ReentryAttack.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReentryAttack *ReentryAttackCallerSession) Owner() (common.Address, error) {
	return _ReentryAttack.Contract.Owner(&_ReentryAttack.CallOpts)
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

// SetUsUpTheBomb is a paid mutator transaction binding the contract method 0xb9ce28a4.
//
// Solidity: function setUsUpTheBomb(address es, bytes32 sh, uint256 refundTimestamp, address participant) payable returns()
func (_ReentryAttack *ReentryAttackTransactor) SetUsUpTheBomb(opts *bind.TransactOpts, es common.Address, sh [32]byte, refundTimestamp *big.Int, participant common.Address) (*types.Transaction, error) {
	return _ReentryAttack.contract.Transact(opts, "setUsUpTheBomb", es, sh, refundTimestamp, participant)
}

// SetUsUpTheBomb is a paid mutator transaction binding the contract method 0xb9ce28a4.
//
// Solidity: function setUsUpTheBomb(address es, bytes32 sh, uint256 refundTimestamp, address participant) payable returns()
func (_ReentryAttack *ReentryAttackSession) SetUsUpTheBomb(es common.Address, sh [32]byte, refundTimestamp *big.Int, participant common.Address) (*types.Transaction, error) {
	return _ReentryAttack.Contract.SetUsUpTheBomb(&_ReentryAttack.TransactOpts, es, sh, refundTimestamp, participant)
}

// SetUsUpTheBomb is a paid mutator transaction binding the contract method 0xb9ce28a4.
//
// Solidity: function setUsUpTheBomb(address es, bytes32 sh, uint256 refundTimestamp, address participant) payable returns()
func (_ReentryAttack *ReentryAttackTransactorSession) SetUsUpTheBomb(es common.Address, sh [32]byte, refundTimestamp *big.Int, participant common.Address) (*types.Transaction, error) {
	return _ReentryAttack.Contract.SetUsUpTheBomb(&_ReentryAttack.TransactOpts, es, sh, refundTimestamp, participant)
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
