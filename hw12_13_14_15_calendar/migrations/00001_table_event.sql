-- +goose Up
-- +goose StatementBegin
create table event(
  id uuid,
  title text,
  startTime timestamp with time zone,
  stopTime timestamp with time zone,
  description text,
  userID bigint,
  reminder interval
);
create unique index xpk_event_id on event (id);   
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table event;
-- +goose StatementEnd
