package http

import (
	"errors"
	"net/http"
)

type HTTPServer struct {
	router http.Handler
}

func NewHTTPServer(router http.Handler) *HTTPServer {
	return &HTTPServer{
		router: router,
	}
}

func (s *HTTPServer) Start(addr string) error {
	if err := http.ListenAndServe(addr, s.router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}
