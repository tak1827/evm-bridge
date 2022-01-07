pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";

contract MockNFT is ERC721URIStorage {
    string public constant TOKEN_NAME = "Crypto Fantacy Token";

    string public constant TOKEN_SYMBOL = "CFT";

    constructor() ERC721(TOKEN_NAME, TOKEN_SYMBOL) {}

    function exists(uint256 tokenId) public view returns (bool) {
        return _exists(tokenId);
    }

    function safeMint(
        uint256 tokenId,
        address to,
        string memory uri
    ) public {
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
    }

    function burn(uint256 tokenId) public {
        _burn(tokenId);
    }
}
