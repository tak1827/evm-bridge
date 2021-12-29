const ERC20 = artifacts.require("WrappedToken");
const Bank = artifacts.require("Bank");
const Bridge = artifacts.require("Bridge");

let admin = process.env.ADMIN_ACCOUNT

const NAME = "DRV_TOKEN"
const SYMBOL = "DTN"
const INT_BALANCE = "1000000000000000000000"

module.exports = function (deployer, network, accounts) {
  if (!process.env.hasOwnProperty("ADMIN_ACCOUNT")) {
    if (accounts.length == 0) {
      throw "please set `ADMIN_ACCOUNT` to the env variable"
    }
    admin = accounts[0]
  }

  deployer.deploy(ERC20, NAME, SYMBOL, admin, INT_BALANCE).then(async function() {
    const erc20 = await ERC20.at(ERC20.address);
    for (let i = 1; i < accounts.length; i++) {
      await erc20.mint(accounts[i], INT_BALANCE)
    }
  });
};
