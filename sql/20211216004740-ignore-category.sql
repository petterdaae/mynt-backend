
-- +migrate Up
alter table categories add column ignore boolean;

-- +migrate Down
alter table categories drop column ignore;
