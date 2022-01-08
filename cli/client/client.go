package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tak1827/nonce-incrementor/nonce"
	"github.com/tak1827/transaction-confirmer/confirm"
)

const (
	DefaultGasPrice = int64(0) // 1 gwai
)

var (
	_ confirm.Client = (*Client)(nil)
	_ nonce.Client   = (*Client)(nil)
)

type Client struct {
	ethclient *ethclient.Client
	GasPrice  *big.Int

	erc20ABI  abi.ABI
	nftABI    abi.ABI
	Bank      *IBank
	ERC20     *IERC20
	NFT       *IERC721
	bankAddr  common.Address
	erc20Addr common.Address
	nftAddr   common.Address

	CustomComfirm func(h string, recept *types.Receipt) error
}

func NewClient(ctx context.Context, endpoint string, bankHex, erc20Hex, nftHex string, opts ...Option) (c Client, err error) {
	rpcclient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		err = fmt.Errorf("failed to conecting endpoint(%s) err: %w", endpoint, err)
		return
	}

	c.ethclient = ethclient.NewClient(rpcclient)
	c.GasPrice = big.NewInt(int64(DefaultGasPrice))
	c.bankAddr = common.HexToAddress(bankHex)
	c.erc20Addr = common.HexToAddress(erc20Hex)
	c.nftAddr = common.HexToAddress(nftHex)

	if c.erc20ABI, err = abi.JSON(strings.NewReader(IERC20ABI)); err != nil {
		return
	}

	if c.Bank, err = NewIBank(c.bankAddr, c.ethclient); err != nil {
		return
	}
	if c.ERC20, err = NewIERC20(c.erc20Addr, c.ethclient); err != nil {
		return
	}
	if c.NFT, err = NewIERC721(c.nftAddr, c.ethclient); err != nil {
		return
	}

	for i := 0; i < len(opts); i++ {
		opts[i].Apply(&c)
	}

	return
}

func (c *Client) Nonce(ctx context.Context, privKey string) (nonce uint64, err error) {
	priv, err := crypto.HexToECDSA(privKey)
	if err != nil {
		err = errors.Wrap(err, "failed to get nonce")
		return
	}

	account := crypto.PubkeyToAddress(priv.PublicKey)
	nonce, err = c.ethclient.NonceAt(ctx, account, nil)
	return
}

func (c *Client) SendTx(ctx context.Context, tx interface{}) (string, error) {
	signedTx := tx.(*types.Transaction)

	if err := c.ethclient.SendTransaction(ctx, signedTx); err != nil {
		return "", errors.Wrap(err, "err SendTransaction")
	}

	return signedTx.Hash().Hex(), nil
}

func (c *Client) ConfirmTx(ctx context.Context, hash string, confirmationBlocks uint64) error {
	recept, err := c.Receipt(ctx, hash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return confirm.ErrTxNotFound
		}
		return errors.Wrap(err, "err TransactionReceipt")
	}

	if c.CustomComfirm != nil {
		if err = c.CustomComfirm(hash, recept); err != nil {
			return err
		}
	}

	if recept.Status != 1 {
		return confirm.ErrTxFailed
	}

	block, err := c.LatestBlockNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "err LatestBlockNumber")
	}

	if recept.BlockNumber.Uint64()+confirmationBlocks > block {
		return confirm.ErrTxConfirmPending
	}

	return nil
}

func (c *Client) Receipt(ctx context.Context, hash string) (*types.Receipt, error) {
	return c.ethclient.TransactionReceipt(ctx, common.HexToHash(hash))
}

