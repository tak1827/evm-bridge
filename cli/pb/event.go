package pb

import (
	"github.com/lithdew/bytesutil"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/go-store/store"
)

var (
	PREFIX_EVENT_ERC20 = []byte(".eventerc20")

	eventStore *store.PrefixStore
)

func (m *EventERC20Deposited) StoreKey() []byte {
	id := m.GetId()
	return bytesutil.AppendUint64BE(nil, id)
}

func GetEventERC20(db store.Store, id uint64) (m EventERC20Deposited, err error) {
	s := GetEventERC20Store(db)

	m.Id = id
	v, err := s.Get(m.StoreKey())
	if err != nil {
		return
	}
	err = m.Unmarshal(v)
	return
}

func (m *EventERC20Deposited) Put(db store.Store) error {
	s := GetEventERC20Store(db)

	value, err := m.Marshal()
	if err != nil {
		return err
	}

	return s.Put(m.StoreKey(), value)
}

func ToEventERC20Deposited(e *client.IBankERC20Deposited) EventERC20Deposited {
	return EventERC20Deposited{
		Id:     uint64(e.Id.Int64()),
		Token:  e.Token.Hex(),
		Sender: e.Sender.Hex(),
		Amount: e.Amount.String(),
		Retry:  0,
		Status: EventStatus_UNDEFINED,
	}
}

func GetEventERC20Store(db store.Store) *store.PrefixStore {
	if eventStore == nil {
		eventStore = store.NewPrefixStore(db, PREFIX_EVENT_ERC20)
	}
	return eventStore
}
