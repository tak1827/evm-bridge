package bridge

import (
	"github.com/tak1827/transaction-confirmer/confirm"
)

type Option interface {
	Apply(*Bridge) error
}

type CustomConfirmedHandler confirm.HashHandler

func (f CustomConfirmedHandler) Apply(b *Bridge) error {
	b.CustomConfirmedHandler = confirm.HashHandler(f)
	return nil
}
func WithCustomConfirmedHandler(f confirm.HashHandler) CustomConfirmedHandler {
	return CustomConfirmedHandler(f)
}

type CustomErrHandler confirm.ErrHandler

func (f CustomErrHandler) Apply(b *Bridge) error {
	b.CustomErrHandler = confirm.ErrHandler(f)
	return nil
}
func WithCustomErrHandler(f confirm.ErrHandler) CustomErrHandler {
	return CustomErrHandler(f)
}
