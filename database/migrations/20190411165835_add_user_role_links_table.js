exports.up = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('user_role_links').then((exists) => {
            if (!exists) {
                return knex.schema.createTable('user_role_links', function (table) {
                    table.integer('user_id').notNullable();
                    table.integer('user_role_id').notNullable();
                    table.text('resource_type').notNullable();
                    table.integer('resource_id').notNullable();
                    table.primary(['user_id', 'resource_type', 'resource_id']);
                    table.foreign('user_id').references('users.id').onDelete('cascade');
                    table.foreign('user_role_id').references('user_roles.id').onDelete('cascade');
                });
            }
        })
    ]);
};

exports.down = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('user_role_links').then((exists) => {
            if (exists) {
                return knex.schema.dropTable('user_role_links')
            }
        })
    ]);
};
