package models_test

import (
	"errors"
	"fmt"
	"github.com/wesdean/story-book-api/app_config"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"io/ioutil"
	"os"
	"testing"
)

var db *database.Database
var config *app_config.Config
var logger *logging.Logger

func TestMain(m *testing.M) {
	var err error

	config, err = app_config.NewConfigFromFile("../../app_config/test.config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = config.Validate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logConfig, err := config.GetLogger("Config.Models.Logger")
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
	logger.Info("Setting up tests")

	db, err = database.NewDatabase(nil)
	if err == nil && db == nil {
		err = errors.New("db is nil")
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.Begin()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Info("Running tests")
	m.Run()

	logger.Info("Cleaning up tests")
	err = db.Commit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.GetDB().Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Info("Tests completed")
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
	sqlFile, err := ioutil.ReadFile(config.Models.DatabaseSeed)
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
