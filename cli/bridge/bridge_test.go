package bridge

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/evm-bridge/cli/pb"
	"github.com/tak1827/transaction-confirmer/confirm"
)

const (
	Endpoint = "http://localhost:8545"
	BankHex  = "0xA7921938DED8fF056f0A77ce1e9Bca2A691e86a1"
	ERC20Hex = "0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D"

	PrivKey   = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
	QueueSize = 256
)

func TestHandleERC20DepositedLogs(t *testing.T) {
	var (
		ctx, cancel    = context.WithCancel(context.Background())
		sentTxs        = make(map[string]pb.EventERC20Deposited)
		custom         confirm.HashHandler
		expectedAmount = int64(10)
	)
	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex)
	require.NoError(t, err)

	rc, err := client.NewReadClient(ctx, Endpoint, BankHex)
	require.NoError(t, err)

	confirmer := confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(100))
	confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(10))
	bridge, err := NewBridge(ctx, &c, &rc, &confirmer, PrivKey, "", custom)
	require.NoError(t, err)
	err = bridge.Start(ctx)
	require.NoError(t, err)

	custom = func(h string) (err error) {
		e, ok := bridge.EventERC20s[h]
		if !ok {
			return
		}
		sentTxs[h] = e
		return
	}

	pair := pb.Pair{
		Inaddr:  ERC20Hex,
		Outaddr: ERC20Hex,
		Intype:  pb.Pair_ORIGINAL,
	}
	err = pair.Put(bridge.db)
	require.NoError(t, err)

	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()

	var (
		counter int
		end     uint64
	)
	for {
		<-timer.C
		if !bridge.canClose() {
			incrementBlock(t, &bridge, ctx, expectedAmount, 3)
			continue
		}

		bridge.StartERC20 = end
		err = bridge.UpdateStartERC20(bridge.StartERC20)
		require.NoError(t, err)

		if counter >= 3 {
			break
		}

		batchDepositERC20(t, &bridge, ctx, expectedAmount, 3)

		end, err = bridge.HandleERC20DepositedLogs(ctx)
		require.NoError(t, err)

		counter++
	}

	for _, e := range sentTxs {
		e, err = e.Get(bridge.db)
		require.NoError(t, err)
		require.Equal(t, pb.EventStatus_SUCCEEDED, e.Status)
		amount, err := strconv.Atoi(e.Amount)
		require.NoError(t, err)
		require.Equal(t, expectedAmount, int64(amount))
	}

	block, err := pb.GetConfirmedBlockERC20(bridge.db)
	require.NoError(t, err)
	require.Equal(t, end, block.Number)

	bridge.Close(cancel, 0, true)
}

// func TestRetry(t *testing.T) {
// 	var (
// 		ctx, cancel    = context.WithCancel(context.Background())
// 		sentTxs        = make(map[string]pb.EventERC20Deposited)
// 		expectedAmount = int64(10)
// 	)
// 	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex)
// 	require.NoError(t, err)

// 	rc, err := client.NewReadClient(ctx, Endpoint, BankHex)
// 	require.NoError(t, err)

// 	confirmer := confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(100))
// 	confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(10))
// 	bridge, err := NewBridge(ctx, &c, &rc, &confirmer, PrivKey, "", nil)
// 	require.NoError(t, err)
// 	err = bridge.Start(ctx)
// 	require.NoError(t, err)

// 	pair := pb.Pair{
// 		Inaddr:  ERC20Hex,
// 		Outaddr: ERC20Hex,
// 		Intype:  pb.Pair_ORIGINAL,
// 	}
// 	err = pair.Put(bridge.db)
// 	require.NoError(t, err)

// 	timer := time.NewTicker(100 * time.Millisecond)
// 	defer timer.Stop()

// 	var (
// 		counter int
// 		end     uint64
// 	)
// 	for {
// 		<-timer.C
// 		if !bridge.canClose() {
// 			incrementBlock(t, &bridge, ctx, expectedAmount, 3)
// 			continue
// 		}

// 		bridge.StartERC20 = end
// 		err = bridge.UpdateStartERC20(bridge.StartERC20)
// 		require.NoError(t, err)

// 		if counter >= 3 {
// 			break
// 		}

// 		batchDepositERC20(t, &bridge, ctx, expectedAmount, 3)

// 		end, err = bridge.HandleERC20DepositedLogs(ctx)
// 		require.NoError(t, err)

// 		counter++
// 	}

// 	for _, e := range sentTxs {
// 		e, err = e.Get(bridge.db)
// 		require.NoError(t, err)
// 		require.Equal(t, pb.EventStatus_SUCCEEDED, e.Status)
// 		amount, err := strconv.Atoi(e.Amount)
// 		require.NoError(t, err)
// 		require.Equal(t, expectedAmount, int64(amount))
// 	}

// 	block, err := pb.GetConfirmedBlockERC20(bridge.db)
// 	require.NoError(t, err)
// 	require.Equal(t, end, block.Number)

// 	bridge.Close(cancel, 0, true)
// }

func incrementBlock(t *testing.T, bridge *Bridge, ctx context.Context, amount int64, size int) {
	for i := 0; i < size; i++ {
		_, err := bridge.client.Deposit(ctx, bridge.wallet.priv, nil, amount)
		require.NoError(t, err)
	}
}

func batchDepositERC20(t *testing.T, bridge *Bridge, ctx context.Context, amount int64, size int) {
	for i := 0; i < size; i++ {
		_, err := bridge.client.DepositERC20(ctx, bridge.wallet.priv, nil, amount)
		require.NoError(t, err)
	}
}
