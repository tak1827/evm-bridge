// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "./AccessController.sol";

/**
 * @dev Bank contract
 * hold deposited coin or erc20 tokens
 * even when exchange contract is upgraded, intended not upgrade to sustain token approved status
 */
contract Bank is Context, ReentrancyGuard, AccessController {
    using Address for address payable;
    using SafeERC20 for IERC20;

    /* access permission of this contract. the exchange contract is assume to be allocated this role */
    bytes32 public constant BANK_ACCESS_ROLE = keccak256("BANK_ACCESS_ROLE");
    /* eth deposit amount map */
    mapping(address => uint256) private _deposits;
    /* erc20 token deposit amount map */
    /* mapping(erc20 address => mapping(owner addresss => amount)) */
    mapping(address => mapping(address => uint256)) private _erc20Deposits;

    event Deposited(address indexed payee, uint256 weiAmount);
    event Withdrawn(address indexed payee, uint256 weiAmount);

    event ERC20Deposited(IERC20 indexed token, address indexed sender, uint256 weiAmount);
    event ERC20Withdrawn(IERC20 indexed token, address indexed to, uint256 weiAmount);

    /**
     * @dev initializes the contract by setting a `transferAgent` and `accessControler`
     */
    constructor(AccessControl _accessControler) AccessController(_accessControler) {}


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
        emit Deposited(payee, amount);
    }

    /**
     * @dev withdraw balance for a payee. equivalent amount of wraped coin should be burned on the bridging chain
     * @param payee The address whose funds will be withdrawn and transferred to.
     * @param amount The amount of withdrawn coin
     */
    function withdraw(address payable payee, uint256 amount) public nonReentrant onlyPermited(BANK_ACCESS_ROLE) {
        require(_deposits[payee] >=  amount, "exceed deposited amout");
        _deposits[payee] -= amount;
        payee.sendValue(amount);
        emit Withdrawn(payee, amount);
    }

    //------------------ ERC20 Token ------------------//

    function erc20DepositsOf(IERC20 token, address owner) public view returns (uint256) {
        return _erc20Deposits[address(token)][owner];
    }

    /**
     * @dev call ERC20 `safeTransferFrom`
     * @param token ERC20 token address
     * @param sender From address
     * @param amount Transfer amount
     */
    function safeTransferFrom(
        IERC20 token,
        address sender,
        uint256 amount
    ) public {
        token.safeTransferFrom(sender, address(this), amount);
        _erc20Deposits[address(token)][sender] += amount;
        emit ERC20Deposited(token, sender, amount);
    }

    /**
     * @dev call ERC20 `safeTransfer`. authenticated contract only
     * @param token ERC20 token address
     * @param to To address
     * @param amount Transfer amount
     */
    function safeTransfer(
        IERC20 token,
        address to,
        uint256 amount
    ) public onlyPermited(BANK_ACCESS_ROLE) {
        require(_erc20Deposits[address(token)][to] >=  amount, "exceed deposited amout");
        token.safeTransfer(to, amount);
        _erc20Deposits[address(token)][to] -= amount;
        emit ERC20Withdrawn(token, to, amount);
    }
}
