package app

import (
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
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

func (a *App) ListEvents(userID int) ([]entity.Event, error) {
	return a.storage.GetEventsList(userID)
}
