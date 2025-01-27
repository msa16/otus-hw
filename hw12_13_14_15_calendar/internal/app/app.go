package app

import (
	"context"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type App struct {
	logger  Logger
	storage Storage
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
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
