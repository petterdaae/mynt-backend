
-- +migrate Up
create table incoming_transactions (
    id serial primary key,
    user_id text not null,
    account_id text not null,
    accounting_date text not null,
    interest_date text not null,
    amount integer not null,
    text text not null,
    constraint fk_user_id foreign key (user_id) references users(id),
    constraint fk_account_id foreign key (account_id) references accounts(id)
);

-- +migrate Down
drop table incoming_transactions;

