create table user_role_links (
    user_id int not null references users(id) on delete cascade,
    user_role_id int not null references user_roles(id) on delete cascade,
    resource_type varchar(20) not null,
    resource_id int not null,
    primary key (user_id, resource_type, resource_id)
);