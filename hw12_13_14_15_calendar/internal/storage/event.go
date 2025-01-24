package storage

import (
	"context"
	"time"
)

type Event struct {
	ID          string
	Title       string
	StartTime   time.Time
	StopTime    time.Time
	Description string
	UserID      string
	Reminder    time.Duration
}

type EventRepo interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (Event, error)
	ListEventsDay(ctx context.Context, StartTime time.Time) ([]Event, error)
	ListEventsWeek(ctx context.Context, StartTime time.Time) ([]Event, error)
	ListEventsMonth(ctx context.Context, StartTime time.Time) ([]Event, error)
}
