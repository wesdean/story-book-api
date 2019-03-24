create table users
(
  id         serial primary key,
  username   varchar(50),
  password   varchar(255),
  created_on timestamp not null default current_timestamp,
  last_login timestamp
);