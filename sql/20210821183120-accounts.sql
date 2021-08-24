
-- +migrate Up
create table accounts (
    id text primary key,
    user_id text,
    external_id text,
    account_number text,
    name text,
    available integer,
    balance integer
);

-- +migrate Down
drop table accounts;
