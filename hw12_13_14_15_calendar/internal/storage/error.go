package storage

import (
	"errors"
)

var (
	ErrDateBusy        = errors.New("time reserved for another event")
	ErrEventNotFound   = errors.New("event not found")
	ErrUpdateUserID    = errors.New("can't change user id")
	ErrInvalidStopTime = errors.New("stop time must be greater than start time")
)
