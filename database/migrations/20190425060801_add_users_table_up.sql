create table users (
    id serial primary key,
    username varchar(150) unique not null,
    password varchar(40) not null,
    disabled boolean not null default false,
    archived boolean not null default false,
    created_at timestamp not null default current_timestamp
);