package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

const (
	badID = "bad_id"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	repo := New()
	event1 := storage.Event{
		Title:       "title 1",
		StartTime:   time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC),
		StopTime:    time.Date(2025, 1, 1, 11, 30, 0, 0, time.UTC),
		Description: "description 1",
		UserID:      "1"}
	event2 := storage.Event{
		Title:       "title 2",
		StartTime:   time.Date(2025, 1, 2, 11, 10, 0, 0, time.UTC),
		StopTime:    time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC),
		Description: "description 2",
		UserID:      "1"}

	t.Run("add event 1", func(t *testing.T) {
		id, err := repo.CreateEvent(ctx, event1)
		event1.ID = id

		require.Equal(t, len(id), 36, "generated id must be 36 symbols")
		require.NoError(t, err)
		require.Equal(t, len(repo.all), 1)
		require.Equal(t, len(repo.byUser), 1)
		require.Equal(t, len(repo.byUser[event1.UserID]), 1)
		require.Equal(t, repo.byUser[event1.UserID][event1.StartTime], &event1)
		require.Equal(t, repo.all[id], &event1)
	})

	t.Run("add event ErrDateBusy", func(t *testing.T) {
		_, err := repo.CreateEvent(ctx, event1)
		require.ErrorIs(t, err, storage.ErrDateBusy)
	})

	t.Run("add event ErrInvalidStopTime", func(t *testing.T) {
		savedStopTime := event1.StopTime
		event1.StopTime = time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		_, err := repo.CreateEvent(ctx, event1)
		event1.StopTime = savedStopTime
		require.ErrorIs(t, err, storage.ErrInvalidStopTime)
	})

	t.Run("add event 2", func(t *testing.T) {
		id, err := repo.CreateEvent(ctx, event2)
		event2.ID = id

		require.NoError(t, err)
		require.Equal(t, len(repo.all), 2)
		require.Equal(t, len(repo.byUser), 1)
		require.Equal(t, len(repo.byUser[event2.UserID]), 2)
		require.Equal(t, repo.byUser[event2.UserID][event2.StartTime], &event2)
		require.Equal(t, repo.all[id], &event2)
	})

	t.Run("get event 1", func(t *testing.T) {
		event, err := repo.GetEvent(ctx, event1.ID)
		require.NoError(t, err)
		require.Equal(t, event, &event1)
	})
	t.Run("get event 2", func(t *testing.T) {
		event, err := repo.GetEvent(ctx, event2.ID)
		require.NoError(t, err)
		require.Equal(t, event, &event2)
	})
	t.Run("get event ErrEventNotFound", func(t *testing.T) {
		_, err := repo.GetEvent(ctx, badID)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
	t.Run("update ErrEventNotFound", func(t *testing.T) {
		savedID := event1.ID
		event1.ID = badID
		err := repo.UpdateEvent(ctx, event1)
		event1.ID = savedID
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
	t.Run("update ErrUpdateUserID", func(t *testing.T) {
		savedUserID := event1.UserID
		event1.UserID = badID
		err := repo.UpdateEvent(ctx, event1)
		event1.UserID = savedUserID
		require.ErrorIs(t, err, storage.ErrUpdateUserID)
	})
	t.Run("update ErrInvalidStopTime", func(t *testing.T) {
		savedStopTime := event1.StopTime
		event1.StopTime = time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		err := repo.UpdateEvent(ctx, event1)
		event1.StopTime = savedStopTime
		require.ErrorIs(t, err, storage.ErrInvalidStopTime)
	})
	t.Run("update ErrDateBusy", func(t *testing.T) {
		savedID := event1.ID
		event1.ID = event2.ID
		err := repo.UpdateEvent(ctx, event1)
		event1.ID = savedID
		require.ErrorIs(t, err, storage.ErrDateBusy)
	})
	t.Run("update event ok", func(t *testing.T) {
		err := repo.UpdateEvent(ctx, event1)
		require.NoError(t, err)
		require.Equal(t, len(repo.all), 2)
		require.Equal(t, len(repo.byUser), 1)
		require.Equal(t, len(repo.byUser[event1.UserID]), 2)
		require.Equal(t, repo.byUser[event1.UserID][event1.StartTime], &event1)
		require.Equal(t, repo.all[event1.ID], &event1)
	})

	t.Run("list events day", func(t *testing.T) {
		events, err := repo.ListEventsDay(ctx, event1.StartTime)
		require.NoError(t, err)
		require.Equal(t, len(events), 1)
		require.Equal(t, events[0], event1)
	})
	t.Run("list events week", func(t *testing.T) {
		events, err := repo.ListEventsWeek(ctx, event1.StartTime)
		require.NoError(t, err)
		require.Equal(t, len(events), 2)
	})
	t.Run("list events month", func(t *testing.T) {
		events, err := repo.ListEventsMonth(ctx, event1.StartTime)
		require.NoError(t, err)
		require.Equal(t, len(events), 2)
	})

	t.Run("delete event ErrEventNotFound", func(t *testing.T) {
		err := repo.DeleteEvent(ctx, badID)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
	t.Run("delete event ok", func(t *testing.T) {
		err := repo.DeleteEvent(ctx, event1.ID)
		require.NoError(t, err)
		require.Equal(t, len(repo.all), 1)
		require.Equal(t, len(repo.byUser), 1)
		require.Equal(t, len(repo.byUser[event2.UserID]), 1)
		require.Equal(t, repo.byUser[event2.UserID][event2.StartTime], &event2)
		require.Equal(t, repo.all[event2.ID], &event2)
	})

}
