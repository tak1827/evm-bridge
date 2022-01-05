package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Option interface {
	Apply(*Client) error
}

type GasPriceOpt int64

func (o GasPriceOpt) Apply(c *Client) error {
	c.GasPrice = big.NewInt(int64(o))
	return nil
}
func WithGasPrice(gasPrice int64) GasPriceOpt {
	return GasPriceOpt(gasPrice)
}

type CustomComfirm func(h string, recept *types.Receipt) error

func (f CustomComfirm) Apply(c *Client) error {
	c.CustomComfirm = CustomComfirm(f)
	return nil
}
func WithCustomComfirm(f func(h string, recept *types.Receipt) error) CustomComfirm {
	return CustomComfirm(f)
}
