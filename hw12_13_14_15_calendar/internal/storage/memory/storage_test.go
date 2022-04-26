package storage

import (
	"testing"
	"time"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	storage := New()

	eventOriginal := entity.Event{
		UserID:       1,
		Title:        "first",
		EventTime:    time.Now().Format("2006-01-02"),
		DurationMin:  120,
		NoticeBefore: 60,
		Description:  "first event",
		CreatedAt:    time.Now().Format("2006-01-02"),
	}

	err := storage.CreateEvent(eventOriginal)
	require.NoError(t, err)

	eventID := 1
	eventUpdated := entity.Event{
		UserID:       1,
		Title:        "second",
		EventTime:    time.Now().Format("2006-01-02"),
		DurationMin:  180,
		NoticeBefore: 30,
		Description:  "second event",
		CreatedAt:    time.Now().Format("2006-01-02"),
	}
	err = storage.UpdateEvent(eventID, eventUpdated)
	require.NoError(t, err)

	eventFromCache := storage.EventList[eventID]
	require.Equal(t, eventFromCache, eventUpdated)

	eventDate := time.Now().Format("2006-01-02")
	eventList, err := storage.GetEventsByDate(eventDate)
	require.NoError(t, err)
	require.Equal(t, 1, len(eventList))

	eventList, err = storage.GetEventsByDateInterval(eventDate, eventDate)
	require.NoError(t, err)
	require.Equal(t, 1, len(eventList))

	err = storage.RemoveEvent(eventID)
	require.NoError(t, err)

	eventList, err = storage.GetEventsByDate(eventDate)
	require.NoError(t, err)
	require.Equal(t, 0, len(eventList))
}
