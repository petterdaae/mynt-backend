
-- +migrate Up

-- +migrate StatementBegin
create or replace function update_transaction(
    _user_id text,
    _transaction_id text,
    _custom_date text
) returns void
as
$$
begin
    update transactions set custom_date = _custom_date
    where user_id = _user_id and id = _transaction_id;
end;
$$
language plpgsql;
-- +migrate StatementEnd

-- +migrate Down

-- +migrate StatementBegin
create or replace function update_transaction(
    _user_id text,
    _transaction_id text,
    _custom_date text
) returns void
as
$$
begin
    update transactions set custom_date = _custom_date
    where user_id = _user_id and _transaction_id = transaction_id;
end;
$$
language plpgsql;
-- +migrate StatementEnd
