
-- +migrate Up
create table users (
    id serial primary key,
    email text not null
);

-- +migrate Down
drop table users;
