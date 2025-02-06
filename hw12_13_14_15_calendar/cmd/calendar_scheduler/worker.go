package main

import (
	"context"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app" //nolint:depguard
)

func worker(ctx context.Context, app *app.App) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			app.Logger.Info("worker closing")
			ticker.Stop()
			return
		case <-ticker.C:
		}
		app.Logger.Info("worker processing events")
	}
}
