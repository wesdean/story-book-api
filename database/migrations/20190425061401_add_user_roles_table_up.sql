create table user_roles (
    id serial primary key,
    name varchar(150) unique not null,
    label varchar(150) unique not null,
    description varchar(150) null
);