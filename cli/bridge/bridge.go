package bridge

import (
	"context"
	"errors"
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

	wallet    Wallet
	confirmer *confirm.Confirmer
	db        store.Store
	logger    zerolog.Logger

	// addrPairMap map[string]string // cache of addr pair

	CustomConfirmedHandler confirm.HashHandler
	CustomErrHandler       confirm.ErrHandler

	EventERC20s map[string]pb.EventERC20Deposited

	StartERC20 uint64 // confirmed block number of depositERC20 event
}

func NewBridge(ctx context.Context, c *client.Client, rc *client.ReadClient, confirmer *confirm.Confirmer, privKey string, path string, opts ...Option) (b *Bridge, err error) {
	b = &Bridge{
		client:      c,
		reaadClient: rc,
		confirmer:   confirmer,
		logger:      log.Bridge(""),
		EventERC20s: make(map[string]pb.EventERC20Deposited),
	}

	b.confirmer.AfterTxConfirmed = b.confirmedHandler
	b.confirmer.ErrHandler = b.confirmerErrHandler

	if b.db, err = store.NewLevelDB(path); err != nil {
		return
	}

	if b.wallet, err = NewWallet(ctx, c, privKey); err != nil {
		return
	}

	block, err := pb.GetConfirmedBlockERC20(b.db)
	if err == nil {
		b.StartERC20 = block.Number
	} else if errors.Is(err, store.ErrNotFound) {
		b.StartERC20 = 0
		err = nil
	} else {
		return
	}

	for i := 0; i < len(opts); i++ {
		opts[i].Apply(b)
	}

	return
}

func (b *Bridge) Start(ctx context.Context) (err error) {
	err = b.confirmer.Start(ctx)
	return
}

func (b *Bridge) Close(cancel context.CancelFunc, retryLimit int, commitStarts bool) {
	if !b.canClose() {
		// wait until all confirmed
		timer := time.NewTicker(1 * time.Second)
		defer timer.Stop()

		var retry int
		for {
			select {
			case <-timer.C:
				if b.canClose() {
					break
				}
				if retry >= retryLimit {
					b.logger.Warn().Msgf("faild closing safely, EventERC20s: %v", b.EventERC20s)
					break
				}
				retry++
			}
		}
	}

	if commitStarts {
		if err := b.UpdateStartERC20(b.StartERC20); err != nil {
			b.logger.Warn().Msgf("faild to put ConfirmedBlockERC20, StartERC20: %d", b.StartERC20)
		}
	}

	b.confirmer.Close(cancel)

	if err := b.db.Close(); err != nil {
		b.logger.Warn().Msg("faild to close db")
	}
}

func (b *Bridge) canClose() bool {
	b.Lock()
	defer b.Unlock()

	return len(b.EventERC20s) == 0
}

func (b *Bridge) UpdateStartERC20(num uint64) error {
	return pb.PutConfirmedBlockERC20(b.db, num)
}

func (b *Bridge) HandleERC20DepositedLogs(ctx context.Context) (end uint64, err error) {
	eventCh := make(chan pb.EventERC20Deposited, 256)
	if end, err = b.fetchERC20(ctx, eventCh); err != nil {
		return
	}

	for e := range eventCh {
		b.logger.Info().Msgf("filtered event: %v", e)

		if storedEvent, err := pb.GetEventERC20(b.db, e.GetId()); err != nil {
			if !errors.Is(err, store.ErrNotFound) {
				return end, err
			}
		} else if storedEvent.Status == pb.EventStatus_SUCCEEDED {
			continue
		}

		if _, err = b.sendERC20(ctx, e); err != nil {
			if errors.Is(err, ErrPairNotFound) {
				b.logger.Warn().Msgf("pir not found, event: %v, err: %v", e, err)
				continue
			}
			return
		}
	}

	return
}

