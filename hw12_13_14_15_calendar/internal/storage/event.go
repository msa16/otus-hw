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
	UserID      int64
	Reminder    time.Duration
}

type EventRepo interface {
	CreateEvent(ctx context.Context, event Event) (string, error)
	// по ТЗ Обновить (ID события, событие)
	UpdateEvent(ctx context.Context, id string, event Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (Event, error)
	ListEventsDay(ctx context.Context, startTime time.Time) ([]Event, error)
	ListEventsWeek(ctx context.Context, startTime time.Time) ([]Event, error)
	ListEventsMonth(ctx context.Context, startTime time.Time) ([]Event, error)
}
