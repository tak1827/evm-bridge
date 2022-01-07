const AccessControlRegistry = artifacts.require("AccessControlRegistry");
const AccessController = artifacts.require("AccessController");
const Bank = artifacts.require("Bank");
const Bridge = artifacts.require("Bridge");
const ERC20 = artifacts.require("WrappedToken");
const NFT = artifacts.require("MockNFT");

const { constants, BN, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { ZERO_ADDRESS } = constants;
const { expect } = require('chai');

const TokenID1 = 100
const TokenID2 = 101
const TokenID3 = 102

contract("Bridge", ([deployer, user1, user2, user3, attacker]) => {
  let controller;
  let bank;
  let bridge;
  let erc20;
  let nft;

  beforeEach(async () => {
    controller = await AccessController.new({from: deployer});
    const adminRole = await controller.DEFAULT_ADMIN_ROLE();

    bank = await Bank.new(controller.address, {from: deployer, gas: 5500000});
    const bankRole = await bank.BANK_ACCESS_ROLE();
    await controller.setupRole(bankRole, deployer, {from: deployer});
    await controller.setRoleAdmin(bankRole, adminRole, {from: deployer});

    bridge = await Bridge.new(controller.address, bank.address, {from: deployer, gas: 5500000});
    const bridgeRole = await bridge.BRIDGE_ACCESS_ROLE();
    await controller.setupRole(bridgeRole, deployer, {from: deployer});
    await controller.setRoleAdmin(bridgeRole, adminRole, {from: deployer});
    await controller.setupRole(bankRole, bridge.address, {from: deployer});

    erc20 = await ERC20.new("NAME", "SYM", deployer, 100000000, {from: deployer});
    await erc20.transfer(user1, 10000, {from: deployer});
    await erc20.transfer(user2, 10000, {from: deployer});
    await erc20.approve(bank.address, 1000, {from: user1});
    await erc20.approve(bank.address, 1000, {from: user2});

    nft = await NFT.new({from: deployer});
    await nft.safeMint(TokenID1, user1, {from: deployer});
    await nft.safeMint(TokenID2, user2, {from: deployer});
    await nft.safeMint(TokenID3, user3, {from: deployer});
    await nft.approve(bank.address, TokenID1, {from: user1});
    await nft.approve(bank.address, TokenID2, {from: user2});
  });

  describe('deploy', () => {
    it('check paramaters', async () => {
      const bankRole = await bank.BANK_ACCESS_ROLE();
      expect(await controller.hasRole(bankRole, bridge.address)).to.be.equal(true);
      expect(await controller.hasRole(bankRole, deployer)).to.be.equal(true);
      const bridgeRole = await bridge.BRIDGE_ACCESS_ROLE();
      expect(await controller.hasRole(bridgeRole, deployer)).to.be.equal(true);

      expect(await bridge.controlVersion()).to.be.bignumber.equal('0');
    });
  });

  describe('setAccessControler', () => {
    it('succedd', async () => {
      const newController = await AccessController.new({from: deployer});
      const receipt = await bridge.setAccessControler(newController.address, {from: deployer});

      expectEvent(receipt, 'AccessControlUpdated', {
        version: new BN(1),
        accessController: newController.address
      });

      expect(await bridge.controlVersion()).to.be.bignumber.equal('1');
      expect(await bridge.accessController(new BN(0))).to.be.equal(controller.address);
      expect(await bridge.accessController(new BN(1))).to.be.equal(newController.address);
    });

    it('fail by no autheticated', async function () {
      const newController = await AccessController.new({from: deployer});

      await expectRevert(
        bridge.setAccessControler(newController.address, {from: attacker}),
        "no access permission"
      );
    });
  });

  describe('deposit', () => {
    it('succedd', async () => {
      const amount = 1000;
      await bridge.deposit({from: user1, value: amount});

      expect(await web3.eth.getBalance(bridge.address)).to.be.bignumber.equal('0');
      expect(await web3.eth.getBalance(bank.address)).to.be.bignumber.equal(amount.toString());

      await bridge.deposit({from: user2, value: amount});

      expect(await web3.eth.getBalance(bank.address)).to.be.bignumber.equal('2000');
    });
  });

  describe('withdraw', () => {
    it('succedd', async () => {
      const amount = 1000;
      await bridge.deposit({from: user1, value: amount});
      await bridge.deposit({from: user2, value: amount});
      await bridge.withdraw(user1, user1, amount, {from: deployer});

      expect(await web3.eth.getBalance(bank.address)).to.be.bignumber.equal('1000');

      await bridge.withdraw(user2, user3, amount, {from: deployer});

      expect(await web3.eth.getBalance(bank.address)).to.be.bignumber.equal('0');
    });

    it('failed by over limit', async () => {
      const amount = 1000;
      await bridge.deposit({from: user1, value: amount});

      await expectRevert(
        bridge.withdraw(user1, user1, amount + 1, {from: deployer}),
        "exceed deposited amount"
      );
    });

    it('failed by witdrowing by owner', async () => {
      const amount = 1000;
      await bridge.deposit({from: user1, value: amount});

      await expectRevert(
        bridge.withdraw(user1, user1, amount, {from: user1}),
        "no access permission"
      );
    });

    it('failed by directly calling bank', async () => {
      const amount = 1000;
      await bridge.deposit({from: user1, value: amount});

      await expectRevert(
        bank.withdraw(user1, user1, amount, {from: attacker}),
        "no access permission"
      );
    });
  });

  describe('ERC20Whitelist', () => {
    it('addERC20Whitelist', async () => {
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})

      expect(await bridge.countERC20Whitelist()).to.be.bignumber.equal('1');
      expect(await bridge.getERC20Whitelist(0)).to.be.bignumber.equal(erc20.address);
    });

    it('fremoveERC20Whitelist', async () => {
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})
      await bridge.removeERC20Whitelist(erc20.address, {from: deployer})

      expect(await bridge.countERC20Whitelist()).to.be.bignumber.equal('0');
    });
  });

  describe('depositERC20', () => {
    it('succedd', async () => {
      const amount = 1000;
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})
      await bridge.depositERC20(erc20.address, amount, {from: user1});

      expect(await erc20.balanceOf(bridge.address)).to.be.bignumber.equal('0');
      expect(await erc20.balanceOf(bank.address)).to.be.bignumber.equal(amount.toString());

      await bridge.depositERC20(erc20.address, amount, {from: user2});

      expect(await erc20.balanceOf(bank.address)).to.be.bignumber.equal('2000');
    });

    it('failed by not whitelisted', async () => {
      const amount = 1000;
      await expectRevert(
        bridge.depositERC20(erc20.address, amount, {from: user1}),
        "not whitelisted"
      );
    });

    it('failed by over allowance', async () => {
      const amount = 1000;
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})
      await expectRevert(
        bridge.depositERC20(erc20.address, amount + 1, {from: user1}),
        "ERC20: transfer amount exceeds allowance"
      );
    });
  });

  describe('withdrawERC20', () => {
    it('succedd', async () => {
      const amount = 1000;
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})
      await bridge.depositERC20(erc20.address, amount, {from: user1});
      await bridge.depositERC20(erc20.address, amount, {from: user2});
      await bridge.withdrawERC20(erc20.address, user1, amount, {from: deployer});

      expect(await erc20.balanceOf(bank.address)).to.be.bignumber.equal('1000');
      expect(await erc20.balanceOf(user1)).to.be.bignumber.equal('10000');

      await bridge.withdrawERC20(erc20.address, user3, amount, {from: deployer});

      expect(await erc20.balanceOf(bank.address)).to.be.bignumber.equal('0');
      expect(await erc20.balanceOf(user2)).to.be.bignumber.equal('9000');
      expect(await erc20.balanceOf(user3)).to.be.bignumber.equal('1000');
    });

    it('failed by over limit', async () => {
      const amount = 1000;
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})
      await bridge.depositERC20(erc20.address, amount, {from: user1});

      await expectRevert(
        bridge.withdrawERC20(erc20.address, user1, amount + 1, {from: deployer}),
        "ERC20: transfer amount exceeds balance"
      );
    });

    it('failed by witdrowing by owner', async () => {
      const amount = 1000;
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})
      await bridge.depositERC20(erc20.address, amount, {from: user1});

      await expectRevert(
        bridge.withdrawERC20(erc20.address, user1, amount, {from: user1}),
        "no access permission"
      );
    });

    it('failed by directly calling bank', async () => {
      const amount = 1000;
      await bridge.addERC20Whitelist(erc20.address, {from: deployer})
      await bridge.depositERC20(erc20.address, amount, {from: user1});

      await expectRevert(
        bank.withdrawERC20(erc20.address, user1, amount, {from: attacker}),
        "no access permission"
      );
    });
  });

  describe('NFTWhitelist', () => {
    it('addNFTWhitelist', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})

      expect(await bridge.countNFTWhitelist()).to.be.bignumber.equal('1');
      expect(await bridge.getNFTWhitelist(0)).to.be.bignumber.equal(nft.address);
    });

    it('fremoveNFTWhitelist', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})
      await bridge.removeNFTWhitelist(nft.address, {from: deployer})

      expect(await bridge.countNFTWhitelist()).to.be.bignumber.equal('0');
    });
  });

  describe('depositNFT', () => {
    it('succedd', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})
      await bridge.depositNFT(nft.address, TokenID1, {from: user1});

      expect(await nft.ownerOf(TokenID1)).to.be.equal(bank.address);

      await bridge.depositNFT(nft.address, TokenID2, {from: user2});

      expect(await nft.balanceOf(bank.address)).to.be.bignumber.equal('2');
    });

    it('failed by not whitelisted', async () => {
      await expectRevert(
        bridge.depositNFT(nft.address, TokenID1, {from: user1}),
        "not whitelisted"
      );
    });

    it('failed by not approved', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})
      await expectRevert(
        bridge.depositNFT(nft.address, TokenID3, {from: user3}),
        "ERC721: transfer caller is not owner nor approved."
      );
    });
  });

  describe('withdrawNFT', () => {
    it('succedd', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})
      await bridge.depositNFT(nft.address, TokenID1, {from: user1});
      await bridge.depositNFT(nft.address, TokenID2, {from: user2});
      await bridge.withdrawNFT(nft.address, user3, TokenID1, {from: deployer});

      expect(await nft.ownerOf(TokenID1)).to.be.equal(user3);

      await bridge.withdrawNFT(nft.address, user2, TokenID2, {from: deployer});

      expect(await nft.ownerOf(TokenID2)).to.be.equal(user2);
      expect(await nft.balanceOf(bank.address)).to.be.bignumber.equal('0');
    });

    it('failed by not owner', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})
      await expectRevert(
        bridge.withdrawNFT(nft.address, user1, TokenID2, {from: deployer}),
        "ERC721: transfer of token that is not own"
      );
    });

    it('failed by witdrowing by owner', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})
      await bridge.depositNFT(nft.address, TokenID1, {from: user1});

      await expectRevert(
        bridge.withdrawNFT(nft.address, user1, TokenID1, {from: user1}),
        "no access permission"
      );
    });

    it('failed by directly calling bank', async () => {
      await bridge.addNFTWhitelist(nft.address, {from: deployer})
      await bridge.depositNFT(nft.address, TokenID1, {from: user1});

      await expectRevert(
        bank.withdrawNFT(nft.address, user1, TokenID1, {from: attacker}),
        "no access permission"
      );
    });
  });
});
