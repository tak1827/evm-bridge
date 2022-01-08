package bridge

const (
	SlotERC20 = 1
	SlotNFT   = 2
)

type Rotator struct {
	size int
	slot int
}

func NewRotator(size int) Rotator {
	return Rotator{
		size: size,
		slot: 1,
	}
}

func (r *Rotator) Rotate() int {
	r.slot++
	if r.size < r.slot {
		r.slot = 1
	}
	return r.slot
}
