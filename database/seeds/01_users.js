exports.seed = function (knex) {
    // Deletes ALL existing entries
    return knex('users').del()
        .then(function () {
            // Inserts seed entries
            return knex('users').insert([
                {id: 1, username: 'admin', password: 'adminpassword'},
                {id: 2, username: 'owner', password: 'ownerpassword'},
                {id: 3, username: 'author', password: 'authorpassword'},
                {id: 4, username: 'editor', password: 'editorpassword'},
                {id: 5, username: 'proofreader', password: 'proofreader'},
                {id: 6, username: 'reader', password: 'readerpassword'},
            ]);
        });
};
