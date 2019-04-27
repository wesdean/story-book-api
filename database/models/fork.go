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
	ParentId    null.Int
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

	useOwner bool
	owner    int

	useAuthor bool
	author    int

	useEditor bool
	editor    int

	useProofreader bool
	proofreader    int

	useReader bool
	reader    int

	useUserCanRead bool
	userCanRead    int
}

func NewForkStore(db *database.Database, logger *logging.Logger) *ForkStore {
	return &ForkStore{Store: NewStore(db, logger)}
}

func NewForkQueryOptions() *ForkQueryOptions {
	return &ForkQueryOptions{}
}

//region QueryOptions
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

func (options *ForkQueryOptions) Owner(owner int) *ForkQueryOptions {
	options.useOwner = true
	options.owner = owner
	return options
}

func (options *ForkQueryOptions) Author(author int) *ForkQueryOptions {
	options.useAuthor = true
	options.author = author
	return options
}

func (options *ForkQueryOptions) Editor(editor int) *ForkQueryOptions {
	options.useEditor = true
	options.editor = editor
	return options
}

func (options *ForkQueryOptions) Proofreader(proofreader int) *ForkQueryOptions {
	options.useProofreader = true
	options.proofreader = proofreader
	return options
}

func (options *ForkQueryOptions) Reader(reader int) *ForkQueryOptions {
	options.useReader = true
	options.reader = reader
	return options
}

func (options *ForkQueryOptions) UserCanRead(userId int) *ForkQueryOptions {
	options.useUserCanRead = true
	options.userCanRead = userId
	return options
}

//endregion

var getForksCommonQuery = `
from forks
left join user_role_links as links 
	on (links.resource_type = 'fork' and (links.resource_id is null or links.resource_id = forks.id))
where ($1 = false or ($1 = true and id = $2))
	and ($3 = false or ($3 = true and (
		($4 <> 0 and forks.parent_id = $4) or ($4 = 0 and forks.parent_id is null)
	)))
	and ($5 = false or ($5 = true and forks.creator_id = $6))
	and ($7 = false or ($7 = true and lower(forks.title) like lower('%' || $8 || '%')))
	and ($9 = false or ($9 = true and lower(forks.description) like lower('%' || $10 || '%')))
	and ($11 = false or ($11 = true and ($12 = true and forks.published is not null or $12 = false and forks.published is null)))
	and ($13 = false or ($13 = true and forks.published >= $14))
	and ($15 = false or ($15 = true and forks.published <= $16))
	and ($17 = false or ($17 = true and links.user_id = $18 and user_role_id = $19))
	and ($20 = false or ($20 = true and links.user_id = $21 and user_role_id = $22))
	and ($23 = false or ($23 = true and links.user_id = $24 and user_role_id = $25))
	and ($26 = false or ($26 = true and links.user_id = $27 and user_role_id = $28))
	and ($29 = false or ($29 = true and links.user_id = $30 and user_role_id = $31))
	and ($32 = false or ($32 = true and links.user_id = $33 and 
		(user_role_id = $19 or user_role_id = $22 or user_role_id = $25 or user_role_id = $28 
		or user_role_id = $31 or user_role_id = $34)))
group by id`

func buildForksQueryArgs(options *ForkQueryOptions) []interface{} {
	return []interface{}{
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
		options.useOwner,
		options.owner,
		USER_ROLE_OWNER,
		options.useAuthor,
		options.author,
		USER_ROLE_AUTHOR,
		options.useEditor,
		options.editor,
		USER_ROLE_EDITOR,
		options.useProofreader,
		options.proofreader,
		USER_ROLE_PROOF,
		options.useReader,
		options.reader,
		USER_ROLE_READER,
		options.useUserCanRead,
		options.userCanRead,
		USER_ROLE_SUPER,
	}
}

func (store *ForkStore) CreateFork(fork *Fork) error {
	err := store.ValidateFork(fork)
	if err != nil {
		return err
	}

	sqlQuery := `insert into forks (parent_id, creator_id, title, description, published) 
		values ($1, $2, $3, $4, $5) returning id, created_at, updated_at`
	err = store.db.Tx.QueryRow(sqlQuery,
		fork.ParentId,
		fork.CreatorId,
		fork.Title,
		fork.Description,
		fork.Published,
	).Scan(&fork.Id, &fork.CreatedAt, &fork.UpdatedAt)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to create fork: %s", err.Error())
		return err
	}
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

	sqlQuery := `select id, parent_id, creator_id, title, description, published, created_at, updated_at ` +
		getForksCommonQuery

	args := buildForksQueryArgs(options)
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
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve fork: %s", err.Error())
		return nil, err
	}
	if len(forks) > 0 {
		return forks[0], nil
	}
	return nil, nil
}

