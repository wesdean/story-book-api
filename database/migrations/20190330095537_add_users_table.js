exports.up = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('users').then((exists) => {
            if (!exists) {
                return knex.schema.createTable('users', function (table) {
                    table.increments('id');
                    table.text('username').notNullable().unique();
                    table.text('password').notNullable();
                    table.boolean('disabled').notNullable().defaultTo(false);
                    table.boolean('archived').notNullable().defaultTo(false);
                    table.timestamp('created_at').notNullable().defaultTo(knex.fn.now());
                });
            }
        })
    ]);
};

exports.down = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('users').then((exists) => {
            if (exists) {
                return knex.schema.dropTable('users')
            }
        })
    ]);
};