func (c *Client) LatestBlockNumber(ctx context.Context) (uint64, error) {
	header, err := c.ethclient.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (c *Client) BuildTx(priv *ecdsa.PrivateKey, nonce uint64, to common.Address, value *big.Int, gasLimit uint64, data []byte) (*types.Transaction, error) {
	var (
		tx     = types.NewTransaction(nonce, to, value, gasLimit, c.GasPrice, data)
		signer = types.HomesteadSigner{}
	)

	sig, err := crypto.Sign(signer.Hash(tx).Bytes(), priv)
	if err != nil {
		return nil, errors.Wrap(err, "err Sign")
	}

	return tx.WithSignature(signer, sig)
}

func (c *Client) BuildERC20MintTx(ctx context.Context, priv *ecdsa.PrivateKey, nonce uint64, account common.Address, amount *big.Int) (*types.Transaction, error) {
	var (
		auth     = bind.NewKeyedTransactor(priv)
		to       = c.erc20Addr
		input, _ = c.erc20ABI.Pack("mint", account, amount)
		msg      = ethereum.CallMsg{
			From:     auth.From,
			To:       &to,
			GasPrice: c.GasPrice,
			Data:     input,
		}
	)

	gas, err := c.ethclient.EstimateGas(ctx, msg)
	if err != nil {
		return nil, err
	}

	return c.BuildTx(priv, nonce, to, nil, gas, input)
}

func (c *Client) BuildNFTMintTx(ctx context.Context, priv *ecdsa.PrivateKey, nonce uint64, account common.Address, tokenid *big.Int) (*types.Transaction, error) {
	var (
		auth     = bind.NewKeyedTransactor(priv)
		to       = c.nftAddr
		input, _ = c.nftABI.Pack("safeMint", tokenid, account)
		msg      = ethereum.CallMsg{
			From:     auth.From,
			To:       &to,
			GasPrice: c.GasPrice,
			Data:     input,
		}
	)

	gas, err := c.ethclient.EstimateGas(ctx, msg)
	if err != nil {
		return nil, err
	}

	return c.BuildTx(priv, nonce, to, nil, gas, input)
}

func (c *Client) DepositERC20(ctx context.Context, priv *ecdsa.PrivateKey, nonce *big.Int, amount int64) (*types.Transaction, error) {
	var (
		auth = bind.NewKeyedTransactor(priv)
		a    = big.NewInt(amount)
		opts = &bind.TransactOpts{
			From:     auth.From,
			Nonce:    nonce, // nil = use pending state
			Signer:   auth.Signer,
			GasPrice: c.GasPrice,
			GasLimit: 0, // estimate
			Context:  ctx,
		}
	)
	return c.Bank.DepositERC20(opts, c.erc20Addr, auth.From, a)
}

func (c *Client) DepositNFT(ctx context.Context, priv *ecdsa.PrivateKey, nonce *big.Int, tokenid int64) (*types.Transaction, error) {
	var (
		auth = bind.NewKeyedTransactor(priv)
		id   = big.NewInt(tokenid)
		opts = &bind.TransactOpts{
			From:     auth.From,
			Nonce:    nonce, // nil = use pending state
			Signer:   auth.Signer,
			GasPrice: c.GasPrice,
			GasLimit: 0, // estimate
			Context:  ctx,
		}
	)
	return c.Bank.DepositNFT(opts, c.nftAddr, auth.From, id)
}

func (c *Client) Deposit(ctx context.Context, priv *ecdsa.PrivateKey, nonce *big.Int, amount int64) (*types.Transaction, error) {
	var (
		auth = bind.NewKeyedTransactor(priv)
		a    = big.NewInt(amount)
		opts = &bind.TransactOpts{
			From:     auth.From,
			Nonce:    nonce, // nil = use pending state
			Signer:   auth.Signer,
			GasPrice: c.GasPrice,
			GasLimit: 0, // estimate
			Context:  ctx,
			Value:    a,
		}
	)
	return c.Bank.Deposit(opts, auth.From)
}

type ReadClient struct {
	ethclient *ethclient.Client

	Bank     *IBank
	bankAddr common.Address
}

func NewReadClient(ctx context.Context, endpoint string, bankHex string) (c ReadClient, err error) {
	rpcclient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		err = fmt.Errorf("failed to conecting endpoint(%s) err: %w", endpoint, err)
		return
	}

	c.ethclient = ethclient.NewClient(rpcclient)
	c.bankAddr = common.HexToAddress(bankHex)

	if c.Bank, err = NewIBank(c.bankAddr, c.ethclient); err != nil {
		return
	}

	return
}

func (c *ReadClient) LatestBlockNumber(ctx context.Context) (uint64, error) {
	header, err := c.ethclient.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (c *ReadClient) FilterERC20Deposited(ctx context.Context, start uint64, end *uint64, handle func(e *IBankERC20Deposited) error) error {
	opt := bind.FilterOpts{
		Start:   start,
		End:     end,
		Context: ctx,
	}

	it, err := c.Bank.FilterERC20Deposited(&opt, nil, nil)
	if err != nil {
		return err
	}

	for it.Next() {
		if err := handle(it.Event); err != nil {
			return err
		}
	}

	return nil
}

func (c *ReadClient) FilterNFTDeposited(ctx context.Context, start uint64, end *uint64, handle func(e *IBankNFTDeposited) error) error {
	opt := bind.FilterOpts{
		Start:   start,
		End:     end,
		Context: ctx,
	}

	it, err := c.Bank.FilterNFTDeposited(&opt, nil, nil)
	if err != nil {
		return err
	}

	for it.Next() {
		if err := handle(it.Event); err != nil {
			return err
		}
	}

	return nil
}

func GenerateAddr() (addr common.Address, err error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	addr = crypto.PubkeyToAddress(priv.PublicKey)
	return
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}
