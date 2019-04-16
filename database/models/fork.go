package models

import (
	"errors"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"gopkg.in/guregu/null.v3"
	"time"
)

type ForkStore struct {
	*Store
}

type Fork struct {
	Id          int
	ParentId    int
	CreatorId   int
	Title       string
	Description string
	Published   null.Time
	CreatedAt   null.Time
	UpdatedAt   null.Time
}

type ForkWithBody struct {
	*Fork
	Body string
}

type ForkQueryOptions struct {
	useId bool
	id    int

	useParentId bool
	parentId    int

	useCreatorId bool
	creatorId    int

	useTitle bool
	title    string

	useDescription bool
	description    string

	useIsPublished bool
	isPublished    bool

	usePublishedStart bool
	publishedStart    time.Time

	usePublishedEnd bool
	publishedEnd    time.Time
}

func NewForkStore(db *database.Database, logger *logging.Logger) *ForkStore {
	return &ForkStore{Store: NewStore(db, logger)}
}

func NewForkQueryOptions() *ForkQueryOptions {
	return &ForkQueryOptions{}
}

func (options *ForkQueryOptions) Id(id int) *ForkQueryOptions {
	options.useId = true
	options.id = id
	return options
}

func (options *ForkQueryOptions) ParentId(id int) *ForkQueryOptions {
	options.useParentId = true
	options.parentId = id
	return options
}

func (options *ForkQueryOptions) CreatorId(id int) *ForkQueryOptions {
	options.useCreatorId = true
	options.creatorId = id
	return options
}

func (options *ForkQueryOptions) Title(title string) *ForkQueryOptions {
	options.useTitle = true
	options.title = title
	return options
}

func (options *ForkQueryOptions) Description(description string) *ForkQueryOptions {
	options.useDescription = true
	options.description = description
	return options
}

func (options *ForkQueryOptions) IsPublished(isPublished bool) *ForkQueryOptions {
	options.useIsPublished = true
	options.isPublished = isPublished
	return options
}

func (options *ForkQueryOptions) Published(publishedStart time.Time, publishedEnd time.Time) *ForkQueryOptions {
	options.usePublishedStart = true
	options.publishedStart = publishedStart
	options.usePublishedEnd = true
	options.publishedEnd = publishedEnd
	return options
}

func (options *ForkQueryOptions) PublishedStart(publishedStart time.Time) *ForkQueryOptions {
	options.usePublishedStart = true
	options.publishedStart = publishedStart
	return options
}

func (options *ForkQueryOptions) PublishedEnd(publishedEnd time.Time) *ForkQueryOptions {
	options.usePublishedEnd = true
	options.publishedEnd = publishedEnd
	return options
}

func (store *ForkStore) CreateFork(fork *Fork) error {
	err := store.ValidateFork(fork)
	if err != nil {
		return err
	}

	sqlQuery := `insert into forks (parent_id, creator_id, title, description, published) 
		values ($1, $2, $3, $4, $5) returning id`
	var id int
	err = store.db.Tx.QueryRow(sqlQuery,
		fork.ParentId,
		fork.CreatorId,
		fork.Title,
		fork.Description,
		fork.Published,
	).Scan(&id)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to create fork: %s", err.Error())
		return err
	}

	if id <= 0 {
		logging.Log(store.logger, logging.LOGLEVEL_ERROR, "failed to get fork id")
		return errors.New("failed to get fork id")
	}

	fork.Id = id
	return nil
}

func (store *ForkStore) UpdateFork(fork *Fork) error {
	err := store.ValidateFork(fork)
	if err != nil {
		return err
	}

	sqlQuery := `update forks set
			parent_id = $1,
			creator_id = $2,
			title = $3,
			description = $4,
			published = $5,
			updated_at = $6
		where id = $7`
	_, err = store.db.Tx.Exec(sqlQuery,
		fork.ParentId,
		fork.CreatorId,
		fork.Title,
		fork.Description,
		fork.Published,
		time.Now(),
		fork.Id)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to update fork: %s", err.Error())
		return err
	}

	return nil
}

