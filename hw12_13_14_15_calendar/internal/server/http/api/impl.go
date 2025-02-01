//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml api.yaml

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*ApiServer)(nil)

type ApiServer struct {
	app *app.App
}

func NewApiServer(app *app.App) *ApiServer {
	return &ApiServer{app: app}
}

func (s *ApiServer) FindEvents(w http.ResponseWriter, r *http.Request, params FindEventsParams) {
	var stEvents []*storage.Event
	var err error
	switch *params.Period {
	case "day":
		stEvents, err = s.app.Storage.ListEventsDay(r.Context(), params.StartTime)
	case "week":
		stEvents, err = s.app.Storage.ListEventsWeek(r.Context(), params.StartTime)
	case "month":
		stEvents, err = s.app.Storage.ListEventsMonth(r.Context(), params.StartTime)
	}
	if err != nil {
		sendApiError(w, storageErrorToApiErrorCode(err), err.Error())
		return
	}

	var result []Event
	for _, stEvent := range stEvents {
		reminder := stEvent.Reminder.String()
		event := Event{
			ID:          stEvent.ID,
			Title:       stEvent.Title,
			UserID:      stEvent.UserID,
			StartTime:   stEvent.StartTime,
			StopTime:    stEvent.StopTime,
			Description: &stEvent.Description,
			Reminder:    &reminder,
		}
		result = append(result, event)
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func (s *ApiServer) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// We expect a NewEvent object in the request body.
	var newEvent NewEvent
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for NewEvent")
		return
	}
	reminder, err := time.ParseDuration(*newEvent.Reminder)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for Reminder")
		return
	}

	storageEvent := storage.Event{
		Title:     newEvent.Title,
		StartTime: newEvent.StartTime,
		StopTime:  newEvent.StopTime,
		UserID:    newEvent.UserID,
		Reminder:  reminder,
	}
	if newEvent.Description != nil {
		storageEvent.Description = *newEvent.Description
	}

	id, err := s.app.Storage.CreateEvent(r.Context(), storageEvent)
	if err != nil {
		sendApiError(w, storageErrorToApiErrorCode(err), err.Error())
		return
	}
	eventID := EventID{ID: id}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(eventID)
}

func (s *ApiServer) DeleteEventById(w http.ResponseWriter, r *http.Request, id string) {
	err := s.app.Storage.DeleteEvent(r.Context(), id)
	if err != nil {
		sendApiError(w, storageErrorToApiErrorCode(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *ApiServer) FindEventById(w http.ResponseWriter, r *http.Request, id string) {
	stEvent, err := s.app.Storage.GetEvent(r.Context(), id)
	if err != nil {
		sendApiError(w, storageErrorToApiErrorCode(err), err.Error())
		return
	}

	reminder := stEvent.Reminder.String()
	resp := Event{
		ID:          stEvent.ID,
		Title:       stEvent.Title,
		UserID:      stEvent.UserID,
		StartTime:   stEvent.StartTime,
		StopTime:    stEvent.StopTime,
		Description: &stEvent.Description,
		Reminder:    &reminder,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *ApiServer) UpdateEventById(w http.ResponseWriter, r *http.Request, id string) {
	// We expect a Event object in the request body.
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for Event")
		return
	}
	reminder, err := time.ParseDuration(*event.Reminder)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for Reminder")
		return
	}

	err = s.app.Storage.UpdateEvent(r.Context(), id, storage.Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: *event.Description,
		StartTime:   event.StartTime,
		StopTime:    event.StopTime,
		UserID:      event.UserID,
		Reminder:    reminder,
	})
	if err != nil {
		sendApiError(w, storageErrorToApiErrorCode(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func sendApiError(w http.ResponseWriter, code int, message string) {
	apiErr := Error{
		Code:    code,
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(apiErr)
}

func storageErrorToApiErrorCode(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	case storage.ErrDateBusy, storage.ErrInvalidArgiments, storage.ErrUpdateUserID, storage.ErrInvalidStopTime:
		return http.StatusBadRequest
	case storage.ErrCreateEvent, storage.ErrUpdateEvent, storage.ErrDeleteEvent, storage.ErrReadEvent:
		return http.StatusInternalServerError
	case storage.ErrEventNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
