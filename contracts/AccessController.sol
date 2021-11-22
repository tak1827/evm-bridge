pragma solidity 0.8.9;

import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @dev impl of AccessControl
 */
contract AccessController is AccessControl {
    constructor() {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _setRoleAdmin(DEFAULT_ADMIN_ROLE, DEFAULT_ADMIN_ROLE);
    }

    function setupRole(bytes32 role, address account) public {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, _msgSender()) ||
                hasRole(role, _msgSender())
        );
        _setupRole(role, account);
    }

    function setRoleAdmin(bytes32 roleId, bytes32 adminRoleId)
        public
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        _setRoleAdmin(roleId, adminRoleId);
    }
}
