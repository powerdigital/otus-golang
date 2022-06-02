package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/powerdigital/project/internal/app"
)

type Server struct {
	app        app.App
	httpServer *http.Server
}

func NewServer(app app.App) *Server {
	return &Server{
		app: app,
	}
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:         ":8888",
		Handler:      s.getHandler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		s.app.Logger.Error().Err(err).Msg("http server error")
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.app.Logger.Error().Err(err).Msg("http server shutdown error")
		return fmt.Errorf("http server shutdown error: %w", err)
	}

	return nil
}

func (s *Server) getHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", s.ResizeImage())

	return mux
}

func (s Server) ResizeImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.app.ResizeImage(w, r)
	})
}