func (b *Bridge) fetchERC20(ctx context.Context, eventCh chan pb.EventERC20Deposited) (uint64, error) {
	end, err := b.reaadClient.LatestBlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	go func() {
		defer close(eventCh)

		if err = b.reaadClient.FilterERC20Deposited(ctx, b.StartERC20, &end, func(e *client.IBankERC20Deposited) error {
			eventCh <- pb.ToEventERC20Deposited(e)
			return nil
		}); err != nil {
			b.logger.Warn().Msgf("fraild filter erc20 logs, err: %v", err)
		}
	}()

	return end, nil
}

func (b *Bridge) sendERC20(ctx context.Context, e pb.EventERC20Deposited) (hash string, err error) {
	pair, err := pb.GetPair(b.db, e.Token)
	if err != nil {
		err = ErrPairNotFound
		return
	}

	var tx *types.Transaction
	switch pair.Intype {
	case pb.Pair_ORIGINAL:
		if tx, err = b.mint(ctx, e); err != nil {
			return
		}
	case pb.Pair_WRAPPED:
		//
	}

	hash = tx.Hash().Hex()
	b.writeEventERC20s(hash, e)

	err = b.confirmer.EnqueueTx(ctx, tx)
	return
}

func (b *Bridge) mint(ctx context.Context, e pb.EventERC20Deposited) (tx *types.Transaction, err error) {
	nonce, err := b.wallet.IncrementNonce()
	if err != nil {
		return
	}

	sender := common.HexToAddress(e.Sender)
	amount := new(big.Int)
	amount.SetString(e.Amount, 10)

	return b.client.BuildMintTx(ctx, b.wallet.priv, nonce, sender, amount)
}

func (b *Bridge) confirmedHandler(h string) (err error) {
	if b.CustomConfirmedHandler != nil {
		if err = b.CustomConfirmedHandler(h); err != nil {
			return
		}
	}
	if err = b.erc20ConfirmedHandler(h); err != nil {
		return
	}
	return
}

func (b *Bridge) erc20ConfirmedHandler(h string) (err error) {
	e, exist := b.readEventERC20s(h)
	if !exist {
		return
	}

	b.deleteEventERC20s(h)

	b.logger.Info().Msgf("erc20 confirmed, h: %s", h)

	e.Status = pb.EventStatus_SUCCEEDED

	if err = e.Put(b.db); err != nil {
		return
	}

	return
}

func (b *Bridge) confirmerErrHandler(h string, err error) {
	if b.CustomErrHandler != nil {
		b.CustomErrHandler(h, err)
	}
	b.erc20ErrHandler(h, err)
}

func (b *Bridge) erc20ErrHandler(h string, err error) {
	e, exist := b.readEventERC20s(h)
	if !exist {
		return
	}
	b.deleteEventERC20s(h)

	if e.Retry >= 3 || !errors.Is(err, confirm.ErrTxFailed) {
		b.logger.Warn().Msgf("failed handle erc20 log(%v), hash: %s, err: %v", e, h, err)
		e.Status = pb.EventStatus_FAILED
		if err = e.Put(b.db); err != nil {
			b.logger.Warn().Msgf("failed to put event(%v), hash: %s, err: %v", e, h, err)
		}
		return
	}

	e.Retry++

	if _, err = b.sendERC20(context.Background(), e); err != nil {
		b.logger.Warn().Msgf("failed mint event(%v), hash: %s, err: %v", e, h, err)
	}
}

func (b *Bridge) readEventERC20s(h string) (pb.EventERC20Deposited, bool) {
	b.Lock()
	e, exist := b.EventERC20s[h]
	b.Unlock()

	return e, exist
}

func (b *Bridge) writeEventERC20s(h string, e pb.EventERC20Deposited) {
	b.Lock()
	b.EventERC20s[h] = e
	b.Unlock()
}

func (b *Bridge) deleteEventERC20s(h string) {
	b.Lock()
	delete(b.EventERC20s, h)
	b.Unlock()
}
