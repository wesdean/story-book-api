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

insert into forks (id, parent_id, creator_id, title, description, body, published, created_at, updated_at)
values (1, null, 1, 'Test Story 1', 'This is a test story', 'Testing 1, 2, blah, blah...', null, '2019-03-01 16:47:32',
        '2019-03-01 16:47:32'),
       (2, 1, 1, 'Test Fork 1', 'This is a test fork', 'Testing 1, 2, blah, blah...', '2019-03-05 16:47:32', '2019-03-02 16:47:32',
        '2019-03-01 16:47:32'),
       (3, 1, 1, 'Test Fork 2', 'Stick a fork in Me!', 'Testing 1, 2, blah, blah...', null, '2019-03-03 16:47:32',
        '2019-03-01 16:47:32'),
       (4, null, 2, 'Test Story 2', 'This is a test story', 'Testing 1, 2, blah, blah...', '2019-04-17 16:47:32',
        '2019-03-02 16:47:32',
        '2019-03-01 16:47:32'),
       (5, 4, 3, 'Test Fork 1', 'This is a test fork', 'Testing 1, 2, blah, blah...', null, '2019-03-03 16:47:32',
        '2019-03-01 16:47:32');