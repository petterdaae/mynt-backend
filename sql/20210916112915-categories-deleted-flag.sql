
-- +migrate Up
alter table categories add column deleted boolean;

-- +migrate Down
alter table categories drop column deleted;
