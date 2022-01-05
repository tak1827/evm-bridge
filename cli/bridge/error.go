package bridge

import (
	"errors"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrPairNotFound  = errors.New("pair not found")
)