func (store *ForkStore) GetForksWithBody(options *ForkQueryOptions) ([]*ForkWithBody, error) {
	if options == nil {
		options = NewForkQueryOptions()
	}

	sqlQuery := `select id, parent_id, creator_id, title, description, body, published, created_at, updated_at ` +
		getForksCommonQuery
	args := buildForksQueryArgs(options)
	rows, err := store.db.Tx.Query(sqlQuery, args...)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve forks with body: %s", err.Error())
		return nil, err
	}

	var forks []*ForkWithBody
	for rows.Next() {
		forkWithBody := &ForkWithBody{}
		fork := &Fork{}
		err := rows.Scan(
			&fork.Id,
			&fork.ParentId,
			&fork.CreatorId,
			&fork.Title,
			&fork.Description,
			&forkWithBody.Body,
			&fork.Published,
			&fork.CreatedAt,
			&fork.UpdatedAt)
		if err != nil {
			logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to scan fork with body: %s", err.Error())
			return nil, err
		}
		forkWithBody.Fork = fork
		forks = append(forks, forkWithBody)
	}
	return forks, nil
}

func (store *ForkStore) GetForkWithBody(options *ForkQueryOptions) (*ForkWithBody, error) {
	forks, err := store.GetForksWithBody(options)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve fork with body: %s", err.Error())
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
	if fork.ParentId.ValueOrZero() == 0 {
		sqlQuery := `select count(id) from forks where parent_id is null and creator_id = $1 and title = $2 and id <> $3`
		err = store.db.Tx.QueryRow(sqlQuery, fork.CreatorId, fork.Title, fork.Id).Scan(&count)
	} else {
		sqlQuery := `select count(id) from forks where parent_id = $1 and creator_id = $2 and title = $3 and id <> $4`
		err = store.db.Tx.QueryRow(sqlQuery, fork.ParentId, fork.CreatorId, fork.Title, fork.Id).Scan(&count)
	}
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("duplicate key")
	}

	return nil
}

func (store *ForkStore) UserCanCreate(userId, forkId int) (bool, error) {
	if userId > 0 && forkId == 0 {
		return true, nil
	}

	sqlQuery := `select count(*) from user_role_links
		where user_id = $1
		and resource_type = 'fork'
		and (resource_id = $2 or resource_id is null)
		and user_role_id in ($3, $4, $5)`
	var count int
	err := store.db.Tx.QueryRow(sqlQuery,
		userId,
		forkId,
		USER_ROLE_SUPER,
		USER_ROLE_OWNER,
		USER_ROLE_AUTHOR).Scan(&count)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve user role link count", err.Error())
		return false, err
	}
	return count > 0, nil
}

func (store *ForkStore) UserCanUpdate(userId, forkId int) (bool, error) {
	if userId <= 0 || forkId <= 0 {
		return false, nil
	}

	sqlQuery := `select count(*) from user_role_links
		where user_id = $1
		and resource_type = 'fork'
		and (resource_id = $2 or resource_id is null)
		and user_role_id in ($3, $4, $5, $6)`
	var count int
	err := store.db.Tx.QueryRow(sqlQuery,
		userId,
		forkId,
		USER_ROLE_SUPER,
		USER_ROLE_OWNER,
		USER_ROLE_AUTHOR,
		USER_ROLE_EDITOR).Scan(&count)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve user role link count", err.Error())
		return false, err
	}
	return count > 0, nil
}

func (store *ForkStore) UserCanDelete(userId, forkId int) (bool, error) {
	if userId <= 0 || forkId <= 0 {
		return false, nil
	}

	sqlQuery := `select count(*) from user_role_links
		where user_id = $1
		and resource_type = 'fork'
		and (resource_id = $2 or resource_id is null)
		and user_role_id in ($3, $4, $5)`
	var count int
	err := store.db.Tx.QueryRow(sqlQuery,
		userId,
		forkId,
		USER_ROLE_SUPER,
		USER_ROLE_OWNER,
		USER_ROLE_AUTHOR).Scan(&count)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve user role link count", err.Error())
		return false, err
	}
	return count > 0, nil
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
