// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "./Bank.sol";

/**
 * @dev bridge contract
 * By deppositing erc20 tokens or native token on this contract, equivalent amount of wrapped token is minted
 * The deposited token is locked, until the minted token is burnd.
 */
contract Bridge is Context, ReentrancyGuard, AccessControl {
    using SafeERC20 for IERC20;

    /* access permission of this contract */
    bytes32 public constant BRIDGE_ACCESS_ROLE =
        keccak256("BRIDGE_ACCESS_ROLE");
    /* hold deposited tokens. event if this contract is upgraded, bank contract keep holding deposited tokens  */
    Bank public bank;

    constructor(Bank _bank) {
        bank = _bank;

        // grant admin role to deployer
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setRoleAdmin(DEFAULT_ADMIN_ROLE, DEFAULT_ADMIN_ROLE);

        // grant access permission of this and bank contract to deployer
        _setupRole(BRIDGE_ACCESS_ROLE, _msgSender());
        _setRoleAdmin(BRIDGE_ACCESS_ROLE, DEFAULT_ADMIN_ROLE);
        _setupRole(bank.BANK_ACCESS_ROLE(), address(this));
        _setRoleAdmin(bank.BANK_ACCESS_ROLE(), DEFAULT_ADMIN_ROLE);
    }

    function depositsOf(address owner) public view returns (uint256) {
        return bank.depositsOf(owner);
    }

    function deposit() public payable {
        bank.deposit(_msgSender());
    }

    function withdraw(address payable payee, uint256 amount) public {
        bank.withdraw(payee, amount);
    }

    function depositERC20(IERC20 token, uint256 amount) public {
        bank.safeTransferFrom(token, _msgSender(), amount);
    }

    function withdrawERC20(
        IERC20 token,
        address to,
        uint256 amount
    ) public {
        bank.safeTransfer(token, to, amount);
    }

    function erc20DepositsOf(IERC20 token, address owner)
        public
        view
        returns (uint256)
    {
        return bank.erc20DepositsOf(token, owner);
    }
}
