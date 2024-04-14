package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler *echo.Echo) *Server {
	const port = "8080" // Константа для порту

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: handler, // Використання Router як Handler
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
