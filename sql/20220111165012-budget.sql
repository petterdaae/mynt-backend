
-- +migrate Up
alter table categories add column budget integer;

-- +migrate Down
alter table categories drop column budget;

