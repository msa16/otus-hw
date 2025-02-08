package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"     //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
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
	app.Logger.Debug("worker processing events")

	events, err := app.Storage.ListEventsReminder(ctx)
	if err != nil {
		app.Logger.Error("failed to get events: " + err.Error())
		return
	}
	for _, event := range events {
		notification := storage.Notification{
			ID:        event.ID,
			Title:     event.Title,
			StartTime: event.StartTime,
			UserID:    event.UserID,
		}
		payload, err := json.Marshal(notification)
		if err != nil {
			app.Logger.Error("failed to marshal notification: " + err.Error())
			continue
		}
		err = app.Broker.Publish(topic, payload)
		if err != nil {
			app.Logger.Error("failed to publish notification: " + err.Error())
			continue
		}
		err = app.Storage.ClearReminderTime(ctx, event.ID)
		if err != nil {
			app.Logger.Error("failed to clear reminder time: " + err.Error())
			continue
		}
		app.Logger.Info("published " + event.Title)
	}
}
