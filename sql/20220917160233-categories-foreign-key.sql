
-- +migrate Up
alter table categories 
add constraint fk_parent_id 
foreign key (parent_id) references categories(id)
on delete cascade;

-- +migrate Down
alter table categories
drop constraint fk_parent_id;
