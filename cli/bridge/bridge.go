package bridge

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/evm-bridge/cli/log"
	"github.com/tak1827/evm-bridge/cli/pb"
	"github.com/tak1827/go-store/store"
	"github.com/tak1827/transaction-confirmer/confirm"
)

type Bridge struct {
	sync.Mutex

	client      *client.Client
	reaadClient *client.ReadClient

	DB store.Store

	wallet    Wallet
	confirmer *confirm.Confirmer
	logger    zerolog.Logger

	// TODO: cache of addr pair

	CustomConfirmedHandler confirm.HashHandler
	CustomErrHandler       confirm.ErrHandler

	EventMapERC20       map[string]*pb.EventERC20Deposited
	EventMapNFT         map[string]*pb.EventNFTDeposited
	ConfirmedBlockERC20 pb.ConfirmedBlock
	ConfirmedBlockNFT   pb.ConfirmedBlock
}

func NewBridge(ctx context.Context, c *client.Client, rc *client.ReadClient, confirmer *confirm.Confirmer, privKey string, path string, opts ...Option) (b *Bridge, err error) {
	b = &Bridge{
		client:        c,
		reaadClient:   rc,
		confirmer:     confirmer,
		logger:        log.Bridge(""),
		EventMapERC20: make(map[string]*pb.EventERC20Deposited),
		EventMapNFT:   make(map[string]*pb.EventNFTDeposited),
	}

	b.confirmer.AfterTxConfirmed = b.confirmedHandler
	b.confirmer.ErrHandler = b.confirmerErrHandler

	if b.DB, err = store.NewLevelDB(path); err != nil {
		return
	}
	if b.wallet, err = NewWallet(ctx, c, privKey); err != nil {
		return
	}
	if err = b.ConfirmedBlockERC20.Get(b.DB, pb.BlockERC20); err != nil {
		return
	}
	if err = b.ConfirmedBlockNFT.Get(b.DB, pb.BlockNFT); err != nil {
		return
	}

	for i := 0; i < len(opts); i++ {
		opts[i].Apply(b)
	}

	return
}

func (b *Bridge) Start(ctx context.Context) (err error) {
	b.logger.Info().Msg("bridge is starting...")
	err = b.confirmer.Start(ctx)
	return
}

func (b *Bridge) Close(cancel context.CancelFunc, retryLimit int, commitStarts bool) {
	if !b.canClose() {
		b.logger.Info().Msg("closing...")
		// wait until all confirmed
		timer := time.NewTicker(1 * time.Second)
		defer timer.Stop()

		var retry int
		for {
			<-timer.C
			if b.canClose() {
				break
			}

			b.logger.Info().Msgf("trying safety close, retry:%d, max limit: %d", retry, retryLimit)
			if retry >= retryLimit {
				b.logger.Warn().Msgf("faild closing safely, EventMapERC20: %v, EventMapNFT: %v", b.EventMapERC20, b.EventMapNFT)
				break
			}
			retry++
		}
	}

	if commitStarts {
		if err := b.ConfirmedBlockERC20.Put(b.DB, pb.BlockERC20); err != nil {
			b.logger.Warn().Msgf("faild to put ConfirmedBlockERC20(%v)", b.ConfirmedBlockERC20)
		}
		if err := b.ConfirmedBlockNFT.Put(b.DB, pb.BlockNFT); err != nil {
			b.logger.Warn().Msgf("faild to put ConfirmedBlockNFT(%v)", b.ConfirmedBlockNFT)
		}
		b.logger.Info().Msgf("commited the last confirmed blocks, erc20: %d, nft: %d", b.ConfirmedBlockERC20.Number, b.ConfirmedBlockNFT.Number)
	}

	b.confirmer.Close(cancel)

	if err := b.DB.Close(); err != nil {
		b.logger.Warn().Msg("faild to close db")
	}
}

func (b *Bridge) canClose() bool {
	b.Lock()
	defer b.Unlock()

	return len(b.EventMapERC20) == 0 && len(b.EventMapNFT) == 0
}

func (b *Bridge) FetchERC20(ctx context.Context) (uint64, error) {
	eventCh := make(chan pb.Event, 256)

	end, err := b.reaadClient.LatestBlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	go func() {
		defer close(eventCh)

		if err = b.reaadClient.FilterERC20Deposited(ctx, b.ConfirmedBlockERC20.Number, &end, func(e *client.IBankERC20Deposited) error {
			eventCh <- pb.ToEventERC20Deposited(e)
			return nil
		}); err != nil {
			b.logger.Warn().Msgf("failed filter erc20 logs, err: %v", err)
		}
	}()

	err = b.handleLogs(ctx, eventCh)
	return end, err
}

func (b *Bridge) FetchNFT(ctx context.Context) (uint64, error) {
	var (
		eventCh  = make(chan pb.Event, 256)
		start    = b.ConfirmedBlockNFT.Number
		end, err = b.reaadClient.LatestBlockNumber(ctx)
	)
	if err != nil {
		return 0, err
	}

	go func() {
		defer close(eventCh)

		if err = b.reaadClient.FilterNFTDeposited(ctx, start, &end, func(e *client.IBankNFTDeposited) error {
			eventCh <- pb.ToEventNFTDeposited(e)
			return nil
		}); err != nil {
			b.logger.Warn().Msgf("failed filter nft logs, err: %v", err)
		}
	}()

	err = b.handleLogs(ctx, eventCh)
	return end, err
}

