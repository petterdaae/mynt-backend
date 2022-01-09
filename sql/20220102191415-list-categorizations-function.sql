-- +migrate Up

-- +migrate StatementBegin
create function list_categorizations(
    _user_id text,
    _from_date text,
    _to_date text
) returns table (
    id int,
    transaction_id text,
    amount int,
    category_id int
)
as
$$
begin
    return query
    select 
        c.id, 
        c.transaction_id,
        c.amount,
        c.category_id
    from
        categorizations as c,
        transactions as t
    where
            c.transaction_id = t.id
        and user_id = _user_id
        and coalesce(t.custom_date, split_part(t.accounting_date, 'T', 1)) >= _from_date
        and coalesce(t.custom_date, split_part(t.accounting_date, 'T', 1)) <= _to_date;
end;
$$
language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
drop function list_categorizations;

