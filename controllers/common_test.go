package controllers_test

import (
	"fmt"
	"github.com/wesdean/story-book-api/database"
	"io/ioutil"
	"os"
	"testing"
)

var db *database.Database

func TestMain(m *testing.M) {
	var err error
	db, err = database.NewDatabase(nil)

	err = db.Begin()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	seedDb()

	err = db.Commit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	err = db.GetDB().Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	m.Run()
}

func setupEnvironment(t *testing.T) {
	err := os.Setenv("AUTH_SECRET", "testing")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("AUTH_TIMEOUT", "3600")
	if err != nil {
		t.Fatal(err)
	}
}

func seedDb() {
	sqlFile, err := ioutil.ReadFile("../database/sql/test_seed_controllers.sql")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	_, err = db.Tx.Exec(string(sqlFile))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}
