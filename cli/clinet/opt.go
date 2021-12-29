package client

import (
	"math/big"
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
