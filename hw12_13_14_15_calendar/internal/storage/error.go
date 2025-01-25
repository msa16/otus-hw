package storage

import (
	"errors"
)

var (
	ErrDateBusy         = errors.New("time reserved for another event")
	ErrEventNotFound    = errors.New("event not found")
	ErrInvalidArgiments = errors.New("invalid arguments")
	ErrUpdateUserID     = errors.New("can't change user id")
	ErrInvalidStopTime  = errors.New("stop time must be greater than start time")
	ErrCreateEvent      = errors.New("can't create event")
	ErrUpdateEvent      = errors.New("can't update event")
	ErrDeleteEvent      = errors.New("can't delete event")
)
