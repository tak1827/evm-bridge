const AccessController = artifacts.require("AccessController");
const Bank = artifacts.require("Bank");
const Bridge = artifacts.require("Bridge");

module.exports = function (deployer) {
  deployer.deploy(AccessController).then(function() {
    return deployer.deploy(Bank, AccessController.address);
  }).then(function() {
    return deployer.deploy(Bridge, AccessController.address, Bank.address)
  });
};
