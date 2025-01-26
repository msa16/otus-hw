package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
)

type Server struct {
	logger *logger.Logger
	server *http.Server
	app    Application
}

type Application interface { // TODO
}

func NewServer(logger *logger.Logger, app Application, host string, port int) *Server {
	server := &http.Server{
		Addr:        net.JoinHostPort(host, strconv.Itoa(port)),
		ReadTimeout: 30 * time.Second,
	}
	helloWorldHandler := http.HandlerFunc(helloWorld)
	http.Handle("/hello", loggingMiddleware(logger, helloWorldHandler))
	return &Server{logger: logger, server: server, app: app}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("starting http server at %s", s.server.Addr))

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
	s.logger.Info("shutting down http server")
	return s.server.Shutdown(ctx)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Calendar application welcome page. Current time %v\n", time.Now())))
}
