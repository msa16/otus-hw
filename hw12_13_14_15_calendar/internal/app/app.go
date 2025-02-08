package app

import (
	"context"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/client"  //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type App struct {
	Logger  Logger
	Storage Storage
	Broker  client.Broker
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) (string, error)
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	ListEventsDay(ctx context.Context, startTime time.Time) ([]*storage.Event, error)
	ListEventsWeek(ctx context.Context, startTime time.Time) ([]*storage.Event, error)
	ListEventsMonth(ctx context.Context, startTime time.Time) ([]*storage.Event, error)
	ListEventsReminder(ctx context.Context) ([]*storage.Event, error)
	ClearReminderTime(ctx context.Context, id string) error
	DeleteEventsBeforeDate(ctx context.Context, time time.Time) error
}

func New(logger Logger, storage Storage, broker client.Broker) *App {
	return &App{Logger: logger, Storage: storage, Broker: broker}
}
