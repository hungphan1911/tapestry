-- +goose up
insert into pft_type (description) values ('non-salary');

create table pft_budget (
    id int generated always as identity primary key,
    month date not null,
    category_id int references pft_category(id) on delete cascade,
    amount numeric not null check (amount >= 0)
);

create unique index pft_budget_total_idx on pft_budget (month) where category_id is null;
create unique index pft_budget_category_idx on pft_budget (month, category_id) where category_id is not null;

create table pft_recurring (
    id int generated always as identity primary key,
    description varchar(256) not null,
    type_id int references pft_type(id) on delete set null,
    category_id int references pft_category(id) on delete set null,
    card_id int references pft_card(id) on delete set null,
    day_of_month int not null check (day_of_month between 1 and 31),
    start_month date not null,
    end_month date
);

create table pft_recurring_amount (
    recurring_id int not null references pft_recurring(id) on delete cascade,
    effective_month date not null,
    amount numeric not null,
    amount_expression varchar(256),
    primary key (recurring_id, effective_month)
);

alter table pft_transaction
    add column recurring_id int references pft_recurring(id) on delete set null;

-- +goose down
alter table pft_transaction drop column if exists recurring_id;
drop table if exists pft_recurring_amount;
drop table if exists pft_recurring;
drop table if exists pft_budget;
delete from pft_type where description = 'non-salary';
