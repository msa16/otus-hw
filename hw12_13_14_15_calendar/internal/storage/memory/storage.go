package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

type userEvents map[time.Time]*storage.Event
type Storage struct {
	mu sync.RWMutex
	// все события
	all map[string]*storage.Event
	// события по пользователям и времени. Ограничение: для одного пользователя в один момент времени может начинаться только одно событие.
	// другие пересечения событий по времени считаем допустимым
	byUser map[string]userEvents
}

func New() *Storage {
	return &Storage{mu: sync.RWMutex{}, all: make(map[string]*storage.Event), byUser: make(map[string]userEvents)}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	if event.StopTime.Before(event.StartTime) {
		return "", storage.ErrInvalidStopTime
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	ue := s.byUser[event.UserID]
	if ue == nil {
		// такого пользователя нет
		ue = make(userEvents)
		s.byUser[event.UserID] = ue
	} else if ue[event.StartTime] != nil {
		return "", storage.ErrDateBusy
	}
	event.ID = uuid.New().String()
	ue[event.StartTime] = &event
	s.all[event.ID] = &event
	return event.ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// проверки
	current := s.all[event.ID]
	if current == nil {
		return storage.ErrEventNotFound
	}
	if current.UserID != event.UserID {
		return storage.ErrUpdateUserID
	}
	if event.StopTime.Before(event.StartTime) {
		return storage.ErrInvalidStopTime
	}

	ue := s.byUser[event.UserID]
	if ce, ok := ue[event.StartTime]; ok && ce.ID != event.ID {
		return storage.ErrDateBusy
	}
	// изменение
	delete(ue, current.StartTime)
	ue[event.StartTime] = current

	current.Title = event.Title
	current.Description = event.Description
	current.StartTime = event.StartTime
	current.StopTime = event.StopTime
	current.Reminder = event.Reminder

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// проверки
	current := s.all[id]
	if current == nil {
		return storage.ErrEventNotFound
	}
	// изменение
	delete(s.byUser[current.UserID], current.StartTime)
	delete(s.all, id)
	return nil
}

func (s *Storage) listEventsInt(StartTime time.Time, StopTime time.Time) []storage.Event {
	result := make([]storage.Event, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.all {
		if v.StartTime.Compare(StartTime) >= 0 && v.StartTime.Compare(StopTime) <= 0 {
			result = append(result, *v)
		}
	}
	return result
}

func (s *Storage) ListEventsDay(ctx context.Context, StartTime time.Time) ([]storage.Event, error) {
	return s.listEventsInt(StartTime, StartTime.Add(time.Hour*24)), nil
}

func (s *Storage) ListEventsWeek(ctx context.Context, StartTime time.Time) ([]storage.Event, error) {
	return s.listEventsInt(StartTime, StartTime.Add(time.Hour*24*7)), nil
}
func (s *Storage) ListEventsMonth(ctx context.Context, StartTime time.Time) ([]storage.Event, error) {
	return s.listEventsInt(StartTime, StartTime.AddDate(0, 1, 0)), nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	current := s.all[id]
	if current == nil {
		return nil, storage.ErrEventNotFound
	}
	return current, nil
}
