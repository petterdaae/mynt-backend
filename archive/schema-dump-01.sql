create table users
(
    id                    text not null primary key, -- TODO: change name to google_sub?
    sbanken_client_id     text,
    sbanken_client_secret text,
    main_budget           integer -- TODO: foreign key, on delete set to null?
);

create table accounts
(
    id             text not null primary key,
    user_id        text, -- TODO: foreign key, on delete cascade
    external_id    text,
    account_number text,
    name           text,
    available      integer,
    balance        integer
);

create table transactions
(
    id              text not null primary key,
    user_id         text, -- TODO: remove field? (is linked to users through accounts)
    account_id      text, -- TODO: foreign key, on delete cascade
    external_id     text,
    accounting_date text,
    interest_date   text,
    amount          integer,
    text            text,
    custom_date     text
);

create table categories
(
    id        serial primary key,
    user_id   text, -- TODO: remove field? (is linked to users through accounts)
    name      text,
    parent_id integer, -- TODO: foreign key, on delete cascade
    deleted   boolean, -- TODO: remove field
    color     text,
    ignore    boolean
);

create table categorizations
(
    id             integer default nextval('transactions_to_categories_id_seq'::regclass) not null
        constraint transactions_to_categories_pkey
            primary key, -- TODO: waddup?
    transaction_id text, -- TODO: foreign key, on cascade delete
    amount         integer,
    category_id    integer, -- TODO: foreign key, on cascade delete
    user_id        text -- TODO: remove field
);

create table budgets
(
    id      serial primary key,
    user_id text, -- TODO: foreign key, on cascade delete
    name    text,
    color   text
);

alter table budgets
    owner to postgres;

create table budget_items
(
    id             serial
        primary key,
    user_id        text,
    budget_id      integer,
    category_id    integer,
    name           text,
    monthly_amount integer,
    kind           text,
    custom_items   text
);

alter table budget_items
    owner to postgres;

create table names
(
    id           serial
        primary key,
    user_id      text,
    name         text,
    regex        text,
    replace_with text
);

alter table names
    owner to postgres;

