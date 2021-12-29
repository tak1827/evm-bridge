package main

import (
	"context"
	"fmt"

	"github.com/tak1827/evm-bridge/cli/client"
)

const (
	Endpoint = "http://localhost:8545"
	BankHex  = "0xA7921938DED8fF056f0A77ce1e9Bca2A691e86a1"
	ERC20Hex = "0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		ctx = context.Background()
	)

	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex)
	handleErr(err)

	err = c.FilterERC20Deposited(ctx, 0, nil, func(e *client.IBankERC20Deposited) (stop bool) {
		fmt.Printf("event: %v\n\n", e)
		return false
	})
	handleErr(err)
}
