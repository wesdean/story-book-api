package models

import (
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
	ParentId    null.Int
	CreatorId   int
	Title       string
	Description string
	Body        string
	Published   null.Time
	CreatedAt   null.Time
	UpdatedAt   null.Time
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

func (store *ForkStore) GetForks(options *ForkQueryOptions) ([]*Fork, error) {
	if options == nil {
		options = NewForkQueryOptions()
	}

	sqlQuery := `select id, parent_id, creator_id, title, description, body, published, created_at, updated_at
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
			&fork.Body,
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
