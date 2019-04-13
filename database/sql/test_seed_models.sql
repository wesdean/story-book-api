truncate users restart identity cascade;
insert into users (id, username, password, disabled, archived)
values (1, 'admin', 'adminpassword', false, false),
       (2, 'owner', 'ownerpassword', false, false),
       (3, 'author', 'authorpassword', false, false),
       (4, 'editor', 'editorpassword', false, false),
       (5, 'reader', 'readerpassword', false, false),
       (6, 'disabledreader', 'readerpassword', true, false),
       (7, 'archivedreader', 'readerpassword', false, true);

truncate user_roles restart identity cascade;
insert into user_roles (name, label, description)
values (/*1*/'superuser', 'Superuser', 'Like Superman only awesomer'),
       (/*2*/'owner', 'Owner', 'Owner of a resource'),
       (/*3*/'author', 'Author', 'Author of a resource'),
       (/*4*/'editor', 'Editor', 'Editor of a resource'),
       (/*5*/'proofreader', 'Proofreader', 'Proofreader of a resource'),
       (/*6*/'reader', 'Reader', 'Reader of a resource');

insert into user_role_links (user_id, user_role_id, resource_type, resource_id)
values (1, 1, 'application', 0),
       (2, 2, 'fork', 0),
       (3, 3, 'fork', 0),
       (4, 4, 'fork', 0),
       (5, 5, 'fork', 0),
       (6, 6, 'fork', 0);