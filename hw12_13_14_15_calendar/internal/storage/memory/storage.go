package storage

import (
	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
)

type MemoryStorage struct {
	eventIncr int
	EventList map[int]entity.Event
}

func New() MemoryStorage {
	storage := MemoryStorage{
		eventIncr: 0,
		EventList: make(map[int]entity.Event),
	}
	return storage
}

func (s MemoryStorage) CreateEvent(event entity.Event) error {
	s.eventIncr++
	event.ID = s.eventIncr
	s.EventList[s.eventIncr] = event

	return nil
}

func (s MemoryStorage) UpdateEvent(eventID int, event entity.Event) error {
	s.EventList[eventID] = event
	return nil
}

func (s MemoryStorage) RemoveEvent(eventID int) error {
	delete(s.EventList, eventID)
	return nil
}

func (s MemoryStorage) GetEventsByDate(eventDate string) ([]entity.Event, error) {
	result := make([]entity.Event, 0)

	for _, event := range s.EventList {
		if event.CreatedAt != eventDate {
			continue
		}

		result = append(result, event)
	}

	return result, nil
}

func (s MemoryStorage) GetEventsByDateInterval(beginDate string, endDate string) ([]entity.Event, error) {
	result := make([]entity.Event, 0)

	for _, event := range s.EventList {
		if event.CreatedAt >= beginDate && event.CreatedAt <= endDate {
			result = append(result, event)
		}
	}

	return result, nil
}
