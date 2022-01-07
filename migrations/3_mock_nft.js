const NFT = artifacts.require("MockNFT");

module.exports = function (deployer, network, accounts) {
  deployer.deploy(NFT).then(async function() {
    const nft = await NFT.at(NFT.address);
    const size = 10
    for (let i = 0; i < accounts.length && i < 1; i++) {
      for (let j = i * size; j < (i+1)*size; j++ ) {
        await nft.safeMint(j, accounts[i], "")
      }
    }
  });
};
