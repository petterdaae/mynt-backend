
-- +migrate Up

-- +migrate StatementBegin
create function list_transactions(
    _user_id text,
    _from_date text,
    _to_date text
) returns table (
    id int,
    account_id int,
    accounting_date text,
    interest_date text,
    custom_date text,
    amount int,
    text text
)
as
$$
begin
    return query
    select 
        id, 
        account_id, 
        split_part(accounting_date, 'T', 1) as accounting_date,
        split_part(interest_date, 'T', 1) as interest_date,
        custom_date,
        amount,
        text,
        category_id
    from
        transactions
    where
            user_id = _user_id
        and coalesce(custom_date, accounting_date) >= _from_date
        and coalesce(custom_date, accounting_date) <= _to_date
    order by
        coalesce(custom_date, accounting_date) desc, id;
end;
$$
language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
drop function list_transactions;
