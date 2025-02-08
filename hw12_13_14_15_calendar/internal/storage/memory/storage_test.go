package memorystorage

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/stretchr/testify/require"                              //nolint:depguard
)

const (
	badEventID = "bad_event_id"
	badUserID  = 777
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	repo := New()
	event1 := storage.Event{
		Title:       "title 1",
		StartTime:   time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC),
		StopTime:    time.Date(2025, 1, 1, 11, 30, 0, 0, time.UTC),
		Description: "description 1",
		UserID:      1,
	}
	event2 := storage.Event{
		Title:       "title 2",
		StartTime:   time.Date(2025, 1, 2, 11, 10, 0, 0, time.UTC),
		StopTime:    time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC),
		Description: "description 2",
		UserID:      1,
	}

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
		_, err := repo.GetEvent(ctx, badEventID)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
	t.Run("update ErrEventNotFound", func(t *testing.T) {
		savedID := event1.ID
		event1.ID = badEventID
		err := repo.UpdateEvent(ctx, badEventID, event1)
		event1.ID = savedID
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
	t.Run("update ErrUpdateUserID", func(t *testing.T) {
		savedUserID := event1.UserID
		event1.UserID = badUserID
		err := repo.UpdateEvent(ctx, event1.ID, event1)
		event1.UserID = savedUserID
		require.ErrorIs(t, err, storage.ErrUpdateUserID)
	})
	t.Run("update ErrInvalidStopTime", func(t *testing.T) {
		savedStopTime := event1.StopTime
		event1.StopTime = time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		err := repo.UpdateEvent(ctx, event1.ID, event1)
		event1.StopTime = savedStopTime
		require.ErrorIs(t, err, storage.ErrInvalidStopTime)
	})
	t.Run("update ErrDateBusy", func(t *testing.T) {
		savedID := event1.ID
		event1.ID = event2.ID
		err := repo.UpdateEvent(ctx, event1.ID, event1)
		event1.ID = savedID
		require.ErrorIs(t, err, storage.ErrDateBusy)
	})
	t.Run("update event ok", func(t *testing.T) {
		err := repo.UpdateEvent(ctx, event1.ID, event1)
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
		require.Equal(t, *events[0], event1)
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
	t.Run("list events reminder", func(t *testing.T) {
		events, err := repo.ListEventsReminder(ctx)
		require.NoError(t, err)
		require.Equal(t, 0, len(events))
	})

	t.Run("delete event ErrEventNotFound", func(t *testing.T) {
		err := repo.DeleteEvent(ctx, badEventID)
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

func TestStorageConcurrency(t *testing.T) {
	const threadCount = 10
	const objectPerThread = 10000
	ctx := context.Background()
	repo := New()
	afterCreate := make(chan string)
	afterReadUpdate := make(chan string)

	var wgCreate sync.WaitGroup
	// создаем объекты в разных потоках
	for i := 0; i < threadCount; i++ {
		wgCreate.Add(1)
		go func(userID int64) {
			defer wgCreate.Done()
			startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			stopTime := startTime.Add(time.Hour)
			for i := 0; i < objectPerThread; i++ {
				id, err := repo.CreateEvent(ctx, storage.Event{UserID: userID, StartTime: startTime, StopTime: stopTime})
				require.NoError(t, err)
				require.NotEmpty(t, id)
				afterCreate <- id
				startTime = stopTime
				stopTime = startTime.Add(time.Hour)
			}
		}(int64(i))
	}

	go func() {
		wgCreate.Wait()
		close(afterCreate)
	}()

	var readUpdateCount uint64
	var wgReadUpdate sync.WaitGroup
	// читаем/меняем объекты в разных потоках
	for i := 0; i < threadCount; i++ {
		wgReadUpdate.Add(1)
		go func() {
			defer wgReadUpdate.Done()
			for id := range afterCreate {
				event, err := repo.GetEvent(ctx, id)
				require.NoError(t, err)
				require.NotNil(t, event)
				require.Equal(t, id, event.ID)
				atomic.AddUint64(&readUpdateCount, 1)

				event.Description = event.ID
				err = repo.UpdateEvent(ctx, id, *event)
				require.NoError(t, err)
				afterReadUpdate <- id
			}
		}()
	}

	go func() {
		wgReadUpdate.Wait()
		close(afterReadUpdate)
	}()

	var readDeleteCount uint64
	var wgReadDelete sync.WaitGroup
	// читаем и удаляем объекты в разных потоках
	for i := 0; i < threadCount; i++ {
		wgReadDelete.Add(1)
		go func() {
			defer wgReadDelete.Done()
			for id := range afterReadUpdate {
				event, err := repo.GetEvent(ctx, id)
				require.NoError(t, err)
				require.NotNil(t, event)
				require.Equal(t, id, event.ID)
				// проверяем что update работает
				require.Equal(t, id, event.Description)

				err2 := repo.DeleteEvent(ctx, id)
				require.NoError(t, err2)
				atomic.AddUint64(&readDeleteCount, 1)
			}
		}()
	}
	wgReadDelete.Wait()

	require.Equal(t, uint64(threadCount*objectPerThread), readUpdateCount)
	require.Equal(t, uint64(threadCount*objectPerThread), readDeleteCount)
}
