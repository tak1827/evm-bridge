package bridge

// import (
// 	"context"
// 	"crypto/ecdsa"
// 	"strconv"
// 	// "sync"
// 	"testing"
// 	"time"

// 	// "github.com/davecgh/go-spew/spew"
// 	// "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/stretchr/testify/require"
// 	"github.com/tak1827/evm-bridge/cli/client"
// 	"github.com/tak1827/evm-bridge/cli/pb"
// 	"github.com/tak1827/transaction-confirmer/confirm"
// )

// const (
// 	Endpoint  = "http://localhost:8545"
// 	BankHex   = "0x4c2310DAdb5Be92a39336316f841e1944DA7bd60"
// 	ERC20Hex  = "0xe868feADdAA8965b6e64BDD50a14cD41e3D5245D"
// 	NFTHexIn  = "0x2518a5D597F670F21Dd4eE989698E18127B3a065"
// 	NFTHexOut = "0x61221d7b7978F45A1b51af5492a02Ae6Fc199320"

// 	PrivKey   = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
// 	PrivKey2  = "8179ce3d00ac1d1d1d38e4f038de00ccd0e0375517164ac5448e3acc847acb34"
// 	QueueSize = 256
// )

// func TestTmp(t *testing.T) {
// 	var (
// 		// ctx, cancel = context.WithCancel(context.Background())
// 		ctx  = context.Background()
// 		err            error
// 	)

// 	c, err := client.NewClient(ctx, Endpoint, BankHex, ERC20Hex, NFTHexIn)
// 	require.NoError(t, err)
// }
