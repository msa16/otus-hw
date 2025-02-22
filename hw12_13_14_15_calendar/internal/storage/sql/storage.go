package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"                                 //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Storage struct {
	driver, dsn string
	db          *sql.DB
}

func New(driver, dsn string) *Storage {
	return &Storage{
		driver: driver,
		dsn:    dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sql.Open(s.driver, s.dsn)
	if err != nil {
		return fmt.Errorf("%w: error while connecting to dsn %v using driver %v", err, s.dsn, s.driver)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	if event.StopTime.Before(event.StartTime) {
		return "", storage.ErrInvalidStopTime
	}

	var reminderTime *time.Time
	if event.Reminder != nil {
		tempTime := event.StartTime.Add(-*event.Reminder)
		reminderTime = &tempTime
	}
	row := s.db.QueryRowContext(ctx, `insert into event (id, title, startTime, stopTime, description, userID, reminder, 
	reminderTime) values (gen_random_uuid(),$1, $2, $3, $4, $5, $6, $7) 
	on conflict (starttime, userid) do nothing
	returning id`,
		event.Title, event.StartTime, event.StopTime, event.Description, event.UserID, event.Reminder, reminderTime)

	var id string
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w: %v", storage.ErrDateBusy, event)
		}
		return "", fmt.Errorf("%w: %v %v", storage.ErrCreateEvent, event, err) //nolint:errorlint
	}
	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	if event.ID != id {
		return fmt.Errorf("%w: id=%v event.ID=%v", storage.ErrInvalidArgiments, id, event.ID)
	}
	var reminderTime *time.Time
	if event.Reminder != nil {
		tempTime := event.StartTime.Add(-*event.Reminder)
		reminderTime = &tempTime
	}
	result, err := s.db.ExecContext(ctx, `update event 
	SET title = $1, startTime = $2, stopTime = $3, description = $4, reminder = $5, reminderTime = $6 
	WHERE id = $7 and userID = $8`,
		event.Title, event.StartTime, event.StopTime, event.Description, event.Reminder, reminderTime,
		event.ID, event.UserID)
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrUpdateEvent, event, err) //nolint:errorlint
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrUpdateEvent, event, err) //nolint:errorlint
	}
	if rowCount != 1 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `delete from event where id=$1;`, id)
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrDeleteEvent, id, err) //nolint:errorlint
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrDeleteEvent, id, err) //nolint:errorlint
	}
	if rows != 1 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*storage.Event, error) {
	row := s.db.QueryRowContext(ctx, `select id, title, starttime, stoptime, description, userid, reminder from "event" 
	where id = $1`, id)
	event := storage.Event{}
	err := row.Scan(&event.ID, &event.Title, &event.StartTime, &event.StopTime, &event.Description, &event.UserID,
		&event.Reminder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrEventNotFound
		}
		return nil, fmt.Errorf("%w: %v %v", storage.ErrReadEvent, id, err) //nolint:errorlint
	}
	return &event, nil
}

func (s *Storage) ListEventsDay(ctx context.Context, startTime time.Time) ([]*storage.Event, error) {
	return s.listEventsInt(ctx, startTime, startTime.Add(time.Hour*24))
}

func (s *Storage) ListEventsWeek(ctx context.Context, startTime time.Time) ([]*storage.Event, error) {
	return s.listEventsInt(ctx, startTime, startTime.Add(time.Hour*24*7))
}

func (s *Storage) ListEventsMonth(ctx context.Context, startTime time.Time) ([]*storage.Event, error) {
	return s.listEventsInt(ctx, startTime, startTime.AddDate(0, 1, 0))
}

func (s *Storage) listEventsInt(ctx context.Context, startTime time.Time, stopTime time.Time) (
	[]*storage.Event, error,
) {
	rows, err := s.db.QueryContext(ctx, `select id, title, starttime, stoptime, description, userid, reminder 
	from event where starttime >= $1 and stoptime <= $2`, startTime, stopTime)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %v", storage.ErrReadEvent, err) //nolint:errorlint
	}
	defer rows.Close()
	return makeEventsFromRows(rows)
}

func (s *Storage) ListEventsReminder(ctx context.Context) ([]*storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `select id, title, starttime, stoptime, description, userid, reminder 
	from event where ReminderTime < CURRENT_TIMESTAMP`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %v", storage.ErrReadEvent, err) //nolint:errorlint
	}
	defer rows.Close()
	return makeEventsFromRows(rows)
}

func makeEventsFromRows(rows *sql.Rows) ([]*storage.Event, error) {
	result := make([]*storage.Event, 0)
	for rows.Next() {
		event := &storage.Event{}
		var reminderStr string
		err := rows.Scan(&event.ID, &event.Title, &event.StartTime, &event.StopTime, &event.Description, &event.UserID,
			&reminderStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", storage.ErrReadEvent, err) //nolint:errorlint
		}
		if reminderStr != "" {
			reminderTime, err := time.Parse(time.TimeOnly, reminderStr)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", storage.ErrReadEvent, err) //nolint:errorlint
			}
			reminder := time.Duration(reminderTime.Hour()*60*60+reminderTime.Minute()*60+reminderTime.Second()) * 1000000000
			event.Reminder = &reminder
		}
		result = append(result, event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%w: %v", storage.ErrReadEvent, rows.Err()) //nolint:errorlint
	}
	return result, nil
}

func (s *Storage) ClearReminderTime(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `update event set reminderTime = null where id=$1;`, id)
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrUpdateEvent, id, err) //nolint:errorlint
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrUpdateEvent, id, err) //nolint:errorlint
	}
	if rows != 1 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) DeleteEventsBeforeDate(ctx context.Context, time time.Time) error {
	_, err := s.db.ExecContext(ctx, `delete from event where starttime < $1;`, time)
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrDeleteEvent, time, err) //nolint:errorlint
	}
	return nil
}

func (s *Storage) SaveNotification(ctx context.Context, notification storage.Notification) error {
	_, err := s.db.ExecContext(
		ctx, `insert into notification (id, title, startTime, userID) values ($1, $2, $3, $4) on conflict (id) do nothing`,
		notification.ID, notification.Title, notification.StartTime, notification.UserID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v %v", storage.ErrCreateNotification, notification, err) //nolint:errorlint
	}
	return nil
}
