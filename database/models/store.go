package models

import (
	"github.com/wesdean/story-book-api/database"
)

type Store struct {
	db *database.Database
}

func NewStore(db *database.Database) *Store {
	return &Store{db: db}
}

func (store *Store) Ping() error {
	return store.db.Ping()
}
