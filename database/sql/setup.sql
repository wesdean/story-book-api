create database storybook;

create user storybook with password 'storybook';
grant select on all tables in schema public to storybook;