-- +goose Up
-- +goose StatementBegin
create table notification(
  id uuid,
  title text,
  startTime timestamp with time zone,
  userID bigint
);
create unique index xpk_notification_id on notification (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table notification;
-- +goose StatementEnd
