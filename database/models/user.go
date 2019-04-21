package models

import (
	"errors"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"github.com/wesdean/story-book-api/utils"
	"gopkg.in/guregu/null.v3"
	"os"
	"strconv"
	"time"
)

type UserStore struct {
	*Store
}

type User struct {
	Id        int
	Username  string
	CreatedOn null.Time
	Disabled  bool
	Archived  bool
}

type AuthenticatedUser struct {
	*User
	Timestamp int
}

func NewUserStore(db *database.Database, logger *logging.Logger) *UserStore {
	return &UserStore{Store: NewStore(db, logger)}
}

type UserQueryOptions struct {
	useId bool
	id    int

	useUsername bool
	username    string

	usePassword bool
	password    string

	useDisabled bool
	disabled    bool

	useArchived bool
	archived    bool

	useCreatedAtStart bool
	createdAtStart    time.Time

	useCreatedAtEnd bool
	createdAtEnd    time.Time
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

func (options *UserQueryOptions) Disabled(disabled bool) *UserQueryOptions {
	options.useDisabled = true
	options.disabled = disabled
	return options
}

func (options *UserQueryOptions) Archived(archived bool) *UserQueryOptions {
	options.useArchived = true
	options.archived = archived
	return options
}

func (options *UserQueryOptions) CreatedAt(createdAtStart *time.Time, createdAtEnd *time.Time) *UserQueryOptions {
	if createdAtStart != nil {
		options.useCreatedAtStart = true
		options.createdAtStart = *createdAtStart
	}

	if createdAtEnd != nil {
		options.useCreatedAtEnd = true
		options.createdAtEnd = *createdAtEnd
	}
	return options
}

func (store *UserStore) GetUsers(options *UserQueryOptions) ([]*User, error) {
	if options == nil {
		options = NewUserQueryOptions()
	}

	var err error
	sqlQuery := `select id, username, created_at, disabled, archived
		from users
		where
		($1 = false or ($1 = true and id = $2))
		and ($3 = false or ($3 = true and username = $4))
		and ($5 = false or ($5 = true and password = $6))
		and ($7 = false or ($7 = true and disabled = $8))
		and ($9 = false or ($9 = true and archived = $10))
		and ($11 = false or ($11 = true and created_at >= $12))
		and ($13 = false or ($13 = true and created_at <= $14))`
	args := []interface{}{
		options.useId,
		options.id,
		options.useUsername,
		options.username,
		options.usePassword,
		options.password,
		options.useDisabled,
		options.disabled,
		options.useArchived,
		options.archived,
		options.useCreatedAtStart,
		options.createdAtStart,
		options.useCreatedAtEnd,
		options.createdAtEnd,
	}
	rows, err := store.db.Tx.Query(sqlQuery, args...)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve users: %s", err.Error())
		return nil, err
	}
	defer (func() {
		err = rows.Close()
	})()

	var users []*User
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.CreatedOn,
			&user.Disabled,
			&user.Archived,
		)
		if err != nil {
			logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to scan user row: %s", err.Error())
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

func (store *UserStore) AuthenticateUser(token string) (*AuthenticatedUser, error) {
	if token == "" {
		return nil, errors.New("missing authentication token")
	}

	claims, err := utils.ParseJWTToken(token, []byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		return nil, err
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user in token")
	}

	timestamp, ok := claims["timestamp"].(float64)
	if !ok {
		return nil, errors.New("invalid timestamp in token")
	}

	authTimeout, err := strconv.Atoi(os.Getenv("AUTH_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	if (time.Now().Unix() - int64(timestamp)) > int64(authTimeout) {
		return nil, errors.New("ioken expired")
	}

	user, err := store.GetUser(NewUserQueryOptions().Id(int(userId)))
	if err != nil {
		return nil, err
	}

	if user == nil || user.Id <= 0 {
		return nil, errors.New("invalid user")
	}

	return &AuthenticatedUser{
		User:      user,
		Timestamp: int(timestamp),
	}, nil
}

func (store *UserStore) DisableUser(userId int) error {
	sqlQuery := "update users set disabled = true where id = $1"
	_, err := store.db.Tx.Exec(sqlQuery, userId)
	return err
}

func (store *UserStore) EnableUser(userId int) error {
	sqlQuery := "update users set disabled = false where id = $1"
	_, err := store.db.Tx.Exec(sqlQuery, userId)
	return err
}

func (store *UserStore) ArchiveUser(userId int) error {
	sqlQuery := "update users set archived = true where id = $1"
	_, err := store.db.Tx.Exec(sqlQuery, userId)
	return err
}

func (store *UserStore) UnarchiveUser(userId int) error {
	sqlQuery := "update users set archived = false where id = $1"
	_, err := store.db.Tx.Exec(sqlQuery, userId)
	return err
}
