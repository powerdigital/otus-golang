package storage

import (
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
)

type MemoryStorage struct{}

func New() MemoryStorage {
	storage := MemoryStorage{}
	return storage
}

func (s MemoryStorage) CreateEvent(event entity.Event) error {
	return nil
}

func (s MemoryStorage) UpdateEvent(eventID int, event entity.Event) error {
	return nil
}

func (s MemoryStorage) RemoveEvent(eventID int) error {
	return nil
}

func (s MemoryStorage) GetEventsList(userID int) ([]entity.Event, error) {
	var result []entity.Event
	return result, nil
}
