# EVM Compatible Chain Bridge
Bridging between EVM compatable chains, like Ethereum and Quorum for moving erc20 tokens back and forth

# PreRequirements
|  Software  |  Version  |
| ---- | ---- |
|  Truffle  |  ^v5.x  |
|  Ganache CLI  |  ^v6.x  |

# Getting start
```bash
# install dependencies
npm install

# run test node
npm run chain

# export env variables
export DEPLOYER_KEY=XXX...
export NODE_URL=https://bsc-dataseed.binance.org/
export ADMIN_ACCOUNT=0x000... # admin will be granted the access roles of bridge and bank contract

# run test
npm run test

# deploy to test net
npm run migrate:bridge:bsctest
```

# How bridge work
```javascript
/*
 * Bridge ERC20
 */
// first, the adim add erc20 to whitelist, only whitelisted erc20 can be bridged
await bridge.addERC20Whitelist(erc20.address, {from: admin});
// second, a user approve bank contract
const amount = 1000;
await erc20.approve(bank.address, amount, {from: user1});
// third, a use deposit erc20 via bridge contract
await bridge.depositERC20(erc20.address, amount, {from: user1});
// then, the bridge cli(in cli directory) will issue wrapped token to the other chain

/*
 * Bridge NFT
 */
// first, the adim add nft to whitelist, only whitelisted nft can be bridged
await bridge.addNFTWhitelist(nft.address, {from: admin});
// second, a user approve bank contract
const tokenid = 101;
await nft.approve(bank.address, tokenid, {from: user1});
// third, a use deposit nft via bridge contract
await bridge.depositNFT(nft.address, tokenid, {from: user1});
// then, the bridge cli(in cli directory) will issue wrapped token to the other chain
```

# How to use bridgecli
```sh
# set path for installed bridgecli
export GOPATH=$(go env GOPATH)
export PATH="$GOPATH/bin:$PATH"

# install bridgecli
make install

# confirm instalation
bridgecli -h

# initalize the home directory in where the configuration file is created
bridgecli init --home ./storage

# edite the configuration file
# NOTE: "in-endpoint" is the source chain, "out-endpoint" is the destination chain
vi ./storage/config.toml

# set the erc20 and nft contract address pairs
# NOTE: set [source-chain-address] [destination-chain-address]
bridgecli pair set 0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D 0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D --home ./storage
bridgecli pair set 0x2518a5D597F670F21Dd4eE989698E18127B3a065 0x61221d7b7978F45A1b51af5492a02Ae6Fc199320  --home ./storage

# confirm the set addresses
bridgecli pair get 0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D --home ./storage
bridgecli pair get 0x2518a5D597F670F21Dd4eE989698E18127B3a065 --home ./storage

# set the private key of the destination chain, used for minting erc20 and nft.
export BRIDGECLI_PRI_KEY=XXXX..

# start bridging service
bridgecli serve --home ./storage
```
