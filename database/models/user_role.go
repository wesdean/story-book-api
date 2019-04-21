package models

import (
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"gopkg.in/guregu/null.v3"
)

type UserRoleStore struct {
	*Store
}

type UserRole struct {
	Id          int
	Name        string
	Label       string
	Description null.String
}

type UserRoleQueryOptions struct {
	useId bool
	id    int

	useName bool
	name    string

	useLabel bool
	label    string

	useDescription bool
	description    string
}

const USER_ROLE_SUPER = 1
const USER_ROLE_OWNER = 2
const USER_ROLE_AUTHOR = 3
const USER_ROLE_EDITOR = 4
const USER_ROLE_PROOF = 5
const USER_ROLE_READER = 6

func NewUserRoleQueryOptions() *UserRoleQueryOptions {
	return &UserRoleQueryOptions{}
}

func (options *UserRoleQueryOptions) Id(id int) *UserRoleQueryOptions {
	options.useId = true
	options.id = id
	return options
}

func (options *UserRoleQueryOptions) Name(name string) *UserRoleQueryOptions {
	options.useName = true
	options.name = name
	return options
}

func (options *UserRoleQueryOptions) Label(label string) *UserRoleQueryOptions {
	options.useLabel = true
	options.label = label
	return options
}

func (options *UserRoleQueryOptions) Description(description string) *UserRoleQueryOptions {
	options.useDescription = true
	options.description = description
	return options
}

func NewUserRoleStore(db *database.Database, logger *logging.Logger) *UserRoleStore {
	return &UserRoleStore{Store: NewStore(db, logger)}
}

func (store *UserRoleStore) CreateRole(name, label string, description null.String) (*UserRole, error) {
	sqlQuery := `insert into user_roles (name, label, description) values ($1, $2, $3) returning id`
	var id int
	err := store.db.Tx.QueryRow(sqlQuery, name, label, description).Scan(&id)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to create user role: %s", err.Error())
		return nil, err
	}

	if id <= 0 {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to get ID after creating user role")
		return nil, err
	}
	return &UserRole{
		Id:          id,
		Name:        name,
		Label:       label,
		Description: description,
	}, nil
}

func (store *UserRoleStore) DeleteRole(roleId int) error {
	sqlQuery := `delete from user_roles where id = $1`
	_, err := store.db.Tx.Exec(sqlQuery, roleId)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to delete role: %s", err.Error())
		return err
	}
	return nil
}

func (store *UserRoleStore) GetRoles(options *UserRoleQueryOptions) ([]*UserRole, error) {
	if options == nil {
		options = NewUserRoleQueryOptions()
	}

	sqlQuery := `select id, name, label, description 
		from user_roles
		where ($1 = false or ($1 = true and id = $2))
		and ($3 = false or ($3 = true and name = $4))
		and ($5 = false or ($5 = true and label = $6))
		and ($7 = false or ($7 = true and lower(description) like lower('%' || $8 || '%')))`
	args := []interface{}{
		options.useId,
		options.id,
		options.useName,
		options.name,
		options.useLabel,
		options.label,
		options.useDescription,
		options.description,
	}

	rows, err := store.db.Tx.Query(sqlQuery, args...)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve user roles: %s", err.Error())
		return nil, err
	}

	var roles []*UserRole
	for rows.Next() {
		role := &UserRole{}
		err = rows.Scan(
			&role.Id,
			&role.Name,
			&role.Label,
			&role.Description,
		)
		if err != nil {
			logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to scan user role: %s", err.Error())
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, err
}

func (store *UserRoleStore) GetRole(options *UserRoleQueryOptions) (*UserRole, error) {
	roles, err := store.GetRoles(options)
	if err != nil {
		return nil, err
	}
	if len(roles) > 0 {
		return roles[0], nil
	}
	return nil, nil
}
