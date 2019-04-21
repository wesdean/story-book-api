package integration_test

import (
	"flag"
	"fmt"
	"github.com/wesdean/story-book-api/app_config"
	"github.com/wesdean/story-book-api/database"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var config *app_config.Config
var netClient = &http.Client{
	Timeout: time.Second * 10,
}

//todo Integration tests for forks
//todo Acceptance tests for forks

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		fmt.Println("Integration tests skipped in short mode")
		return
	}

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

	setupEnvironment()

	fmt.Println("CommandHook:BeforeTest")
	err = config.IntegrationTest.Hooks.RunHook("Hooks.CommandHooks.BeforeTest")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.Run()

	fmt.Println("CommandHook:AfterTest")
	err = config.IntegrationTest.Hooks.RunHook("Hooks.CommandHooks.AfterTest")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupEnvironment() {
	var err error

	err = os.Setenv("CONFIG_FILENAME", "../app_config/test.config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Setenv("AUTH_SECRET", "testing")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Setenv("AUTH_TIMEOUT", "3600")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func seedDb() {
	var err error

	db, err := database.NewDatabase(nil)
	err = db.Begin()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	sqlFile, err := ioutil.ReadFile("../database/sql/test_seed_integration.sql")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = db.Tx.Exec(string(sqlFile))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
}
