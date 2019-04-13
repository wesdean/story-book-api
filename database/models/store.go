package models

import (
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
)

type Store struct {
	db     *database.Database
	logger *logging.Logger
}

func NewStore(db *database.Database, logger *logging.Logger) *Store {
	return &Store{
		db:     db,
		logger: logger,
	}
}

func (store *Store) Ping() error {
	return store.db.Ping()
}
