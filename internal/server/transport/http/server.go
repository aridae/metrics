package http

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/pkg/logger"
	"net/http"
)

type Server struct {
	server  *http.Server
	address string
}

func NewServer(address string, mux http.Handler, mws ...func(http.Handler) http.Handler) *Server {
	for _, mw := range mws {
		mux = mw(mux)
	}

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return &Server{address: address, server: server}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		err := s.server.Shutdown(ctx)
		if err != nil {
			logger.Errorf("error shutting down http server: %v", err)
		}
	}()

	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}
