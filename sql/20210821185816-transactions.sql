
-- +migrate Up
create table transactions (
    id text primary key,
    user_id text,
    account_id text,
    external_id text,
    accounting_date text,
    interest_date text,
    amount integer,
    text text
);

-- +migrate Down
drop table transactions;
