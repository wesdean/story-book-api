exports.up = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('user_roles').then((exists) => {
            if (!exists) {
                return knex.schema.createTable('user_roles', function (table) {
                    table.increments('id');
                    table.text('name').notNullable().unique();
                    table.text('label').notNullable().unique();
                    table.text('description').nullable()
                });
            }
        })
    ]);
};

exports.down = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('user_roles').then((exists) => {
            if (exists) {
                return knex.schema.dropTable('user_roles')
            }
        })
    ]);
};
