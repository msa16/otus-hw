//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml api.yaml

package api

import (
	"encoding/json"
	"net/http"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"
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

}

func (s *ApiServer) CreateEvent(w http.ResponseWriter, r *http.Request) {
}

func (s *ApiServer) DeleteEventById(w http.ResponseWriter, r *http.Request, id string) {
}

func (s *ApiServer) FindEventById(w http.ResponseWriter, r *http.Request, id string) {
	event, err := s.app.Storage.GetEvent(r.Context(), id)
	if err != nil {
		sendApiError(w, http.StatusNotFound, err.Error())
		return
	}

	resp := Event{
		ID:    event.ID,
		Title: event.Title,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *ApiServer) UpdateEventById(w http.ResponseWriter, r *http.Request, id string) {
}

func sendApiError(w http.ResponseWriter, code int, message string) {
	apiErr := Error{
		Code:    code,
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(apiErr)
}
