exports.seed = function (knex) {
    // Deletes ALL existing entries
    return knex('user_roles').del()
        .then(function () {
            // Inserts seed entries
            return knex('user_roles').insert([
                {id: 1, name: 'superuser', label: 'Superuser', description: 'Like superman only awesomer'},
                {id: 2, name: 'owner', label: 'Owner', description: 'User owns the target resource'},
                {id: 3, name: 'author', label: 'Author', description: 'User is an author of the target resource'},
                {id: 4, name: 'editor', label: 'Editor', description: 'User is an editor of the target resource'},
                {
                    id: 5,
                    name: 'contributor',
                    label: 'Contributor',
                    description: 'User is a contributor to the target resource'
                },
                {
                    id: 6,
                    name: 'proofreader',
                    label: 'Proofreader',
                    description: 'User is a proofreader for the target resource'
                },
                {id: 7, name: 'reader', label: 'Reader', description: 'User is a reader of the target resource'}
            ]);
        });
};
