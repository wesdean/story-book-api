package models

import (
	"errors"
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"net/http"
)

type Stores struct {
	db                *database.Database
	logger            *logging.Logger
	UserStore         *UserStore
	UserRoleStore     *UserRoleStore
	UserRoleLinkStore *UserRoleLinkStore
	ForkStore         *ForkStore
}

func NewStores(db *database.Database, logger *logging.Logger) *Stores {
	return &Stores{
		db:                db,
		UserStore:         NewUserStore(db, logger),
		UserRoleStore:     NewUserRoleStore(db, logger),
		UserRoleLinkStore: NewUserRoleLinkStore(db, logger),
		ForkStore:         NewForkStore(db, logger),
	}
}

func GetStoresFromRequest(r *http.Request) (*Stores, error) {
	var stores *Stores
	storesContext, ok := context.GetOk(r, "Stores")
	if ok {
		stores = storesContext.(*Stores)
	} else {
		return nil, errors.New("missing database stores")
	}

	return stores, nil
}

func (stores *Stores) Commit() error {
	return stores.db.Commit()
}
