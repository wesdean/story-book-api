package models

import (
	"github.com/wesdean/story-book-api/database"
)

type Stores struct {
	db        *database.Database
	UserStore *UserStore
}

func NewStores(db *database.Database) *Stores {
	return &Stores{
		db:        db,
		UserStore: NewUserStore(db),
	}
}