func (store *ForkStore) DeleteForks(options *ForkQueryOptions) (int, error) {
	sqlQuery := `with deleted as (delete from forks
		where ($1 = false or ($1 = true and id = $2))
			and ($3 = false or ($3 = true and parent_id = $4))
			and ($5 = false or ($5 = true and creator_id = $6))
			and ($7 = false or ($7 = true and lower(title) like lower('%' || $8 || '%')))
			and ($9 = false or ($9 = true and lower(description) like lower('%' || $10 || '%')))
			and ($11 = false or ($11 = true and ($12 = true and published is not null or $12 = false and published is null)))
			and ($13 = false or ($13 = true and published >= $14))
			and ($15 = false or ($15 = true and published <= $16))
		returning *)
		select count(*) from deleted`
	args := []interface{}{
		options.useId,
		options.id,
		options.useParentId,
		options.parentId,
		options.useCreatorId,
		options.creatorId,
		options.useTitle,
		options.title,
		options.useDescription,
		options.description,
		options.useIsPublished,
		options.isPublished,
		options.usePublishedStart,
		options.publishedStart,
		options.usePublishedEnd,
		options.publishedEnd,
	}
	var count int
	err := store.db.Tx.QueryRow(sqlQuery, args...).Scan(&count)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to delete forks: %s", err.Error())
		return 0, err
	}
	return count, nil
}

func (store *ForkStore) GetForks(options *ForkQueryOptions) ([]*Fork, error) {
	if options == nil {
		options = NewForkQueryOptions()
	}

	sqlQuery := `select id, parent_id, creator_id, title, description, published, created_at, updated_at
		from forks
		where ($1 = false or ($1 = true and id = $2))
			and ($3 = false or ($3 = true and parent_id = $4))
			and ($5 = false or ($5 = true and creator_id = $6))
			and ($7 = false or ($7 = true and lower(title) like lower('%' || $8 || '%')))
			and ($9 = false or ($9 = true and lower(description) like lower('%' || $10 || '%')))
			and ($11 = false or ($11 = true and ($12 = true and published is not null or $12 = false and published is null)))
			and ($13 = false or ($13 = true and published >= $14))
			and ($15 = false or ($15 = true and published <= $16))`
	args := []interface{}{
		options.useId,
		options.id,
		options.useParentId,
		options.parentId,
		options.useCreatorId,
		options.creatorId,
		options.useTitle,
		options.title,
		options.useDescription,
		options.description,
		options.useIsPublished,
		options.isPublished,
		options.usePublishedStart,
		options.publishedStart,
		options.usePublishedEnd,
		options.publishedEnd,
	}
	rows, err := store.db.Tx.Query(sqlQuery, args...)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve forks: %s", err.Error())
		return nil, err
	}

	var forks []*Fork
	for rows.Next() {
		fork := &Fork{}
		err := rows.Scan(
			&fork.Id,
			&fork.ParentId,
			&fork.CreatorId,
			&fork.Title,
			&fork.Description,
			&fork.Published,
			&fork.CreatedAt,
			&fork.UpdatedAt)
		if err != nil {
			logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to scan fork: %s", err.Error())
			return nil, err
		}
		forks = append(forks, fork)
	}
	return forks, nil
}

func (store *ForkStore) GetFork(options *ForkQueryOptions) (*Fork, error) {
	forks, err := store.GetForks(options)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve fork %s", err.Error())
		return nil, err
	}
	if len(forks) > 0 {
		return forks[0], nil
	}
	return nil, nil
}

func (store *ForkStore) ValidateFork(fork *Fork) error {
	err := fork.Validate()
	if err != nil {
		return err
	}

	var count int
	sqlQuery := `select count(id) from forks where parent_id = $1 and creator_id = $2 and title = $3 and id <> $4`
	err = store.db.Tx.QueryRow(sqlQuery, fork.ParentId, fork.CreatorId, fork.Title, fork.Id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("duplicate key forks_parent_id_creator_id_title_unique")
	}

	return nil
}

func (fork *Fork) Validate() error {
	if fork.CreatorId <= 0 {
		return errors.New("invalid creator id")
	}

	if fork.Title == "" {
		return errors.New("invalid title")
	}

	if fork.Description == "" {
		return errors.New("invalid description")
	}

	return nil
}
