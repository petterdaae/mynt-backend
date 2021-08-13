
-- +migrate Up
alter table users add column sbanken_client_id text;
alter table users add column sbanken_client_secret text;

-- +migrate Down
alter table users drop column sbanken_client_id;
alter table users drop column sbanken_client_secret;
