truncate users;
insert into users (id, username, password)
values (1, 'admin', 'adminpassword'),
       (2, 'owner', 'ownerpassword'),
       (3, 'author', 'authorpassword'),
       (4, 'editor', 'editorpassword'),
       (5, 'reader', 'readerpassword');