-- +goose Up
create table orders (
    order_id bigint primary key,
    client_id integer not null,
    store_until timestamptz not null,
    status varchar(50) not null,
    cost integer not null,
    weight integer not null,
    packages varchar[],
    pick_up_time timestamptz
);

-- +goose Down
drop table if exists orders;