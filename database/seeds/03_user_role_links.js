exports.seed = function (knex) {
    // Deletes ALL existing entries
    return knex('user_role_links').del()
        .then(function () {
            // Inserts seed entries
            return knex('user_role_links').insert([
                {user_id: 1, user_role_id: 1, resource_type: 'application', resource_id: 0},
                {user_id: 2, user_role_id: 2, resource_type: 'fork', resource_id: 0},
                {user_id: 3, user_role_id: 3, resource_type: 'fork', resource_id: 0},
                {user_id: 4, user_role_id: 4, resource_type: 'fork', resource_id: 0},
                {user_id: 5, user_role_id: 5, resource_type: 'fork', resource_id: 0},
                {user_id: 6, user_role_id: 6, resource_type: 'fork', resource_id: 0}
            ]);
        });
};
