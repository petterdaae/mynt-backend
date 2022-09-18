
-- +migrate Up
alter table categorizations 
add constraint fk_category_id 
foreign key (category_id) references categories(id)
on delete cascade;

-- +migrate Down
alter table categorizations
drop constraint fk_category_id;
