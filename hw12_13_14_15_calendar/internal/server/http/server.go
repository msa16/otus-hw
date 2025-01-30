package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app" //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/server/http/api"
)

type Server struct {
	server *http.Server
	app    *app.App
}

func NewServer(app *app.App, host string, port int) *Server {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	apiServer := api.NewApiServer(app)
	mux := http.NewServeMux()
	// get an `http.Handler` that we can use
	h := loggingMiddleware(app.Logger, api.HandlerFromMux(apiServer, mux))

	server := &http.Server{
		Addr:        net.JoinHostPort(host, strconv.Itoa(port)),
		ReadTimeout: 30 * time.Second,
		Handler:     h,
	}
	helloWorldHandler := http.HandlerFunc(helloWorld)
	mux.Handle("/hello", loggingMiddleware(app.Logger, helloWorldHandler))
	return &Server{server: server, app: app}
}

func (s *Server) Start(ctx context.Context) error {
	s.app.Logger.Info(fmt.Sprintf("starting http server at %s", s.server.Addr))

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	s.app.Logger.Info("shutting down http server")
	return s.server.Shutdown(ctx)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Calendar application welcome page. Current time %v\n", time.Now())))
}
