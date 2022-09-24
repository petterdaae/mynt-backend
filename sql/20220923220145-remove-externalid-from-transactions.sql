
-- +migrate Up
alter table transactions drop column external_id;

-- +migrate Down
alter table transactions add column external_id text;