func (b *Bridge) handleLogs(ctx context.Context, eventCh chan pb.Event) error {
	for e := range eventCh {
		b.logger.Info().Msgf("handling event: %v", e)

		if err := e.Get(b.DB); err != nil {
			if !errors.Is(err, store.ErrNotFound) {
				return err
			}
		} else if e.GetStatus() == pb.EventStatus_SUCCEEDED {
			continue
		}

		if _, err := b.send(ctx, e); err != nil {
			if errors.Is(err, ErrPairNotFound) {
				b.logger.Warn().Msgf("pir not found, event: %v, err: %v", e, err)
				continue
			}
			return err
		}
	}

	return nil
}

func (b *Bridge) send(ctx context.Context, e pb.Event) (hash string, err error) {
	pair, err := pb.GetPair(b.DB, e.GetToken())
	if err != nil {
		err = ErrPairNotFound
		return
	}

	var (
		tx *types.Transaction
		to = common.HexToAddress(pair.Outaddr)
	)
	switch pair.Intype {
	case pb.Pair_ORIGINAL:
		if tx, err = b.mint(ctx, e, to); err != nil {
			return
		}
	case pb.Pair_WRAPPED:
		// TODO:
	}

	hash = tx.Hash().Hex()
	b.writeEventMap(hash, e)

	err = b.confirmer.EnqueueTx(ctx, tx)
	return
}

func (b *Bridge) mint(ctx context.Context, e pb.Event, to common.Address) (tx *types.Transaction, err error) {
	nonce, err := b.wallet.IncrementNonce()
	if err != nil {
		return
	}

	switch v := e.(type) {
	case *pb.EventERC20Deposited:
		sender := common.HexToAddress(v.Sender)
		amount := new(big.Int)
		amount.SetString(v.Amount, 10)
		tx, err = b.client.BuildERC20MintTx(ctx, b.wallet.priv, nonce, to, sender, amount)
	case *pb.EventNFTDeposited:
		sender := common.HexToAddress(v.Sender)
		tokenid := big.NewInt(int64(v.Tokenid))
		tx, err = b.client.BuildNFTMintTx(ctx, b.wallet.priv, nonce, to, sender, tokenid)
	default:
		panic(fmt.Sprintf("unexpected type(%T)\n", v))
	}

	return
}

func (b *Bridge) confirmedHandler(h string) (err error) {
	if b.CustomConfirmedHandler != nil {
		if err = b.CustomConfirmedHandler(h); err != nil {
			return
		}
	}

	e, exist := b.readEventMap(h)
	if !exist {
		return
	}
	b.deleteEventMap(h)

	e.SetStatus(pb.EventStatus_SUCCEEDED)

	if err = e.Put(b.DB); err != nil {
		return
	}

	b.logger.Info().Msgf("confirmed, hash: %s, event: %v", h, e)

	return
}

func (b *Bridge) confirmerErrHandler(h string, err error) {
	if b.CustomErrHandler != nil {
		b.CustomErrHandler(h, err)
	}

	e, exist := b.readEventMap(h)
	if !exist {
		return
	}
	b.deleteEventMap(h)

	if e.GetRetry() >= 3 || !errors.Is(err, confirm.ErrTxFailed) {
		b.logger.Warn().Msgf("failed handle erc20 log(%v), hash: %s, err: %v", e, h, err)
		e.SetStatus(pb.EventStatus_FAILED)
		if err = e.Put(b.DB); err != nil {
			b.logger.Warn().Msgf("failed to put event(%v), hash: %s, err: %v", e, h, err)
		}
		return
	}

	e.SetRetry(e.GetRetry() + 1)

	if _, err = b.send(context.Background(), e); err != nil {
		b.logger.Warn().Msgf("failed mint event(%v), hash: %s, err: %v", e, h, err)
	}
}

func (b *Bridge) writeEventMap(h string, e pb.Event) {
	b.Lock()
	defer b.Unlock()

	switch v := e.(type) {
	case *pb.EventERC20Deposited:
		b.EventMapERC20[h] = v
	case *pb.EventNFTDeposited:
		b.EventMapNFT[h] = v
	default:
		panic(fmt.Sprintf("unexpected type(%T)\n", v))
	}
}

func (b *Bridge) readEventMap(h string) (e pb.Event, exist bool) {
	b.Lock()
	defer b.Unlock()

	e, exist = b.EventMapERC20[h]
	if exist {
		return
	}
	e, exist = b.EventMapNFT[h]
	if exist {
		return
	}
	return
}

func (b *Bridge) deleteEventMap(h string) {
	b.Lock()
	defer b.Unlock()

	exist := false
	if _, exist = b.EventMapERC20[h]; exist {
		delete(b.EventMapERC20, h)
		return
	}
	if _, exist = b.EventMapNFT[h]; exist {
		delete(b.EventMapNFT, h)
		return
	}

	return
}
