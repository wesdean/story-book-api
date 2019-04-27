truncate users restart identity cascade;
insert into users (username, password, disabled, archived)
values (/*1*/'admin', 'adminpassword', false, false),
       (/*2*/'owner', 'ownerpassword', false, false),
       (/*3*/'author', 'authorpassword', false, false),
       (/*4*/'editor', 'editorpassword', false, false),
       (/*5*/'proofreader', 'proofreaderpassword', false, false),
       (/*6*/'reader', 'readerpassword', false, false),
       (/*7*/'disabledreader', 'readerpassword', true, false),
       (/*8*/'archivedreader', 'readerpassword', false, true);

truncate user_roles restart identity cascade;
insert into user_roles (name, label, description)
values (/*1*/'superuser', 'Superuser', 'Like Superman only awesomer'),
       (/*2*/'owner', 'Owner', 'Owner of a resource'),
       (/*3*/'author', 'Author', 'Author of a resource'),
       (/*4*/'editor', 'Editor', 'Editor of a resource'),
       (/*5*/'proofreader', 'Proofreader', 'Proofreader of a resource'),
       (/*6*/'reader', 'Reader', 'Reader of a resource');

truncate user_role_links restart identity cascade;
insert into user_role_links (user_id, user_role_id, resource_type, resource_id)
values (1, 1, 'fork', null),
       (2, 2, 'fork', 4),
       (3, 3, 'fork', 2),
       (4, 4, 'fork', 4),
       (5, 5, 'fork', 3),
       (6, 6, 'fork', 1),
       (6, 6, 'fork', 2),
       (6, 6, 'fork', 3),
       (2, 2, 'fork', 5),
       (2, 2, 'fork', 6),
       (2, 6, 'fork', 7),
       (3, 3, 'fork', 1),
       (3, 3, 'fork', 8),
       (4, 4, 'fork', 3),
       (5, 5, 'fork', 7),
       (5, 5, 'fork', 8),
       (6, 6, 'fork', 8),
       (4, 4, 'fork', 9);

truncate forks restart identity cascade;
insert into forks (parent_id, creator_id, title, description, body, published, created_at, updated_at)
values (/*1*/ null, 1, 'Test Story 1', 'This is a story about a girl', 'Testing 1, 2, blah, blah...', null,
              '2019-03-01 16:47:32',
              '2019-03-01 16:47:32'),
       (/*2*/ 1, 1, 'Test Fork 1', 'This is a test fork', 'Testing 1, 2, blah, blah...', '2019-03-05 16:47:32',
              '2019-03-02 16:47:32',
              '2019-03-01 16:47:32'),
       (/*3*/ 1, 1, 'Test Fork 2', 'Stick a fork in Me!', 'Testing 1, 2, blah, blah...', null, '2019-03-03 16:47:32',
              '2019-03-01 16:47:32'),
       (/*4*/ null, 2, 'Test Story 2', 'This is a story about a boy', 'Testing 1, 2, blah, blah...', '2019-04-17 16:47:32',
              '2019-03-02 16:47:32', '2019-03-01 16:47:32'),
       (/*5*/ 4, 3, 'Test Fork 1', 'This is a test fork', 'Testing 1, 2, blah, blah...', null, '2019-03-03 16:47:32',
              '2019-03-01 16:47:32'),
       (/*6*/ null, 2, 'Test story 3', 'This is a test story', 'Testing 1, 2, blah, blah...', null, '2019-01-23 03:43:32',
              '2019-01-23 03:43:32'),
       (/*7*/ null, 3, 'Test story 4', 'This is a test story about a chinchilla', 'Testing 1, 2, blah, blah...', null,
              '2019-01-23 03:43:32',
              '2019-01-23 03:43:32'),
       (/*8*/ null, 2, 'Test Story 5', 'This is a test story', 'Testing 1, 2, blah, blah...', '2019-01-23 03:43:32',
              '2019-01-21 03:43:32', '2019-01-22 03:43:32'),
       (/*9*/ null, 2, 'Test Story 6', 'This is a test story', 'Testing 1, 2, blah, blah...', null,
              '2019-01-21 03:43:32', '2019-01-22 03:43:32');