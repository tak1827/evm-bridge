// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @dev exchange NFT saving gas cost.
 * a sell or a buyer only pay gas cost in the last step when the nft is transferd.
 * an exchange operator earn protocol fee, paied by seller or buyer.
 */
contract AccessController {
    /* access permission holder */
    AccessControl public accessControler;

    /**
     * @dev checks access permission
     */
    modifier onlyPermited(bytes32 role) {
        require(
            accessControler.hasRole(role, address(this)),
            "no access permission"
        );
        _;
    }

    /**
     * @dev initializes the contract by setting an `accessControler`
     */
    constructor(AccessControl _accessControler) {
        accessControler = _accessControler;
    }

    /**
     * @dev set new AccessControler address. Authenticated contract or person only
     * @param newaccessControler AccessControler address
     */
    function _setAccessControler(AccessControl newaccessControler)
        internal
        virtual
    {
        accessControler = newaccessControler;
    }
}
