-- +goose up
create table pft_type (
    id int generated always as identity primary key,
    description varchar(256)
);

insert into pft_type (description) values
    ('spend'),
    ('fixed'),
    ('income');

create table pft_category (
    id int generated always as identity primary key,
    description varchar(256) not null
);

create table pft_card (
    id int generated always as identity primary key,
    card_name varchar(256) not null
);

create table pft_transaction (
    id int generated always as identity primary key,
    date timestamptz not null default current_timestamp,
    description varchar(256),
    amount numeric,
    amount_expression varchar(256),
    type_id int references pft_type(id) on delete set null,
    category_id int references pft_category(id) on delete set null,
    card_id int references pft_card(id) on delete set null
);

-- +goose down
drop table if exists pft_transaction;
drop table if exists pft_card;
drop table if exists pft_category;
drop table if exists pft_type;
