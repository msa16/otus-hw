package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app" //nolint:depguard
)

func worker(ctx context.Context, app *app.App, topic string) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			app.Logger.Info("worker closing")
			ticker.Stop()
			return
		case <-ticker.C:
			processEvents(ctx, app, topic)
		}

	}
}

func processEvents(ctx context.Context, app *app.App, topic string) {
	app.Logger.Info("worker processing events")

	events, err := app.Storage.ListEventsReminder(ctx)
	if err != nil {
		app.Logger.Error("failed to get events: " + err.Error())
		return
	}
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			app.Logger.Error("failed to marshal event: " + err.Error())
			continue
		}
		err = app.Broker.Publish(topic, payload)
		if err != nil {
			app.Logger.Error("failed to publish event: " + err.Error())
			continue
		}
		app.Logger.Info("published " + event.Title)
	}
}
