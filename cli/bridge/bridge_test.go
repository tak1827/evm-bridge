package bridge

import (
	"context"
	"crypto/ecdsa"
	"strconv"
	"sync"
	"testing"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
	PrivKey2  = "8179ce3d00ac1d1d1d38e4f038de00ccd0e0375517164ac5448e3acc847acb34"
	QueueSize = 256
)

func TestHandleERC20DepositedLogs(t *testing.T) {
	var (
		ctx, _ = context.WithCancel(context.Background())
		pair   = pb.Pair{
			Inaddr:  ERC20Hex,
			Outaddr: ERC20Hex,
			Intype:  pb.Pair_ORIGINAL,
		}
		priv2, _       = crypto.HexToECDSA(PrivKey2)
		sentTxs        = make(map[string]pb.EventERC20Deposited)
		expectedAmount = int64(10)
		err            error
	)

	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex)
	require.NoError(t, err)

	rc, err := client.NewReadClient(ctx, Endpoint, BankHex)
	require.NoError(t, err)

	confirmer := confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(100))

	bridge, err := NewBridge(ctx, &c, &rc, &confirmer, PrivKey, "")
	require.NoError(t, err)
	err = bridge.Start(ctx)
	require.NoError(t, err)

	err = pair.Put(bridge.db)
	require.NoError(t, err)

	bridge.CustomConfirmedHandler = func(h string) (err error) {
		bridge.Lock()
		e, ok := bridge.EventERC20s[h]
		bridge.Unlock()
		if !ok {
			return
		}
		sentTxs[h] = e
		return
	}

	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()

	var (
		counter int
		end     uint64
	)
	for {
		<-timer.C
		if !bridge.canClose() {
			incrementBlock(t, bridge, ctx, priv2, expectedAmount, 3)
			continue
		}

		err = bridge.UpdateStartERC20(end)
		require.NoError(t, err)

		if counter >= 3 {
			break
		}

		batchDepositERC20(t, bridge, ctx, expectedAmount, 3)

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

	// bridge.Close(cancel, 0, true)
}

func TestRetry(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		pair        = pb.Pair{
			Inaddr:  ERC20Hex,
			Outaddr: ERC20Hex,
			Intype:  pb.Pair_ORIGINAL,
		}
		m           sync.Mutex
		retryFlg    = make(map[uint64]uint32)
		sentCounter uint32
		amount      = int64(10)
		priv2, _    = crypto.HexToECDSA(PrivKey2)
		c, _        = client.NewClient(ctx, Endpoint, BankHex, ERC20Hex)
		rc, _       = client.NewReadClient(ctx, Endpoint, BankHex)
		confirmer   = confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(100))
		bridge, _   = NewBridge(ctx, &c, &rc, &confirmer, PrivKey, "")
		err         error
	)

	confirmer.AfterTxSent = func(h string) (err error) {
		bridge.Lock()
		event, _ := bridge.EventERC20s[h]
		bridge.Unlock()

		m.Lock()
		defer m.Unlock()

		_, ok := retryFlg[event.Id]
		if ok {
			return
		}
		retryFlg[event.Id] = sentCounter % 4
		sentCounter++
		return
	}

	c.CustomComfirm = func(h string, recept *types.Receipt) error {
		bridge.Lock()
		event := bridge.EventERC20s[h]
		bridge.Unlock()

		m.Lock()
		retry := retryFlg[event.Id]
		m.Unlock()

		if event.Retry < retry || retry == 3 {
			recept.Status = 0
		}

		return nil
	}

	bridge.Start(ctx)
	pair.Put(bridge.db)

	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()

	var (
		counter int
		end     uint64
	)
	for {
		<-timer.C
		if !bridge.canClose() {
			incrementBlock(t, bridge, ctx, priv2, amount, 3)
			continue
		}

		bridge.UpdateStartERC20(end)

		if counter >= 3 {
			break
		}

		batchDepositERC20(t, bridge, ctx, amount, 3)

		end, err = bridge.HandleERC20DepositedLogs(ctx)
		require.NoError(t, err)

		counter++
	}

	for id, retry := range retryFlg {
		event, _ := pb.GetEventERC20(bridge.db, id)
		require.Equal(t, retry, event.Retry)
		if retry >= 3 {
			require.Equal(t, pb.EventStatus_FAILED, event.Status)
		} else {
			require.Equal(t, pb.EventStatus_SUCCEEDED, event.Status)
		}
	}

	bridge.Close(cancel, 0, true)
}

func incrementBlock(t *testing.T, bridge *Bridge, ctx context.Context, priv *ecdsa.PrivateKey, amount int64, size int) {
	for i := 0; i < size; i++ {
		_, err := bridge.client.Deposit(ctx, priv, nil, amount)
		require.NoError(t, err)
	}
}

func batchDepositERC20(t *testing.T, bridge *Bridge, ctx context.Context, amount int64, size int) {
	for i := 0; i < size; i++ {
		_, err := bridge.client.DepositERC20(ctx, bridge.wallet.priv, nil, amount)
		require.NoError(t, err)
	}
}
