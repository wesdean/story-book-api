create table forks (
    id serial primary key,
    parent_id int null references forks(id) on delete cascade,
    creator_id int not null references users(id) on delete cascade,
    title varchar(150) not null,
    description text not null,
    body text null,
    published timestamp null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index on forks(parent_id);
create index on forks(creator_id);
create unique index on forks(parent_id, creator_id, title) where parent_id is not null;
create unique index on forks(creator_id, title) where parent_id is null;
create index on forks(creator_id, title);