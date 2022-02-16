
-- +migrate Up
alter table budget_items add column kind text;
alter table budget_items add column custom_items text;

update budget_items set kind = 'monthly';

-- +migrate Down
alter table budget_items drop column kind;
alter table budget_items drop column custom_items;

