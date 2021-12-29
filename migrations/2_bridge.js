const AccessController = artifacts.require("AccessController");
const Bank = artifacts.require("Bank");
const Bridge = artifacts.require("Bridge");

let admin = process.env.ADMIN_ACCOUNT

module.exports = function (deployer, network, accounts) {
  if (!process.env.hasOwnProperty("ADMIN_ACCOUNT")) {
    if (accounts.length == 0) {
      throw "please set `ADMIN_ACCOUNT` to the env variable"
    }
    admin = accounts[0]
  }

  deployer.deploy(AccessController).then(function() {
    return deployer.deploy(Bank, AccessController.address);
  }).then(function() {
    return deployer.deploy(Bridge, AccessController.address, Bank.address)
  }).then(async function() {
    // grant access roles
    const controller = await AccessController.at(AccessController.address);
    const bank = await Bank.at(Bank.address);
    const bridge = await Bridge.at(Bridge.address);

    const bankAccessRole = await bank.BANK_ACCESS_ROLE()
    const bridgeAccessRole = await bridge.BRIDGE_ACCESS_ROLE()

    await controller.setupRole(bankAccessRole, Bridge.address);
    let has = await controller.hasRole(bankAccessRole, Bridge.address);
    if (!has) {
      throw "faild to grant access role"
    }
    await controller.setupRole(bankAccessRole, admin);
    has = await controller.hasRole(bankAccessRole, admin);
    if (!has) {
      throw "faild to grant access role"
    }

    await controller.setupRole(bridgeAccessRole, admin);
    has = await controller.hasRole(bridgeAccessRole, admin);
    if (!has) {
      throw "faild to grant access role"
    }
  });
};
