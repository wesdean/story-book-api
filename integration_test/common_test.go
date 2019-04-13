package integration_test

import (
	"fmt"
	"github.com/wesdean/story-book-api/app_config"
	"github.com/wesdean/story-book-api/database"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

var config *app_config.Config

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

	setupEnvironment()

	if config.IntegrationTest.StartDocker != nil {
		cmd := exec.Command(config.IntegrationTest.StartDocker.Command, config.IntegrationTest.StartDocker.Arguments...)
		cmd.Dir = config.IntegrationTest.StartDocker.Directory
		fmt.Println("Building Docker container")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	m.Run()

	if config.IntegrationTest.StopDocker != nil {
		cmd := exec.Command(config.IntegrationTest.StopDocker.Command, config.IntegrationTest.StopDocker.Arguments...)
		cmd.Dir = config.IntegrationTest.StopDocker.Directory
		fmt.Println("Stopping Docker container")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func setupEnvironment() {
	var err error

	err = os.Setenv("CONFIG_FILENAME", "../app_config/test.config.json")

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
