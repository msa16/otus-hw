package main

import (
	"context"
	"encoding/json"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"     //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

var (
	calendar  *app.App
	workerCtx context.Context
)

func processNotification(raw *[]byte) error {
	notification := storage.Notification{}
	err := json.Unmarshal(*raw, &notification)
	if err != nil {
		// When a handler returns an error, the default behavior is to send a Nack (negative-acknowledgement).
		// The message will be processed again.
		// если не смогли разобрать что прилетело, то нет смысла получать это снова
		calendar.Logger.Error("unmarshal " + err.Error())
		return nil //nolint:nilerr
	}

	calendar.Logger.Info("received event " + notification.ID)
	err = calendar.Storage.SaveNotification(workerCtx, notification)
	if err != nil {
		calendar.Logger.Error("failed to save notification: " + err.Error())
		return err
	}
	return nil
}
