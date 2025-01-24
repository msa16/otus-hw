package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct { // TODO
	driver, dsn string
	db          *sql.DB
}

func New(driver, dsn string) (*Storage, error) {
	return &Storage{
		driver: driver,
		dsn:    dsn,
	}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sql.Open(s.driver, s.dsn)
	if err != nil {
		return fmt.Errorf("%w: error while connecting to dsn %v using driver %v", err, s.dsn, s.driver)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}
