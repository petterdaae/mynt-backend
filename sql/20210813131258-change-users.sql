
-- +migrate Up
drop table users;
create table users (
    id text primary key
);

-- +migrate Down
drop table users;
create table users (
    id serial primary key,
    email text not null
);

