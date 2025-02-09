package main

import (
	"encoding/json"
	"strconv"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"     //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

var calendar *app.App

func processNotification(raw *[]byte) error {
	consumedPayload := storage.Notification{}
	err := json.Unmarshal(*raw, &consumedPayload)
	if err != nil {
		// When a handler returns an error, the default behavior is to send a Nack (negative-acknowledgement).
		// The message will be processed again.
		// если не смогли разобрать что прилетело, то нет смысла получать это снова
		calendar.Logger.Error("unmarshal " + err.Error())
		return nil //nolint:nilerr
	}

	calendar.Logger.Info("received event " + consumedPayload.Title + " " +
		consumedPayload.ID + " " + consumedPayload.StartTime.String() + " " + strconv.Itoa(int(consumedPayload.UserID)))
	return nil
}
