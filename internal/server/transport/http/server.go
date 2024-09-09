package http

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	address string
	server  *http.Server
}

func NewServer(address string, mux http.Handler) *Server {
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return &Server{address: address, server: server}
}

func (s *Server) Run(_ context.Context) error {
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
