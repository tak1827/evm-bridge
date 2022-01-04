package pb

import (
	"github.com/tak1827/go-store/store"
)

var (
	PREFIX_TCONIRMED_BLOCK = []byte(".confirmedblok")

	KEY_ERC20 = []byte(".erc20")
	KEY_COIN  = []byte(".coin") // native coin, like ETH

	blockStore *store.PrefixStore
)

func GetConfirmedBlockERC20(db store.Store) (b ConfirmedBlock, err error) {
	s := GetBlockStore(db)

	v, err := s.Get(KEY_ERC20)
	if err != nil {
		return
	}
	err = b.Unmarshal(v)
	return
}

func PutConfirmedBlockERC20(db store.Store, num uint64) error {
	var (
		s = GetBlockStore(db)
		m ConfirmedBlock
	)

	m.Number = num

	value, err := m.Marshal()
	if err != nil {
		return err
	}

	return s.Put(KEY_ERC20, value)
}

func GetBlockStore(db store.Store) *store.PrefixStore {
	if blockStore == nil {
		blockStore = store.NewPrefixStore(db, PREFIX_TCONIRMED_BLOCK)
	}
	return blockStore
}
