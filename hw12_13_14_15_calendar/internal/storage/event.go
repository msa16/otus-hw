package storage

import (
	"time"
)

type Event struct {
	ID          string
	Title       string
	StartTime   time.Time
	StopTime    time.Time
	Description string
	UserID      int64
	Reminder    time.Duration
}
