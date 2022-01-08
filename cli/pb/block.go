package pb

import (
	"errors"

	"github.com/tak1827/go-store/store"
)

var (
	PREFIX_CONFIRMED_BLOCK = []byte(".confirmedblok")

	KEY_ERC20 = []byte(".erc20")
	KEY_COIN  = []byte(".coin") // native coin, like ETH
	KEY_NFT   = []byte(".nft")

	blockStore *store.PrefixStore
)

type BlockType int

const (
	BlockERC20 BlockType = iota
	BlockNFT
)

func GetConfirmedBlock(db store.Store, t BlockType) (b ConfirmedBlock, err error) {
	var (
		s = getBlockStore(db)
		v []byte
	)
	switch t {
	case BlockERC20:
		v, err = s.Get(KEY_ERC20)
	case BlockNFT:
		v, err = s.Get(KEY_NFT)
	}

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			err = nil
		}
		return
	}
	err = b.Unmarshal(v)
	return
}

func (m ConfirmedBlock) Put(db store.Store, t BlockType) error {
	s := getBlockStore(db)
	value, err := m.Marshal()
	if err != nil {
		return err
	}

	switch t {
	case BlockERC20:
		s.Put(KEY_ERC20, value)
	case BlockNFT:
		s.Put(KEY_NFT, value)
	}

	return nil
}

func getBlockStore(db store.Store) *store.PrefixStore {
	if blockStore == nil {
		blockStore = store.NewPrefixStore(db, PREFIX_CONFIRMED_BLOCK)
	}
	return blockStore
}
