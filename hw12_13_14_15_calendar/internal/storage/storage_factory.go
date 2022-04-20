package storage

import (
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
	memstore "github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	dbstore "github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
)

type DataHandler interface {
	CreateEvent(event entity.Event) error
	UpdateEvent(eventID int, event entity.Event) error
	RemoveEvent(eventID int) error
	GetEventsByDate(eventDate string) ([]entity.Event, error)
	GetEventsByDateInterval(beginDate string, endDate string) ([]entity.Event, error)
}

func New(config config.Config) DataHandler {
	var storage DataHandler
	switch config.Store {
	case "sql":
		storage = dbstore.New(config.Connection)
	case "memory":
		storage = memstore.New()
	default:
		panic("storage type is not configured")
	}

	return storage
}
