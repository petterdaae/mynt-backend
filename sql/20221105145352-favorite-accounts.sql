
-- +migrate Up
alter table accounts add column favorite boolean default true;

-- +migrate Down
alter table accounts drop column favorite;
