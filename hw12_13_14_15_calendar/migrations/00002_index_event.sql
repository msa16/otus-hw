-- +goose Up
-- +goose StatementBegin
create unique index xak0_event_startTime_UserID on event (startTime, UserID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index xak0_event_startTime_UserID;
-- +goose StatementEnd
