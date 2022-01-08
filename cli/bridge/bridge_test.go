package bridge

import (
	"context"
	"crypto/ecdsa"
	"strconv"
	// "sync"
	"testing"
	"time"

	// "github.com/davecgh/go-spew/spew"
	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/evm-bridge/cli/pb"
	"github.com/tak1827/transaction-confirmer/confirm"
)

const (
	Endpoint = "http://localhost:8545"
	BankHex  = "0xA428F7F590C846F290Ee819F13158AdB298fc6c7"
	ERC20Hex = "0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D"
	NFTHex   = "0x61221d7b7978F45A1b51af5492a02Ae6Fc199320"

	PrivKey   = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
	PrivKey2  = "8179ce3d00ac1d1d1d38e4f038de00ccd0e0375517164ac5448e3acc847acb34"
	QueueSize = 256
)

func TestHandleERC20DepositedLogs(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		// ctx  = context.Background()
		rotator = NewRotator(2)
		pair    = pb.Pair{
			Inaddr:  ERC20Hex,
			Outaddr: ERC20Hex,
			Intype:  pb.Pair_ORIGINAL,
		}
		priv2, _       = crypto.HexToECDSA(PrivKey2)
		sentTxs        = make(map[string]pb.Event)
		expectedAmount = int64(10)
		err            error
	)

	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex, NFTHex)
	require.NoError(t, err)

	rc, err := client.NewReadClient(ctx, Endpoint, BankHex)
	require.NoError(t, err)

	confirmer := confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(100))

	bridge, err := NewBridge(ctx, &c, &rc, &confirmer, PrivKey, "")
	require.NoError(t, err)
	require.NoError(t, bridge.Start(ctx))

	err = pair.Put(bridge.db)
	require.NoError(t, err)

	bridge.CustomConfirmedHandler = func(h string) (err error) {
		e, ok := bridge.readEventMap(h)
		if !ok {
			return
		}
		sentTxs[h] = e
		return
	}

	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()

	var (
		counterERC20 int
		counterNFT   int
		breakERC20   = false
		breakNFT     = false
		tokenid      int64
	)
	for {
		<-timer.C
		incrementBlock(t, bridge, ctx, priv2, 3)

		switch rotator.Rotate() {
		case SlotERC20:
			if !(len(bridge.EventMapERC20) == 0) {
				continue
			}

			if counterERC20 >= 3 || breakERC20 {
				breakERC20 = true
				break
			}
			counterERC20++

			batchDepositERC20(t, bridge, ctx, expectedAmount, 3)

			bridge.ConfirmedBlockERC20.Number, err = bridge.FetchERC20(ctx)
			require.NoError(t, err)

		case SlotNFT:
			if !(len(bridge.EventMapNFT) == 0) {
				continue
			}

			if counterNFT >= 3 {
				breakNFT = true
				break
			}
			counterNFT++

			tokenid = batchDepositNFT(t, bridge, ctx, tokenid, 3)

			bridge.ConfirmedBlockNFT.Number, err = bridge.FetchNFT(ctx)
			require.NoError(t, err)

		default:
			panic("unexpexted slot")
		}

		if breakERC20 && breakNFT {
			break
		}
	}

	require.NoError(t, bridge.ConfirmedBlockERC20.Put(bridge.db, pb.BlockERC20))
	require.NoError(t, bridge.ConfirmedBlockNFT.Put(bridge.db, pb.BlockNFT))

	tokenids := make(map[uint64]struct{})
	for _, e := range sentTxs {
		e, err = e.Get(bridge.db)
		require.NoError(t, err)
		require.Equal(t, pb.EventStatus_SUCCEEDED, e.GStatus())

		switch v := e.(type) {
		case pb.EventERC20Deposited:
			amount, err := strconv.Atoi(v.Amount)
			require.NoError(t, err)
			require.Equal(t, expectedAmount, int64(amount))
		case pb.EventNFTDeposited:
			_, ok := tokenids[v.Tokenid]
			require.Equal(t, false, ok)
			tokenids[v.Tokenid] = struct{}{}
		}
	}

	block, err := pb.GetConfirmedBlock(bridge.db, pb.BlockERC20)
	require.NoError(t, err)
	require.Equal(t, true, block.Number > 0)

	block, err = pb.GetConfirmedBlock(bridge.db, pb.BlockNFT)
	require.NoError(t, err)
	require.Equal(t, true, block.Number > 0)

	bridge.Close(cancel, 0, true)
}

