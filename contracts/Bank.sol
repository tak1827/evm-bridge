// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "./AccessControlRegistry.sol";

/**
 * @dev Bank contract
 * hold deposited coin or erc20 tokens
 * even when exchange contract is upgraded, intended not upgrade to sustain token approved status
 */
contract Bank is Context, ReentrancyGuard, AccessControlRegistry {
    using Address for address payable;
    using SafeERC20 for IERC20;
    using Counters for Counters.Counter;

    /* access permission of this contract. the exchange contract is assume to be allocated this role */
    bytes32 public constant BANK_ACCESS_ROLE = keccak256("BANK_ACCESS_ROLE");
    /* index of `Deposited` evnet */
    Counters.Counter public indexDeposited;
    /* index of `ERC20Deposited` evnet */
    Counters.Counter public indexERC20Deposited;
    /* eth deposit amount map */
    mapping(address => uint256) private _deposits;
    /* erc20 token deposit amount map */
    /* mapping(erc20 address => mapping(owner addresss => amount)) */
    mapping(address => mapping(address => uint256)) private _erc20Deposits;

    event Deposited(uint256 indexed index, address payee, uint256 weiAmount);
    event ERC20Deposited(
        uint256 indexed index,
        IERC20 indexed token,
        address sender,
        uint256 amount
    );
    event Withdrawn(
        address indexed payee,
        address recipient,
        uint256 weiAmount
    );
    event ERC20Withdrawn(
        IERC20 indexed token,
        address indexed from,
        address to,
        uint256 amount
    );

    /**
     * @dev initializes the contract by setting a `transferAgent` and `accessControler`
     */
    constructor(AccessControl control) AccessControlRegistry(control) {}

    /**
     * @dev set new AccessControler address. Authenticated contract or person only
     * @param newaccessControler AccessControler address
     */
    function setAccessControler(AccessControl newaccessControler)
        public
        onlyPermited(BANK_ACCESS_ROLE)
    {
        _setAccessControler(newaccessControler);
    }

    //------------------ Native Coin ------------------//

    function depositsOf(address owner) public view returns (uint256) {
        return _deposits[owner];
    }

    /**
     * @dev stores the sent amount. equivalent amout of wrapped coin is minted on the bridging chain
     * @param payee The destination address of the funds.
     */
    function deposit(address payee) public payable {
        uint256 amount = msg.value;
        _deposits[payee] += amount;
        emit Deposited(assignIndex(indexDeposited), payee, amount);
    }

    /**
     * @dev withdraw balance for a payee. equivalent amount of wraped coin should be burned on the bridging chain
     * @param payee The address whose funds will be withdrawn
     * @param recipient The address of transferred to
     * @param amount The amount of withdrawn coin
     */
    function withdraw(
        address payee,
        address payable recipient,
        uint256 amount
    ) public nonReentrant onlyPermited(BANK_ACCESS_ROLE) {
        require(_deposits[payee] >= amount, "exceed deposited amount");
        _deposits[payee] -= amount;
        recipient.sendValue(amount);
        emit Withdrawn(payee, recipient, amount);
    }

    //------------------ ERC20 Token ------------------//

    function erc20DepositsOf(IERC20 token, address owner)
        public
        view
        returns (uint256)
    {
        return _erc20Deposits[address(token)][owner];
    }

    /**
     * @dev call ERC20 `safeTransferFrom`
     * @param token ERC20 token address
     * @param sender From address
     * @param amount Transfer amount
     */
    function depositERC20(
        IERC20 token,
        address sender,
        uint256 amount
    ) public {
        token.safeTransferFrom(sender, address(this), amount);
        _erc20Deposits[address(token)][sender] += amount;
        emit ERC20Deposited(
            assignIndex(indexERC20Deposited),
            token,
            sender,
            amount
        );
    }

    /**
     * @dev call ERC20 `safeTransfer`. authenticated contract only
     * @param token ERC20 token address
     * @param from To address of the fund is withdrawn
     * @param to The address of the fund is transferred
     * @param amount Transfer amount
     */
    function withdrawERC20(
        IERC20 token,
        address from,
        address to,
        uint256 amount
    ) public onlyPermited(BANK_ACCESS_ROLE) {
        require(
            _erc20Deposits[address(token)][from] >= amount,
            "exceed deposited amount"
        );
        _erc20Deposits[address(token)][from] -= amount;
        token.safeTransfer(to, amount);
        emit ERC20Withdrawn(token, from, to, amount);
    }

    function assignIndex(Counters.Counter storage counter)
        internal
        returns (uint256)
    {
        counter.increment();
        return counter.current() - 1;
    }
}
