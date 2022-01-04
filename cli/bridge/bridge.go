package bridge

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/evm-bridge/cli/log"
	"github.com/tak1827/evm-bridge/cli/pb"
	"github.com/tak1827/go-store/store"
	"github.com/tak1827/transaction-confirmer/confirm"
)

type Bridge struct {
	client    *client.Client
	wallet    Wallet
	confirmer *confirm.Confirmer
	db        store.Store
	logger    zerolog.Logger

	// addrPairMap map[string]string // cache of addr pair

	eventERC20s map[string]pb.EventERC20Deposited
	startERC20  uint64 // confirmed block number of depositERC20 event
}

func NewBridge(ctx context.Context, c *client.Client, confirmer *confirm.Confirmer, privKey string) (b Bridge, err error) {
	b.client = c

	b.confirmer = confirmer
	b.confirmer.AfterTxConfirmed = b.confirmedHandler
	b.confirmer.ErrHandler = b.confirmerErrHandler
	// b.confirmer.Start(ctx)

	if b.db, err = store.NewLevelDB(""); err != nil {
		return
	}

	if b.wallet, err = NewWallet(ctx, c, privKey); err != nil {
		return
	}

	b.logger = log.Bridge("")

	block, err := pb.GetConfirmedBlockERC20(b.db)
	if err == nil {
		b.startERC20 = block.Number
	} else if errors.Is(err, store.ErrNotFound) {
		b.startERC20 = 0
	} else {
		return
	}

	// b.addrPairMap = make(map[string]string)

	return
}

func (b *Bridge) HandleERC20DepositedLogs(ctx context.Context) (err error) {
	eventCh := make(chan pb.EventERC20Deposited, 256)
	if err = b.fetchERC20(ctx, eventCh); err != nil {
		return
	}

	for e := range eventCh {
		b.logger.Info().Msgf("filtered event: %v", e)

		if e, err = pb.GetEventERC20(b.db, e.GetId()); err != nil {
			if !errors.Is(err, store.ErrNotFound) {
				return
			}
		} else if e.Status != pb.EventStatus_UNDEFINED {
			continue
		}

		var hash string
		if hash, err = b.sendERC20(ctx, e); err != nil {
			return
		}

		b.eventERC20s[hash] = e
	}
	// <-closeCh
	// pb.PutERC20ConfirmedBlock(b.db, end)
	return
}

func (b *Bridge) fetchERC20(ctx context.Context, eventCh chan pb.EventERC20Deposited) error {
	end, err := b.client.LatestBlockNumber(ctx)
	if err != nil {
		return err
	}

	go func() {
		defer close(eventCh)

		if err = b.client.FilterERC20Deposited(ctx, b.startERC20, &end, func(e *client.IBankERC20Deposited) error {
			eventCh <- pb.ToEventERC20Deposited(e)
			return nil
		}); err != nil {
			b.logger.Warn().Msgf("fraild filter erc20 logs, err: %v", err)
		}
	}()

	return nil
}

func (b *Bridge) sendERC20(ctx context.Context, e pb.EventERC20Deposited) (hash string, err error) {
	pair, err := pb.GetPair(b.db, e.Token)
	if err != nil {
		return
	}

	switch pair.Intype {
	case pb.Pair_ORIGINAL:
		if hash, err = b.mint(ctx, e); err != nil {
			return
		}
	case pb.Pair_WRAPPED:
		//
	}

	return
}

func (b *Bridge) mint(ctx context.Context, e pb.EventERC20Deposited) (hash string, err error) {
	nonce, err := b.wallet.IncrementNonce()
	if err != nil {
		return
	}

	sender := common.HexToAddress(e.Sender)
	amount := new(big.Int)
	amount.SetString(e.Amount, 10)

	tx, err := b.client.BuildMintTx(ctx, b.wallet.priv, nonce, sender, amount)
	if err != nil {
		return
	}

	if err = b.confirmer.EnqueueTx(ctx, tx); err != nil {
		return
	}

	hash = tx.Hash().Hex()
	return
}

func (b *Bridge) erc20ConfirmedHandler(h string) (err error) {
	e, ok := b.eventERC20s[h]
	if !ok {
		return nil
	}
	defer delete(b.eventERC20s, h)

	e.Status = pb.EventStatus_SUCCEEDED

	if err := e.Put(b.db); err != nil {
		return err
	}

	return
}

func (b *Bridge) erc20ErrHandler(h string, err error) {
	e, ok := b.eventERC20s[h]
	if !ok {
		return
	}

	if !errors.Is(err, confirm.ErrTxFailed) || e.Retry >= 3 {
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

func (b *Bridge) confirmedHandler(h string) (err error) {
	if err = b.erc20ConfirmedHandler(h); err != nil {
		return
	}
	return
}

func (b *Bridge) confirmerErrHandler(h string, err error) {
	b.erc20ErrHandler(h, err)
}