// func TestRetry(t *testing.T) {
// 	var (
// 		ctx, cancel = context.WithCancel(context.Background())
// 		pair        = pb.Pair{
// 			Inaddr:  ERC20Hex,
// 			Outaddr: ERC20Hex,
// 			Intype:  pb.Pair_ORIGINAL,
// 		}
// 		m           sync.Mutex
// 		retryFlg    = make(map[uint64]uint32)
// 		sentCounter uint32
// 		amount      = int64(10)
// 		priv2, _    = crypto.HexToECDSA(PrivKey2)
// 		c, _        = client.NewClient(ctx, Endpoint, BankHex, ERC20Hex, NFTHex)
// 		rc, _       = client.NewReadClient(ctx, Endpoint, BankHex)
// 		confirmer   = confirm.NewConfirmer(&c, QueueSize, confirm.WithWorkers(2), confirm.WithWorkerInterval(100))
// 		bridge, _   = NewBridge(ctx, &c, &rc, &confirmer, PrivKey, "")
// 		err         error
// 	)

// 	confirmer.AfterTxSent = func(h string) (err error) {
// 		bridge.Lock()
// 		event, _ := bridge.EventMapERC20[h]
// 		bridge.Unlock()

// 		m.Lock()
// 		defer m.Unlock()

// 		_, ok := retryFlg[event.Id]
// 		if ok {
// 			return
// 		}
// 		retryFlg[event.Id] = sentCounter % 4
// 		sentCounter++
// 		return
// 	}

// 	c.CustomComfirm = func(h string, recept *types.Receipt) error {
// 		bridge.Lock()
// 		event := bridge.EventMapERC20[h]
// 		bridge.Unlock()

// 		m.Lock()
// 		retry := retryFlg[event.Id]
// 		m.Unlock()

// 		if event.Retry < retry || retry == 3 {
// 			recept.Status = 0
// 		}

// 		return nil
// 	}

// 	bridge.Start(ctx)
// 	pair.Put(bridge.db)

// 	timer := time.NewTicker(100 * time.Millisecond)
// 	defer timer.Stop()

// 	var (
// 		counter int
// 		end     uint64
// 	)
// 	for {
// 		<-timer.C
// 		if !bridge.canClose() {
// 			incrementBlock(t, bridge, ctx, priv2, amount, 3)
// 			continue
// 		}

// 		bridge.UpdateStartERC20(end)

// 		if counter >= 3 {
// 			break
// 		}

// 		batchDepositERC20(t, bridge, ctx, amount, 3)

// 		end, err = bridge.HandleERC20DepositedLogs(ctx)
// 		require.NoError(t, err)

// 		counter++
// 	}

// 	for id, retry := range retryFlg {
// 		event, _ := pb.GetEventERC20(bridge.db, id)
// 		require.Equal(t, retry, event.Retry)
// 		if retry >= 3 {
// 			require.Equal(t, pb.EventStatus_FAILED, event.Status)
// 		} else {
// 			require.Equal(t, pb.EventStatus_SUCCEEDED, event.Status)
// 		}
// 	}

// 	bridge.Close(cancel, 0, true)
// }

func incrementBlock(t *testing.T, bridge *Bridge, ctx context.Context, priv *ecdsa.PrivateKey, size int) {
	for i := 0; i < size; i++ {
		_, err := bridge.client.Deposit(ctx, priv, nil, int64(1))
		require.NoError(t, err)
	}
}

func batchDepositERC20(t *testing.T, bridge *Bridge, ctx context.Context, amount int64, size int) {
	for i := 0; i < size; i++ {
		_, err := bridge.client.DepositERC20(ctx, bridge.wallet.priv, nil, amount)
		require.NoError(t, err)
	}
}

func batchDepositNFT(t *testing.T, bridge *Bridge, ctx context.Context, tokenid, size int64) (id int64) {
	for id = tokenid; id < tokenid+size; id++ {
		_, err := bridge.client.DepositNFT(ctx, bridge.wallet.priv, nil, id)
		require.NoError(t, err)
	}
	return
}
