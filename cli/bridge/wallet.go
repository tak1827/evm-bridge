package bridge

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tak1827/nonce-incrementor/nonce"
)

type Wallet struct {
	Nonce   *nonce.Nonce
	priv    *ecdsa.PrivateKey
	privStr string
}

func NewWallet(ctx context.Context, client nonce.Client, privKey string) (w Wallet, err error) {
	w.Nonce, err = nonce.NewNonce(ctx, client, privKey, true)
	if err != nil {
		return
	}

	w.privStr = privKey

	w.priv, err = crypto.HexToECDSA(privKey)
	return
}

func (w Wallet) IncrementNonce() (uint64, error) {
	return w.Nonce.Increment()
}
