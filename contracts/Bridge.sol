// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./AccessController.sol";
import "./AccessControlRegistry.sol";
import "./Bank.sol";

/**
 * @dev bridge contract
 * By deppositing erc20 tokens or native token on this contract, equivalent amount of wrapped token is minted
 * The deposited token is locked, until the minted token is burnd.
 */
contract Bridge is Context, AccessControlRegistry {
    using EnumerableSet for EnumerableSet.AddressSet;

    /* access permission of this contract */
    bytes32 public constant BRIDGE_ACCESS_ROLE =
        keccak256("BRIDGE_ACCESS_ROLE");
    /* hold deposited tokens. event if this contract is upgraded, bank contract keep holding deposited tokens  */
    Bank public bank;
    /* whitelisted erc20 token addresses */
    EnumerableSet.AddressSet private erc20Whitelist;
    /* whitelisted nft addresses */
    EnumerableSet.AddressSet private nftWhitelist;

    /**
     * @dev checks whitelist
     */
    modifier onlyWhitelistedERC20(address erc20) {
        require(erc20Whitelist.contains(erc20), "not whitelisted");
        _;
    }

    /**
     * @dev checks whitelist
     */
    modifier onlyWhitelistedNFT(address nft) {
        require(nftWhitelist.contains(nft), "not whitelisted");
        _;
    }

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

    function depositERC20(IERC20 token, uint256 amount)
        public
        onlyWhitelistedERC20(address(token))
    {
        bank.depositERC20(token, _msgSender(), amount);
    }

    function withdrawERC20(
        IERC20 token,
        address to,
        uint256 amount
    ) public onlyPermited(BRIDGE_ACCESS_ROLE) {
        bank.withdrawERC20(token, to, amount);
    }

    function getERC20Whitelist(uint256 index) public view returns (address) {
        return erc20Whitelist.at(index);
    }

    function countERC20Whitelist() public view returns (uint256) {
        return erc20Whitelist.length();
    }

    function addERC20Whitelist(address erc20)
        public
        onlyPermited(BRIDGE_ACCESS_ROLE)
    {
        erc20Whitelist.add(erc20);
    }

    function removeERC20Whitelist(address erc20)
        public
        onlyPermited(BRIDGE_ACCESS_ROLE)
    {
        erc20Whitelist.remove(erc20);
    }

    function depositNFT(IERC721 token, uint256 tokenid)
        public
        onlyWhitelistedNFT(address(token))
    {
        bank.depositNFT(token, _msgSender(), tokenid);
    }

    function withdrawNFT(
        IERC721 token,
        address to,
        uint256 tokenid
    ) public onlyPermited(BRIDGE_ACCESS_ROLE) {
        bank.withdrawNFT(token, to, tokenid);
    }

    function getNFTWhitelist(uint256 index) public view returns (address) {
        return nftWhitelist.at(index);
    }

    function countNFTWhitelist() public view returns (uint256) {
        return nftWhitelist.length();
    }

    function addNFTWhitelist(address nft)
        public
        onlyPermited(BRIDGE_ACCESS_ROLE)
    {
        nftWhitelist.add(nft);
    }

    function removeNFTWhitelist(address nft)
        public
        onlyPermited(BRIDGE_ACCESS_ROLE)
    {
        nftWhitelist.remove(nft);
    }
}
