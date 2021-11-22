// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Context.sol";

/**
 * @dev Registry of AccessController
 */
contract AccessControlRegistry is Context {
    /* current access controller version */
    uint16 public controlVersion;
    /* mapping of access controller */
    mapping(uint16 => AccessControl) public accessController;

    /**
     * @dev Emitted when access controller updated
     */
    event AccessControlUpdated(
        uint16 indexed version,
        AccessControl indexed accessController
    );

    /**
     * @dev checks access permission
     */
    modifier onlyPermited(bytes32 role) {
        require(
            accessController[controlVersion].hasRole(role, _msgSender()),
            "no access permission"
        );
        _;
    }

    /**
     * @dev initializes the contract by setting an `accessControler`
     */
    constructor(AccessControl _accessControler) {
        accessController[controlVersion] = _accessControler;
    }

    /**
     * @dev set new AccessControler address. Authenticated contract or person only
     * @param newaccessControler AccessControler address
     */
    function _setAccessControler(AccessControl newaccessControler)
        internal
        virtual
    {
        controlVersion++;
        accessController[controlVersion] = newaccessControler;
        emit AccessControlUpdated(controlVersion, newaccessControler);
    }
}
