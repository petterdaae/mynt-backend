
-- +migrate Up
alter table transactions_to_categories
rename to categorizations;

-- +migrate Down
alter table categorizations
rename to transactions_to_categories;
