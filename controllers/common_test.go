package controllers_test

import (
	"fmt"
	"github.com/wesdean/story-book-api/app_config"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"io/ioutil"
	"os"
	"testing"
)

var config *app_config.Config
var logger *logging.Logger

func TestMain(m *testing.M) {
	var err error

	config, err = app_config.NewConfigFromFile("../app_config/test.config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = config.Validate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logConfig, err := config.GetLogger("Config.Controllers.Logger")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger, err = logging.NewLogger(logConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logging.CloseLogger(logger)

	logger.Info("Running tests")
	m.Run()

	logger.Info("Tests completed")
}

func setupEnvironment(t *testing.T) {
	var err error

	err = os.Setenv("CONFIG_FILENAME", "../app_config/test.config.json")

	err = os.Setenv("AUTH_SECRET", "testing")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("AUTH_TIMEOUT", "3600")
	if err != nil {
		t.Fatal(err)
	}
}

func openDB() *database.Database {
	db, err := database.NewDatabase(nil)

	err = db.Begin()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}
	return db
}

func closeDB(db *database.Database) {
	if db == nil {
		return
	}

	err := db.Commit()
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

	db = nil
}

func seedDb() {
	db := openDB()
	defer closeDB(db)

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
