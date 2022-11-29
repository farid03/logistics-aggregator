create table user_type
(
    id        serial primary key,
    type_name varchar not null
);
create table user_state
(
    id         serial primary key,
    state_name varchar not null
);
create table order_state
(
    id         serial primary key,
    state_name varchar not null
);

create table sys_user
(
    id            serial primary key,
    username      varchar not null,
    password      varchar not null,
    social_rating float,
    name          varchar,
    surname       varchar,
    user_state_id integer references user_state (id),
    user_type_id  integer references user_type (id)
);

create table user_rights
(
    sys_user_id integer references sys_user (id) on delete cascade primary key,
    r           bool not null,
    w           bool not null,
    x           bool not null
);

create table color
(
    code integer primary key,
    name varchar not null
);

create table car_body_type
(
    id   serial primary key,
    name varchar not null
);

create table loading_places
(
    id   serial primary key,
    name varchar not null
);

create table position
(
    id        serial primary key,
    latitude  float8 not null,
    longitude float8 not null
);

create table car
(
    id            serial primary key,
    owner_id      integer references sys_user (id) on delete cascade not null,
    license_plate varchar                                            not null,
    position_id   integer references position (id) on update cascade not null
);

create table specification
(
    id                serial primary key,
    car_id            integer references car (id) on delete cascade not null,
    length            integer check ( length > 0 )                  not null,
    height            integer check ( length > 0 )                  not null,
    color_code        integer references color (code)               not null,
    car_body_type_id  integer references car_body_type (id)         not null,
    loading_places_id integer references loading_places (id)        not null
);

create table feedback
(
    id          serial primary key,
    author_id   integer references sys_user (id) on delete no action not null,
    receiver_id integer references sys_user (id) on delete no action not null,
    order_id    integer,
    description varchar,
    points      integer                                              not null
);

create table executor_order
(
    id             serial primary key,
    owner_id       integer references sys_user (id) on delete cascade not null,
    title          varchar                                            not null,
    description    varchar,
    price          integer,
    order_state_id integer references order_state (id)                not null
);

create table required_specification
(
    id                serial primary key,
    length            integer check ( length > 0 )           not null,
    height            integer check ( length > 0 )           not null,
    car_body_type_id  integer references car_body_type (id)  not null,
    loading_places_id integer references loading_places (id) not null,
    color_code        integer references color (code)
);

create table customer_order
(
    id                        serial primary key,
    owner_id                  integer references sys_user (id) on delete cascade               not null,
    title                     varchar                                                          not null,
    description               varchar,
    price                     integer,
    order_state_id            integer references order_state (id)                              not null,
    from_id                   integer references position (id) on delete cascade               not null,
    to_id                     integer references position (id) on delete cascade               not null,
    required_specification_id integer references required_specification (id) on delete cascade not null
);
-- убедиться в корректности типов связей
