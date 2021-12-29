package pb

var (
	// PREFIX_CONFIRMED_BLOCK = []byte(".confirmedblok")
	storeKey = []byte(".confirmedblok")
)

func (x *ConfirmedBlock) StoreKey() []byte {
	return storeKey
}
