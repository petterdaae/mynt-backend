
-- +migrate Up
create table categories (
    id serial primary key,
    user_id text,
    name text,
    parent_id integer
);

-- +migrate Down
drop table categories;

