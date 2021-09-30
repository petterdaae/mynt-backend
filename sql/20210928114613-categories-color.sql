
-- +migrate Up
alter table categories add column color text;

-- +migrate Down
alter table categories drop column color;
