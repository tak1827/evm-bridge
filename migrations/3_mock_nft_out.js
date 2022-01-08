const NFT = artifacts.require("MockNFT");

module.exports = function (deployer, network, accounts) {
  deployer.deploy(NFT)
};
