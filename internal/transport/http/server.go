package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
	logger *logrus.Logger
}

func NewServer(addr string, handler http.Handler, logger *logrus.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,

			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting server at (", s.server.Addr, ")")

	err := s.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
