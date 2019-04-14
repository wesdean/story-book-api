exports.up = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('forks').then((exists) => {
            if (!exists) {
                return knex.schema.createTable('forks', function (table) {
                    table.increments('id').primary();
                    table.integer('parent_id').nullable().index();
                    table.integer('creator_id').notNullable().index();
                    table.text('title').notNullable();
                    table.text('description').notNullable();
                    table.text('body').notNullable();
                    table.timestamp('published').nullable();
                    table.unique(['parent_id', 'title']);
                    table.timestamps(true, true);
                    table.foreign('parent_id').references('forks.id').onDelete('cascade');
                    table.foreign('creator_id').references('users.id').onDelete('cascade');
                });
            }
        })
    ]);
};

exports.down = function (knex, Promise) {
    return Promise.all([
        knex.schema.hasTable('forks').then((exists) => {
            if (exists) {
                return knex.schema.dropTable('forks')
            }
        })
    ]);
};
