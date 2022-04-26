package app

import (
	"time"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
)

const (
	weekDays  = 7
	monthDays = 31
)

type App struct {
	logger  logger.Logger
	storage storage.DataHandler
}

func New(logger *logger.Logger, storage storage.DataHandler) *App {
	return &App{
		logger:  *logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(event entity.Event) error {
	return a.storage.CreateEvent(event)
}

func (a *App) UpdateEvent(eventID int, event entity.Event) error {
	return a.storage.UpdateEvent(eventID, event)
}

func (a *App) RemoveEvent(eventID int) error {
	return a.storage.RemoveEvent(eventID)
}

func (a *App) EventsListDay(eventDate string) ([]entity.Event, error) {
	return a.storage.GetEventsByDate(eventDate)
}

func (a *App) EventsListWeek(weekBegin string) ([]entity.Event, error) {
	weekEnd, err := time.Parse("2006-01-02", weekBegin)
	if err != nil {
		return nil, err
	}

	return a.storage.GetEventsByDateInterval(weekBegin, weekEnd.Add(weekDays).Format("2006-01-02"))
}

func (a *App) EventsListMonth(monthBegin string) ([]entity.Event, error) {
	monthEnd, err := time.Parse("2006-01-02", monthBegin)
	if err != nil {
		return nil, err
	}

	return a.storage.GetEventsByDateInterval(monthBegin, monthEnd.Add(monthDays).Format("2006-01-02"))
}
