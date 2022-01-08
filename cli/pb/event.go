package pb

import (
	"github.com/lithdew/bytesutil"
	"github.com/tak1827/evm-bridge/cli/client"
	"github.com/tak1827/go-store/store"
)

var (
	PREFIX_EVENT_ERC20 = []byte(".eventerc20")
	PREFIX_EVENT_NFT   = []byte(".eventnft")

	eventStoreERC20 *store.PrefixStore
	eventStoreNFT   *store.PrefixStore

	_ Event = (*EventERC20Deposited)(nil)
	_ Event = (*EventNFTDeposited)(nil)
)

type Event interface {
	GetRetry() uint32
	SetRetry(retry uint32)
	GetStatus() EventStatus
	SetStatus(status EventStatus)
	GetToken() string
	Get(db store.Store) error
	Put(db store.Store) error
}

func (m *EventERC20Deposited) StoreKey() []byte {
	id := m.GetId()
	return bytesutil.AppendUint64BE(nil, id)
}

func (m *EventERC20Deposited) SetRetry(retry uint32) {
	m.Retry = retry
}

func (m *EventERC20Deposited) SetStatus(status EventStatus) {
	m.Status = status
}

func (m *EventERC20Deposited) Get(db store.Store) error {
	s := getEventERC20Store(db)
	v, err := s.Get(m.StoreKey())
	if err != nil {
		return err
	}
	return m.Unmarshal(v)
}

func (m *EventERC20Deposited) Put(db store.Store) error {
	s := getEventERC20Store(db)
	value, err := m.Marshal()
	if err != nil {
		return err
	}
	return s.Put(m.StoreKey(), value)
}

func ToEventERC20Deposited(e *client.IBankERC20Deposited) *EventERC20Deposited {
	return &EventERC20Deposited{
		Id:     uint64(e.Id.Int64()),
		Token:  e.Token.Hex(),
		Sender: e.Sender.Hex(),
		Amount: e.Amount.String(),
		Retry:  0,
		Status: EventStatus_UNDEFINED,
	}
}

func getEventERC20Store(db store.Store) *store.PrefixStore {
	if eventStoreERC20 == nil {
		eventStoreERC20 = store.NewPrefixStore(db, PREFIX_EVENT_ERC20)
	}
	return eventStoreERC20
}

func (m *EventNFTDeposited) StoreKey() []byte {
	id := m.GetId()
	return bytesutil.AppendUint64BE(nil, id)
}

func (m *EventNFTDeposited) SetRetry(retry uint32) {
	m.Retry = retry
}

func (m *EventNFTDeposited) SetStatus(status EventStatus) {
	m.Status = status
}

func (m *EventNFTDeposited) Get(db store.Store) error {
	s := getEventNFTStore(db)
	v, err := s.Get(m.StoreKey())
	if err != nil {
		return err
	}
	return m.Unmarshal(v)
}

func (m *EventNFTDeposited) Put(db store.Store) error {
	s := getEventNFTStore(db)
	value, err := m.Marshal()
	if err != nil {
		return err
	}
	return s.Put(m.StoreKey(), value)
}

func ToEventNFTDeposited(e *client.IBankNFTDeposited) *EventNFTDeposited {
	return &EventNFTDeposited{
		Id:      uint64(e.Id.Int64()),
		Token:   e.Token.Hex(),
		Sender:  e.Sender.Hex(),
		Tokenid: uint64(e.Tokenid.Int64()),
		Retry:   0,
		Status:  EventStatus_UNDEFINED,
	}
}

func getEventNFTStore(db store.Store) *store.PrefixStore {
	if eventStoreNFT == nil {
		eventStoreNFT = store.NewPrefixStore(db, PREFIX_EVENT_NFT)
	}
	return eventStoreNFT
}
