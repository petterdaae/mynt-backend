
-- +migrate Up
alter table transactions add column custom_date text;

-- +migrate Down
alter table transactions drop column custom_date;
