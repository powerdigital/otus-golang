package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
)

type Server struct {
	app    Application
	logger logger.Logger
}

type RequestHandler struct{}

type Application interface {
	CreateEvent(event entity.Event) error
	UpdateEvent(eventID int, event entity.Event) error
	RemoveEvent(eventID int) error
	ListEvents(userID int) ([]entity.Event, error)
}

func NewServer(logger logger.Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:         ":8888",
		Handler:      s.getHandler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx.Done()
	return nil
}

func (s *Server) getHandler() *http.ServeMux {
	handler := &RequestHandler{}

	mux := http.NewServeMux()
	mux.Handle("/list", loggingMiddleware(handler.List(*s), *s))
	mux.Handle("/create", loggingMiddleware(handler.Create(*s), *s))
	mux.Handle("/update", loggingMiddleware(handler.Update(*s), *s))
	mux.Handle("/remove", loggingMiddleware(handler.Remove(*s), *s))

	return mux
}

func (h *RequestHandler) Create(s Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var event entity.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}

		s.app.CreateEvent(event)
	})
}

func (h *RequestHandler) Update(s Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventID, err := strconv.Atoi(r.URL.Query().Get("event_id"))
		if err != nil {
			err := errors.New("required event_id param is not provided")
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}

		var event entity.Event
		err = json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}

		s.app.UpdateEvent(eventID, event)
	})
}

func (h *RequestHandler) Remove(s Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventID, err := strconv.Atoi(r.URL.Query().Get("event_id"))
		if err != nil {
			err := errors.New("required event_id param is not provided")
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}

		s.app.RemoveEvent(eventID)
	})
}

func (h *RequestHandler) List(s Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			err := errors.New("required user_id param is not provided")
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}

		result, err := s.app.ListEvents(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.logger.Error(err.Error())
			return
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.logger.Error(err.Error())
			return
		}

		w.Write(jsonData)
	})
}
