package pb

import (
	"github.com/tak1827/go-store/store"
)

var (
	PREFIX_ADDR_PAIR = []byte(".pair")

	pairStore *store.PrefixStore
)

func (m *Pair) StoreKey() []byte {
	return []byte(m.GetInaddr())
}

func GetPair(db store.Store, inaddr string) (m Pair, err error) {
	s := GetPairStore(db)

	m.Inaddr = inaddr

	value, err := s.Get(m.StoreKey())
	if err != nil {
		return
	}
	err = m.Unmarshal(value)
	return
}

func (m *Pair) Put(db store.Store) error {
	s := GetPairStore(db)

	value, err := m.Marshal()
	if err != nil {
		return err
	}

	return s.Put(m.StoreKey(), value)
}

func GetPairStore(db store.Store) *store.PrefixStore {
	if pairStore == nil {
		pairStore = store.NewPrefixStore(db, PREFIX_ADDR_PAIR)
	}
	return pairStore
}
