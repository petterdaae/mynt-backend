
-- +migrate Up
alter table users add column main_budget integer;

-- +migrate Down
alter table users drop column main_budget;

