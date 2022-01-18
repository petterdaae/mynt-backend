
-- +migrate Up
alter table categories drop column budget;

-- +migrate Down
alter table categories add column budget integer;
