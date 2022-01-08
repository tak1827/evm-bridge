// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package client

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

// IBankABI is the input ABI used to generate the binding from.
const IBankABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiAmount\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ERC20Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ERC20Withdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenid\",\"type\":\"uint256\"}],\"name\":\"NFTDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenid\",\"type\":\"uint256\"}],\"name\":\"NFTWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiAmount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenid\",\"type\":\"uint256\"}],\"name\":\"depositNFT\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenid\",\"type\":\"uint256\"}],\"name\":\"withdrawNFT\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IBankFuncSigs maps the 4-byte function signature to its string representation.
var IBankFuncSigs = map[string]string{
	"f340fa01": "deposit(address)",
	"1cad5a40": "depositERC20(address,address,uint256)",
	"548bfbaa": "depositNFT(address,address,uint256)",
	"d9caed12": "withdraw(address,address,uint256)",
	"44004cc1": "withdrawERC20(address,address,uint256)",
	"bff863ec": "withdrawNFT(address,address,uint256)",
}

// IBank is an auto generated Go binding around an Ethereum contract.
type IBank struct {
	IBankCaller     // Read-only binding to the contract
	IBankTransactor // Write-only binding to the contract
	IBankFilterer   // Log filterer for contract events
}

// IBankCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBankCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBankTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBankTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBankFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBankFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBankSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBankSession struct {
	Contract     *IBank            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBankCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBankCallerSession struct {
	Contract *IBankCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IBankTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBankTransactorSession struct {
	Contract     *IBankTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBankRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBankRaw struct {
	Contract *IBank // Generic contract binding to access the raw methods on
}

// IBankCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBankCallerRaw struct {
	Contract *IBankCaller // Generic read-only contract binding to access the raw methods on
}

// IBankTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBankTransactorRaw struct {
	Contract *IBankTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBank creates a new instance of IBank, bound to a specific deployed contract.
func NewIBank(address common.Address, backend bind.ContractBackend) (*IBank, error) {
	contract, err := bindIBank(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBank{IBankCaller: IBankCaller{contract: contract}, IBankTransactor: IBankTransactor{contract: contract}, IBankFilterer: IBankFilterer{contract: contract}}, nil
}

// NewIBankCaller creates a new read-only instance of IBank, bound to a specific deployed contract.
func NewIBankCaller(address common.Address, caller bind.ContractCaller) (*IBankCaller, error) {
	contract, err := bindIBank(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBankCaller{contract: contract}, nil
}

// NewIBankTransactor creates a new write-only instance of IBank, bound to a specific deployed contract.
func NewIBankTransactor(address common.Address, transactor bind.ContractTransactor) (*IBankTransactor, error) {
	contract, err := bindIBank(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBankTransactor{contract: contract}, nil
}

// NewIBankFilterer creates a new log filterer instance of IBank, bound to a specific deployed contract.
func NewIBankFilterer(address common.Address, filterer bind.ContractFilterer) (*IBankFilterer, error) {
	contract, err := bindIBank(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBankFilterer{contract: contract}, nil
}

// bindIBank binds a generic wrapper to an already deployed contract.
func bindIBank(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IBankABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBank *IBankRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBank.Contract.IBankCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBank *IBankRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBank.Contract.IBankTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBank *IBankRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBank.Contract.IBankTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBank *IBankCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBank.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBank *IBankTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBank.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBank *IBankTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBank.Contract.contract.Transact(opts, method, params...)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address payee) payable returns()
func (_IBank *IBankTransactor) Deposit(opts *bind.TransactOpts, payee common.Address) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "deposit", payee)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address payee) payable returns()
func (_IBank *IBankSession) Deposit(payee common.Address) (*types.Transaction, error) {
	return _IBank.Contract.Deposit(&_IBank.TransactOpts, payee)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address payee) payable returns()
func (_IBank *IBankTransactorSession) Deposit(payee common.Address) (*types.Transaction, error) {
	return _IBank.Contract.Deposit(&_IBank.TransactOpts, payee)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x1cad5a40.
//
// Solidity: function depositERC20(address token, address sender, uint256 amount) returns()
func (_IBank *IBankTransactor) DepositERC20(opts *bind.TransactOpts, token common.Address, sender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "depositERC20", token, sender, amount)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x1cad5a40.
//
// Solidity: function depositERC20(address token, address sender, uint256 amount) returns()
func (_IBank *IBankSession) DepositERC20(token common.Address, sender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.DepositERC20(&_IBank.TransactOpts, token, sender, amount)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0x1cad5a40.
//
// Solidity: function depositERC20(address token, address sender, uint256 amount) returns()
func (_IBank *IBankTransactorSession) DepositERC20(token common.Address, sender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.DepositERC20(&_IBank.TransactOpts, token, sender, amount)
}

// DepositNFT is a paid mutator transaction binding the contract method 0x548bfbaa.
//
// Solidity: function depositNFT(address token, address sender, uint256 tokenid) returns()
func (_IBank *IBankTransactor) DepositNFT(opts *bind.TransactOpts, token common.Address, sender common.Address, tokenid *big.Int) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "depositNFT", token, sender, tokenid)
}

// DepositNFT is a paid mutator transaction binding the contract method 0x548bfbaa.
//
// Solidity: function depositNFT(address token, address sender, uint256 tokenid) returns()
func (_IBank *IBankSession) DepositNFT(token common.Address, sender common.Address, tokenid *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.DepositNFT(&_IBank.TransactOpts, token, sender, tokenid)
}

// DepositNFT is a paid mutator transaction binding the contract method 0x548bfbaa.
//
// Solidity: function depositNFT(address token, address sender, uint256 tokenid) returns()
func (_IBank *IBankTransactorSession) DepositNFT(token common.Address, sender common.Address, tokenid *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.DepositNFT(&_IBank.TransactOpts, token, sender, tokenid)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address payee, address recipient, uint256 amount) returns()
func (_IBank *IBankTransactor) Withdraw(opts *bind.TransactOpts, payee common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "withdraw", payee, recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address payee, address recipient, uint256 amount) returns()
func (_IBank *IBankSession) Withdraw(payee common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.Withdraw(&_IBank.TransactOpts, payee, recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address payee, address recipient, uint256 amount) returns()
func (_IBank *IBankTransactorSession) Withdraw(payee common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.Withdraw(&_IBank.TransactOpts, payee, recipient, amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0x44004cc1.
//
// Solidity: function withdrawERC20(address token, address to, uint256 amount) returns()
func (_IBank *IBankTransactor) WithdrawERC20(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "withdrawERC20", token, to, amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0x44004cc1.
//
// Solidity: function withdrawERC20(address token, address to, uint256 amount) returns()
func (_IBank *IBankSession) WithdrawERC20(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.WithdrawERC20(&_IBank.TransactOpts, token, to, amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0x44004cc1.
//
// Solidity: function withdrawERC20(address token, address to, uint256 amount) returns()
func (_IBank *IBankTransactorSession) WithdrawERC20(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.WithdrawERC20(&_IBank.TransactOpts, token, to, amount)
}

// WithdrawNFT is a paid mutator transaction binding the contract method 0xbff863ec.
//
// Solidity: function withdrawNFT(address token, address to, uint256 tokenid) returns()
func (_IBank *IBankTransactor) WithdrawNFT(opts *bind.TransactOpts, token common.Address, to common.Address, tokenid *big.Int) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "withdrawNFT", token, to, tokenid)
}

// WithdrawNFT is a paid mutator transaction binding the contract method 0xbff863ec.
//
// Solidity: function withdrawNFT(address token, address to, uint256 tokenid) returns()
func (_IBank *IBankSession) WithdrawNFT(token common.Address, to common.Address, tokenid *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.WithdrawNFT(&_IBank.TransactOpts, token, to, tokenid)
}

// WithdrawNFT is a paid mutator transaction binding the contract method 0xbff863ec.
//
// Solidity: function withdrawNFT(address token, address to, uint256 tokenid) returns()
func (_IBank *IBankTransactorSession) WithdrawNFT(token common.Address, to common.Address, tokenid *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.WithdrawNFT(&_IBank.TransactOpts, token, to, tokenid)
}

// IBankDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the IBank contract.
type IBankDepositedIterator struct {
	Event *IBankDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IBankDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IBankDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IBankDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankDeposited represents a Deposited event raised by the IBank contract.
type IBankDeposited struct {
	Id        *big.Int
	Payee     common.Address
	WeiAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x1599c0fcf897af5babc2bfcf707f5dc050f841b044d97c3251ecec35b9abf80b.
//
// Solidity: event Deposited(uint256 indexed id, address payee, uint256 weiAmount)
func (_IBank *IBankFilterer) FilterDeposited(opts *bind.FilterOpts, id []*big.Int) (*IBankDepositedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "Deposited", idRule)
	if err != nil {
		return nil, err
	}
	return &IBankDepositedIterator{contract: _IBank.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x1599c0fcf897af5babc2bfcf707f5dc050f841b044d97c3251ecec35b9abf80b.
//
// Solidity: event Deposited(uint256 indexed id, address payee, uint256 weiAmount)
func (_IBank *IBankFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *IBankDeposited, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "Deposited", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankDeposited)
				if err := _IBank.contract.UnpackLog(event, "Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposited is a log parse operation binding the contract event 0x1599c0fcf897af5babc2bfcf707f5dc050f841b044d97c3251ecec35b9abf80b.
//
// Solidity: event Deposited(uint256 indexed id, address payee, uint256 weiAmount)
func (_IBank *IBankFilterer) ParseDeposited(log types.Log) (*IBankDeposited, error) {
	event := new(IBankDeposited)
	if err := _IBank.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBankERC20DepositedIterator is returned from FilterERC20Deposited and is used to iterate over the raw logs and unpacked data for ERC20Deposited events raised by the IBank contract.
type IBankERC20DepositedIterator struct {
	Event *IBankERC20Deposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IBankERC20DepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankERC20Deposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IBankERC20Deposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IBankERC20DepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankERC20DepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankERC20Deposited represents a ERC20Deposited event raised by the IBank contract.
type IBankERC20Deposited struct {
	Id     *big.Int
	Token  common.Address
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterERC20Deposited is a free log retrieval operation binding the contract event 0x54b55b0dca2281471c1e9b3db0a1397b951657e8c4c8551f4b26316d1408097d.
//
// Solidity: event ERC20Deposited(uint256 indexed id, address indexed token, address sender, uint256 amount)
func (_IBank *IBankFilterer) FilterERC20Deposited(opts *bind.FilterOpts, id []*big.Int, token []common.Address) (*IBankERC20DepositedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "ERC20Deposited", idRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &IBankERC20DepositedIterator{contract: _IBank.contract, event: "ERC20Deposited", logs: logs, sub: sub}, nil
}

// WatchERC20Deposited is a free log subscription operation binding the contract event 0x54b55b0dca2281471c1e9b3db0a1397b951657e8c4c8551f4b26316d1408097d.
//
// Solidity: event ERC20Deposited(uint256 indexed id, address indexed token, address sender, uint256 amount)
func (_IBank *IBankFilterer) WatchERC20Deposited(opts *bind.WatchOpts, sink chan<- *IBankERC20Deposited, id []*big.Int, token []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "ERC20Deposited", idRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankERC20Deposited)
				if err := _IBank.contract.UnpackLog(event, "ERC20Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseERC20Deposited is a log parse operation binding the contract event 0x54b55b0dca2281471c1e9b3db0a1397b951657e8c4c8551f4b26316d1408097d.
//
// Solidity: event ERC20Deposited(uint256 indexed id, address indexed token, address sender, uint256 amount)
func (_IBank *IBankFilterer) ParseERC20Deposited(log types.Log) (*IBankERC20Deposited, error) {
	event := new(IBankERC20Deposited)
	if err := _IBank.contract.UnpackLog(event, "ERC20Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBankERC20WithdrawnIterator is returned from FilterERC20Withdrawn and is used to iterate over the raw logs and unpacked data for ERC20Withdrawn events raised by the IBank contract.
type IBankERC20WithdrawnIterator struct {
	Event *IBankERC20Withdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IBankERC20WithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankERC20Withdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IBankERC20Withdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IBankERC20WithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankERC20WithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankERC20Withdrawn represents a ERC20Withdrawn event raised by the IBank contract.
type IBankERC20Withdrawn struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterERC20Withdrawn is a free log retrieval operation binding the contract event 0xbfed55bdcd242e3dd0f60ddd7d1e87c67f61c34cd9527b3e6455d841b1025362.
//
// Solidity: event ERC20Withdrawn(address indexed token, address to, uint256 amount)
func (_IBank *IBankFilterer) FilterERC20Withdrawn(opts *bind.FilterOpts, token []common.Address) (*IBankERC20WithdrawnIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "ERC20Withdrawn", tokenRule)
	if err != nil {
		return nil, err
	}
	return &IBankERC20WithdrawnIterator{contract: _IBank.contract, event: "ERC20Withdrawn", logs: logs, sub: sub}, nil
}

// WatchERC20Withdrawn is a free log subscription operation binding the contract event 0xbfed55bdcd242e3dd0f60ddd7d1e87c67f61c34cd9527b3e6455d841b1025362.
//
// Solidity: event ERC20Withdrawn(address indexed token, address to, uint256 amount)
func (_IBank *IBankFilterer) WatchERC20Withdrawn(opts *bind.WatchOpts, sink chan<- *IBankERC20Withdrawn, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "ERC20Withdrawn", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankERC20Withdrawn)
				if err := _IBank.contract.UnpackLog(event, "ERC20Withdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseERC20Withdrawn is a log parse operation binding the contract event 0xbfed55bdcd242e3dd0f60ddd7d1e87c67f61c34cd9527b3e6455d841b1025362.
//
// Solidity: event ERC20Withdrawn(address indexed token, address to, uint256 amount)
func (_IBank *IBankFilterer) ParseERC20Withdrawn(log types.Log) (*IBankERC20Withdrawn, error) {
	event := new(IBankERC20Withdrawn)
	if err := _IBank.contract.UnpackLog(event, "ERC20Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBankNFTDepositedIterator is returned from FilterNFTDeposited and is used to iterate over the raw logs and unpacked data for NFTDeposited events raised by the IBank contract.
type IBankNFTDepositedIterator struct {
	Event *IBankNFTDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IBankNFTDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankNFTDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IBankNFTDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IBankNFTDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankNFTDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankNFTDeposited represents a NFTDeposited event raised by the IBank contract.
type IBankNFTDeposited struct {
	Id      *big.Int
	Token   common.Address
	Sender  common.Address
	Tokenid *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNFTDeposited is a free log retrieval operation binding the contract event 0x327b4964b21908843dad3036ddb4742ba7dac55409b0f9a23772775f7fe99afd.
//
// Solidity: event NFTDeposited(uint256 indexed id, address indexed token, address sender, uint256 tokenid)
func (_IBank *IBankFilterer) FilterNFTDeposited(opts *bind.FilterOpts, id []*big.Int, token []common.Address) (*IBankNFTDepositedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "NFTDeposited", idRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &IBankNFTDepositedIterator{contract: _IBank.contract, event: "NFTDeposited", logs: logs, sub: sub}, nil
}

// WatchNFTDeposited is a free log subscription operation binding the contract event 0x327b4964b21908843dad3036ddb4742ba7dac55409b0f9a23772775f7fe99afd.
//
// Solidity: event NFTDeposited(uint256 indexed id, address indexed token, address sender, uint256 tokenid)
func (_IBank *IBankFilterer) WatchNFTDeposited(opts *bind.WatchOpts, sink chan<- *IBankNFTDeposited, id []*big.Int, token []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "NFTDeposited", idRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankNFTDeposited)
				if err := _IBank.contract.UnpackLog(event, "NFTDeposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNFTDeposited is a log parse operation binding the contract event 0x327b4964b21908843dad3036ddb4742ba7dac55409b0f9a23772775f7fe99afd.
//
// Solidity: event NFTDeposited(uint256 indexed id, address indexed token, address sender, uint256 tokenid)
func (_IBank *IBankFilterer) ParseNFTDeposited(log types.Log) (*IBankNFTDeposited, error) {
	event := new(IBankNFTDeposited)
	if err := _IBank.contract.UnpackLog(event, "NFTDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBankNFTWithdrawnIterator is returned from FilterNFTWithdrawn and is used to iterate over the raw logs and unpacked data for NFTWithdrawn events raised by the IBank contract.
type IBankNFTWithdrawnIterator struct {
	Event *IBankNFTWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IBankNFTWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankNFTWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IBankNFTWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IBankNFTWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankNFTWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankNFTWithdrawn represents a NFTWithdrawn event raised by the IBank contract.
type IBankNFTWithdrawn struct {
	Token   common.Address
	To      common.Address
	Tokenid *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNFTWithdrawn is a free log retrieval operation binding the contract event 0xbbde41973f9ce4890f7ad9762c23d8191f261fd643bdf13ed8bbc10549b49fcb.
//
// Solidity: event NFTWithdrawn(address indexed token, address to, uint256 tokenid)
func (_IBank *IBankFilterer) FilterNFTWithdrawn(opts *bind.FilterOpts, token []common.Address) (*IBankNFTWithdrawnIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "NFTWithdrawn", tokenRule)
	if err != nil {
		return nil, err
	}
	return &IBankNFTWithdrawnIterator{contract: _IBank.contract, event: "NFTWithdrawn", logs: logs, sub: sub}, nil
}

// WatchNFTWithdrawn is a free log subscription operation binding the contract event 0xbbde41973f9ce4890f7ad9762c23d8191f261fd643bdf13ed8bbc10549b49fcb.
//
// Solidity: event NFTWithdrawn(address indexed token, address to, uint256 tokenid)
func (_IBank *IBankFilterer) WatchNFTWithdrawn(opts *bind.WatchOpts, sink chan<- *IBankNFTWithdrawn, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "NFTWithdrawn", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankNFTWithdrawn)
				if err := _IBank.contract.UnpackLog(event, "NFTWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNFTWithdrawn is a log parse operation binding the contract event 0xbbde41973f9ce4890f7ad9762c23d8191f261fd643bdf13ed8bbc10549b49fcb.
//
// Solidity: event NFTWithdrawn(address indexed token, address to, uint256 tokenid)
func (_IBank *IBankFilterer) ParseNFTWithdrawn(log types.Log) (*IBankNFTWithdrawn, error) {
	event := new(IBankNFTWithdrawn)
	if err := _IBank.contract.UnpackLog(event, "NFTWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBankWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the IBank contract.
type IBankWithdrawnIterator struct {
	Event *IBankWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IBankWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IBankWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IBankWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankWithdrawn represents a Withdrawn event raised by the IBank contract.
type IBankWithdrawn struct {
	Payee     common.Address
	Recipient common.Address
	WeiAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed payee, address recipient, uint256 weiAmount)
func (_IBank *IBankFilterer) FilterWithdrawn(opts *bind.FilterOpts, payee []common.Address) (*IBankWithdrawnIterator, error) {

	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "Withdrawn", payeeRule)
	if err != nil {
		return nil, err
	}
	return &IBankWithdrawnIterator{contract: _IBank.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed payee, address recipient, uint256 weiAmount)
func (_IBank *IBankFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *IBankWithdrawn, payee []common.Address) (event.Subscription, error) {

	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "Withdrawn", payeeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankWithdrawn)
				if err := _IBank.contract.UnpackLog(event, "Withdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawn is a log parse operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed payee, address recipient, uint256 weiAmount)
func (_IBank *IBankFilterer) ParseWithdrawn(log types.Log) (*IBankWithdrawn, error) {
	event := new(IBankWithdrawn)
	if err := _IBank.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
