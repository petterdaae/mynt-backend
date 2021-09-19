
-- +migrate Up
create table transactions_to_categories (
    id serial primary key,
    transaction_id text,
    amount integer,
    category_id integer,
    user_id text
);

-- +migrate Down
drop table transactions_to_categories;
