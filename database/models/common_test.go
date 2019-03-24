package models_test

import (
	"github.com/wesdean/story-book-api/database"
	"testing"
)

var db *database.Database

func setupTest(t *testing.T) {
	var err error
	db, err = database.NewDatabase(nil)
	err = db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
}

func tearDown(t *testing.T) func() {
	return func() {
		var err error
		err = db.Rollback()
		if err != nil {
			t.Fatal(err)
			return
		}

		err = db.GetDB().Close()
		if err != nil {
			t.Fatal(err)
			return
		}
	}
}
