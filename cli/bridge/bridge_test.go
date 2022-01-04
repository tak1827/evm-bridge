package bridge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/transaction-confirmer/confirm"
)

const (
	Endpoint = "http://localhost:8545"
	BankHex  = "0xA7921938DED8fF056f0A77ce1e9Bca2A691e86a1"
	ERC20Hex = "0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D"

	PrivKey   = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
	QueueSize = 65536
)

func TestA(t *testing.T) {
	var (
		ctx = context.Background()
	)

	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex)
	require.NoError(t, err)

	confirmer := confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(100))
	// bridge, err := NewBridge(ctx, &c, &confirmer, PrivKey)
	_, err = NewBridge(ctx, &c, &confirmer, PrivKey)
	require.NoError(t, err)
}
