// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
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
    /* index of `NFTDeposited` evnet */
    Counters.Counter public indexNFTDeposited;

    event Deposited(uint256 indexed id, address payee, uint256 weiAmount);
    event ERC20Deposited(
        uint256 indexed id,
        IERC20 indexed token,
        address sender,
        uint256 amount
    );
    event NFTDeposited(
        uint256 indexed id,
        IERC721 indexed token,
        address sender,
        uint256 tokenid
    );
    event Withdrawn(
        address indexed payee,
        address recipient,
        uint256 weiAmount
    );
    event ERC20Withdrawn(IERC20 indexed token, address to, uint256 amount);
    event NFTWithdrawn(IERC721 indexed token, address to, uint256 tokenid);

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

    /**
     * @dev stores the sent amount. equivalent amout of wrapped coin is minted on the bridging chain
     * @param payee The destination address of the funds.
     */
    function deposit(address payee) public payable {
        uint256 amount = msg.value;
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
        require(address(this).balance >= amount, "exceed deposited amount");
        recipient.sendValue(amount);
        emit Withdrawn(payee, recipient, amount);
    }

    //------------------ ERC20 Token ------------------//

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
     * @param to The address of the fund is transferred
     * @param amount Transfer amount
     */
    function withdrawERC20(
        IERC20 token,
        address to,
        uint256 amount
    ) public onlyPermited(BANK_ACCESS_ROLE) {
        token.safeTransfer(to, amount);
        emit ERC20Withdrawn(token, to, amount);
    }

    //------------------ NFT ------------------//

    /**
     * @dev call NFT `transferFrom`
     * @param token NFT token address
     * @param sender From address
     * @param tokenid the token id
     */
    function depositNFT(
        IERC721 token,
        address sender,
        uint256 tokenid
    ) public {
        token.transferFrom(sender, address(this), tokenid);
        emit NFTDeposited(
            assignIndex(indexNFTDeposited),
            token,
            sender,
            tokenid
        );
    }

    /**
     * @dev call NFT `safeTransfer`. authenticated contract only
     * @param token NFT token address
     * @param to The address of the fund is transferred
     * @param tokenid The token id
     */
    function withdrawNFT(
        IERC721 token,
        address to,
        uint256 tokenid
    ) public onlyPermited(BANK_ACCESS_ROLE) {
        token.safeTransferFrom(address(this), to, tokenid);
        emit NFTWithdrawn(token, to, tokenid);
    }

    //------------------ inner functions ------------------//

    function assignIndex(Counters.Counter storage counter)
        internal
        returns (uint256)
    {
        counter.increment();
        return counter.current() - 1;
    }
}
