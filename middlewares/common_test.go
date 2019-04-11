package middlewares_test

import (
	"os"
	"testing"
)

func setupEnvironment(t *testing.T) {
	var err error

	err = os.Setenv("CONFIG_FILENAME", "../app_config/test.config.json")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("AUTH_TIMEOUT", "3")
	if err != nil {
		t.Error(t)
		return
	}
}
