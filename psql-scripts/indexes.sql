create index user_idx on sys_user using hash(id);

create index social_rating_idx on sys_user using btree(social_rating);

create index customer_order_idx on customer_order using btree(price);

create index executor_order_idx on executor_order using btree(price);

create index car_idx on car using hash(id);

create index position_idx on position using hash(id);