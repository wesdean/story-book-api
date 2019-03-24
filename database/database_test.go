package database_test

import (
	"github.com/wesdean/story-book-api/database"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db, err := database.NewDatabase(nil)
	if err != nil {
		t.Error(err)
		return
	}

	if db == nil {
		t.Error("expected Database got nil")
		return
	}
}

func TestDatabase_GetDB(t *testing.T) {
	db, err := database.NewDatabase(nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Rollback()

	if db.GetDB() == nil {
		t.Error("expected *sql.DB got nil")
		return
	}
}

func TestDatabase_Ping(t *testing.T) {
	db, err := database.NewDatabase(nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Rollback()

	err = db.Ping()
	if err != nil {
		t.Error(err)
		return
	}
}
