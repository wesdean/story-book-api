truncate users;
insert into users (id, username, password, disabled, archived)
values (1, 'admin', 'adminpassword', false, false),
       (2, 'owner', 'ownerpassword', false, false),
       (3, 'author', 'authorpassword', false, false),
       (4, 'editor', 'editorpassword', false, false),
       (5, 'reader', 'readerpassword', false, false),
       (6, 'disabledreader', 'readerpassword', true, false),
       (7, 'archivedreader', 'readerpassword', false, true);