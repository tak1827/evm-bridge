// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "./AccessController.sol";
import "./AccessControlRegistry.sol";
import "./Bank.sol";

/**
 * @dev bridge contract
 * By deppositing erc20 tokens or native token on this contract, equivalent amount of wrapped token is minted
 * The deposited token is locked, until the minted token is burnd.
 */
contract Bridge is Context, AccessControlRegistry {
    using SafeERC20 for IERC20;

    /* access permission of this contract */
    bytes32 public constant BRIDGE_ACCESS_ROLE =
        keccak256("BRIDGE_ACCESS_ROLE");
    /* hold deposited tokens. event if this contract is upgraded, bank contract keep holding deposited tokens  */
    Bank public bank;

    constructor(AccessControl control, Bank _bank)
        AccessControlRegistry(control)
    {
        bank = _bank;
    }

    /**
     * @dev set new AccessControler address. Authenticated contract or person only
     * @param newaccessControler AccessControler address
     */
    function setAccessControler(AccessControl newaccessControler)
        public
        onlyPermited(BRIDGE_ACCESS_ROLE)
    {
        _setAccessControler(newaccessControler);
    }

    function depositsOf(address owner) public view returns (uint256) {
        return bank.depositsOf(owner);
    }

    function deposit() public payable {
        bank.deposit{value: msg.value}(_msgSender());
    }

    function withdraw(
        address payee,
        address payable recepient,
        uint256 amount
    ) public onlyPermited(BRIDGE_ACCESS_ROLE) {
        bank.withdraw(payee, recepient, amount);
    }

    function depositERC20(IERC20 token, uint256 amount) public {
        bank.depositERC20(token, _msgSender(), amount);
    }

    function withdrawERC20(
        IERC20 token,
        address from,
        address to,
        uint256 amount
    ) public onlyPermited(BRIDGE_ACCESS_ROLE) {
        bank.withdrawERC20(token, from, to, amount);
    }

    function erc20DepositsOf(IERC20 token, address owner)
        public
        view
        returns (uint256)
    {
        return bank.erc20DepositsOf(token, owner);
    }
}
