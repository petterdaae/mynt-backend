
-- +migrate Up
create table budgets (
    id serial primary key,
    user_id text,
    name text,
    color text
);

create table budget_items (
    id serial primary key,
    user_id text,
    budget_id integer,
    category_id integer,
    negative_amount integer,
    positive_amount integer
);

-- +migrate Down
drop table budgets;
drop table budget_items;

