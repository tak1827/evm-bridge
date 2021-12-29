package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tak1827/evm-bridge/cli/client"
)

const (
	Endpoint = "http://localhost:8545"
	BankHex  = "0xA7921938DED8fF056f0A77ce1e9Bca2A691e86a1"
	ERC20Hex = "0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D"

	PrivKey = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		ctx          = context.Background()
		priv, _      = crypto.HexToECDSA(PrivKey)
		recipient, _ = client.GenerateAddr()
		amount       = big.NewInt(1)
	)

	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex)
	handleErr(err)

	// err = c.FilterERC20Deposited(ctx, 0, nil, func(e *client.IBankERC20Deposited) (stop bool) {
	// 	fmt.Printf("event: %v\n\n", e)
	// 	return false
	// })
	// handleErr(err)

	nonce, err := c.Nonce(ctx, PrivKey)
	handleErr(err)

	fmt.Printf("nonce: %d", nonce)

	tx, err := c.TransferTx(ctx, priv, nonce, recipient, amount)
	handleErr(err)

	hash, err := c.SendTx(ctx, tx)
	handleErr(err)

	spew.Dump(hash)
}
