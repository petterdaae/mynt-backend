
-- +migrate Up
alter table budget_items add column monthly_amount integer;
update budget_items set monthly_amount = positive_amount where positive_amount is not null and negative_amount is null;
update budget_items set monthly_amount = -negative_amount where negative_amount is not null and positive_amount is null;
update budget_items set monthly_amount = positive_amount - negative_amount where negative_amount is not null and positive_amount is not null;
alter table budget_items drop column positive_amount;
alter table budget_items drop column negative_amount;

-- +migrate Down
alter table budget_items add column positive_amount integer;
alter table budget_items add column negative_amount integer;
update budget_items set positive_amount = monthly_amount where monthly_amount >= 0;
update budget_items set negative_amount = -monthly_amount where monthly_amount < 0;
alter table budget_items drop column monthly_amount;
