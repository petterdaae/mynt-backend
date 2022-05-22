
-- +migrate Up
create table names (
    id serial primary key,
    user_id text,
    name text,
    regex text,
    replace_with text
);

-- +migrate Down
drop table names;
