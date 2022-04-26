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
	EventsListDay(eventDate string) ([]entity.Event, error)
	EventsListWeek(weekBegin string) ([]entity.Event, error)
	EventsListMonth(monthBegin string) ([]entity.Event, error)
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
	mux.Handle("/create-event", loggingMiddleware(handler.Create(*s), *s))
	mux.Handle("/update-event", loggingMiddleware(handler.Update(*s), *s))
	mux.Handle("/remove-event", loggingMiddleware(handler.Remove(*s), *s))
	mux.Handle("/events-list-day", loggingMiddleware(handler.EventsListDay(*s), *s))
	mux.Handle("/events-list-week", loggingMiddleware(handler.EventsListWeek(*s), *s))
	mux.Handle("/events-list-month", loggingMiddleware(handler.EventsListMonth(*s), *s))

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

func (h *RequestHandler) EventsListDay(s Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventDate := r.URL.Query().Get("event_date")
		result, err := s.app.EventsListDay(eventDate)
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

func (h *RequestHandler) EventsListWeek(s Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventDate := r.URL.Query().Get("event_date")
		result, err := s.app.EventsListWeek(eventDate)
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

func (h *RequestHandler) EventsListMonth(s Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventDate := r.URL.Query().Get("event_date")
		result, err := s.app.EventsListMonth(eventDate)
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
