-- +goose Up
-- +goose StatementBegin
alter table event add column if not exists ReminderTime timestamp with time zone null;
comment on column event.ReminderTime is 'Время отправки напоминания';
update event set ReminderTime = starttime - reminder where reminder is not null;
create index if not exists xie0_event_ReminderTime on event (ReminderTime);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists xie0_event_ReminderTime;
alter table event drop column if exists ReminderTime;
-- +goose StatementEnd
