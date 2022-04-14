package storage

import (
	"testing"
	"time"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	storage := New()

	event := entity.Event{
		UserID:       1,
		Title:        "first",
		EventTime:    time.Time{},
		DurationMin:  120,
		NoticeBefore: 60,
		Description:  "first event",
		CreatedAt:    time.Time{},
	}

	err := storage.CreateEvent(event)
	require.NoError(t, err)

	eventID := 1
	event = entity.Event{
		UserID:       1,
		Title:        "second",
		DurationMin:  180,
		NoticeBefore: 30,
	}
	err = storage.UpdateEvent(eventID, event)
	require.NoError(t, err)

	event = storage.EventList[eventID]
	require.Equal(t, event.UserID, 1)
	require.Equal(t, event.Title, "second")
	require.Equal(t, event.DurationMin, 180)
	require.Equal(t, event.NoticeBefore, 30)

	userID := 1
	eventList, err := storage.GetEventsList(userID)
	require.NoError(t, err)
	require.Equal(t, 1, len(eventList))

	err = storage.RemoveEvent(eventID)
	require.NoError(t, err)

	eventList, err = storage.GetEventsList(userID)
	require.NoError(t, err)
	require.Equal(t, 0, len(eventList))
}
