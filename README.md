# quorum-bridge
Bridging between Ethereum and Quorum for moving erc20 tokens back and forth

# PreRequirements
|  Software  |  Version  |
| ---- | ---- |
|  Truffle  |  ^v5.x  |
|  Ganache CLI  |  ^v6.x  |

# Getting start
```bash
# install dependencies
npm install

# run testing chain
npm run chain

# export envs
export DEPLOYER_KEY=XXX...
export NODE_URL=https://bsc-dataseed.binance.org/

# deploy bridge contract
npm run migrate:bsctest
```
