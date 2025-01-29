//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml api.yaml

package api

import (
	"net/http"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) FindEvents(w http.ResponseWriter, r *http.Request, params FindEventsParams) {
}

func (Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
}

func (Server) DeleteEventById(w http.ResponseWriter, r *http.Request, id string) {
}

func (Server) FindEventById(w http.ResponseWriter, r *http.Request, id string) {
}

func (Server) UpdateEventById(w http.ResponseWriter, r *http.Request, id string) {
}
