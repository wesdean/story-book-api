package models

import (
	"github.com/wesdean/story-book-api/database"
	"gopkg.in/guregu/null.v3"
	"time"
)

type UserStore struct {
	*Store
}

type User struct {
	Id        null.Int
	Username  null.String
	CreatedOn null.Time
	LastLogin null.Time
}

func NewUserStore(db *database.Database) *UserStore {
	return &UserStore{Store: NewStore(db)}
}

type UserQueryOptions struct {
	useId bool
	id    int

	useUsername bool
	username    string

	usePassword bool
	password    string

	useCreatedOnStart bool
	createdOnStart    time.Time

	useCreatedOnEnd bool
	createdOnEnd    time.Time

	useLastLoginStart bool
	lastLoginStart    time.Time

	useLastLoginEnd bool
	lastLoginEnd    time.Time
}

func NewUserQueryOptions() *UserQueryOptions {
	return &UserQueryOptions{}
}

func (options *UserQueryOptions) Id(id int) *UserQueryOptions {
	options.useId = true
	options.id = id
	return options
}

func (options *UserQueryOptions) Username(username string) *UserQueryOptions {
	options.useUsername = true
	options.username = username
	return options
}

func (options *UserQueryOptions) Password(password string) *UserQueryOptions {
	options.usePassword = true
	options.password = password
	return options
}

func (options *UserQueryOptions) CreatedOn(createdOnStart *time.Time, createdOnEnd *time.Time) *UserQueryOptions {
	if createdOnStart != nil {
		options.useCreatedOnStart = true
		options.createdOnStart = *createdOnStart
	}

	if createdOnEnd != nil {
		options.useCreatedOnEnd = true
		options.createdOnEnd = *createdOnEnd
	}
	return options
}

func (options *UserQueryOptions) LastLogin(lastLoginStart *time.Time, lastLoginEnd *time.Time) *UserQueryOptions {
	if lastLoginStart != nil {
		options.useLastLoginStart = true
		options.lastLoginStart = *lastLoginStart
	}

	if lastLoginEnd != nil {
		options.useLastLoginEnd = true
		options.lastLoginEnd = *lastLoginEnd
	}
	return options
}

func (store *UserStore) GetUsers(options *UserQueryOptions) ([]*User, error) {
	if options == nil {
		options = NewUserQueryOptions()
	}

	var err error
	sqlQuery := "select id, username, created_on, last_login " +
		"from users " +
		"where" +
		"($1 = false or ($1 = true and id = $2)) " +
		"and ($3 = false or ($3 = true and username = $4)) " +
		"and ($5 = false or ($5 = true and password = $6)) " +
		"and ($7 = false or ($7 = true and created_on >= $8)) " +
		"and ($9 = false or ($9 = true and created_on <= $10)) " +
		"and ($11 = false or ($11 = true and last_login >= $12)) " +
		"and ($13 = false or ($13 = true and last_login <= $14))"
	args := []interface{}{
		options.useId,
		options.id,
		options.useUsername,
		options.username,
		options.usePassword,
		options.password,
		options.useCreatedOnStart,
		options.createdOnStart,
		options.useCreatedOnEnd,
		options.createdOnEnd,
		options.useLastLoginStart,
		options.lastLoginStart,
		options.useLastLoginEnd,
		options.lastLoginEnd,
	}
	rows, err := store.db.Tx.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer (func() {
		err = rows.Close()
	})()

	users := []*User{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.CreatedOn,
			&user.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, err
}

func (store *UserStore) GetUser(options *UserQueryOptions) (*User, error) {
	users, err := store.GetUsers(options)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, nil
}
