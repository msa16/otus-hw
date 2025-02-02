//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml api.yaml

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"     //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check.
var _ ServerInterface = (*Server)(nil)

type Server struct {
	app *app.App
}

func NewAPIServer(app *app.App) *Server {
	return &Server{app: app}
}

func (s *Server) FindEvents(w http.ResponseWriter, r *http.Request, params FindEventsParams) {
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
		sendAPIError(w, storageErrorToAPIErrorCode(err), err.Error())
		return
	}

	result := make([]Event, 0, len(stEvents))
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

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// We expect a NewEvent object in the request body.
	var newEvent NewEvent
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		sendAPIError(w, http.StatusBadRequest, "Invalid format for NewEvent")
		return
	}
	reminder, err := time.ParseDuration(*newEvent.Reminder)
	if err != nil {
		sendAPIError(w, http.StatusBadRequest, "Invalid format for Reminder")
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
		sendAPIError(w, storageErrorToAPIErrorCode(err), err.Error())
		return
	}
	eventID := EventID{ID: id}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(eventID)
}

func (s *Server) DeleteEventByID(w http.ResponseWriter, r *http.Request, id string) {
	err := s.app.Storage.DeleteEvent(r.Context(), id)
	if err != nil {
		sendAPIError(w, storageErrorToAPIErrorCode(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) FindEventByID(w http.ResponseWriter, r *http.Request, id string) {
	stEvent, err := s.app.Storage.GetEvent(r.Context(), id)
	if err != nil {
		sendAPIError(w, storageErrorToAPIErrorCode(err), err.Error())
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

func (s *Server) UpdateEventByID(w http.ResponseWriter, r *http.Request, id string) {
	// We expect a Event object in the request body.
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		sendAPIError(w, http.StatusBadRequest, "Invalid format for Event")
		return
	}
	reminder, err := time.ParseDuration(*event.Reminder)
	if err != nil {
		sendAPIError(w, http.StatusBadRequest, "Invalid format for Reminder")
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
		sendAPIError(w, storageErrorToAPIErrorCode(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func sendAPIError(w http.ResponseWriter, code int, message string) {
	apiErr := Error{
		Code:    code,
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(apiErr)
}

func storageErrorToAPIErrorCode(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, storage.ErrDateBusy) ||
		errors.Is(err, storage.ErrInvalidArgiments) ||
		errors.Is(err, storage.ErrUpdateUserID) ||
		errors.Is(err, storage.ErrInvalidStopTime):
		return http.StatusBadRequest
	case errors.Is(err, storage.ErrCreateEvent) ||
		errors.Is(err, storage.ErrUpdateEvent) ||
		errors.Is(err, storage.ErrDeleteEvent) ||
		errors.Is(err, storage.ErrReadEvent):
		return http.StatusInternalServerError
	case errors.Is(err, storage.ErrEventNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
