package middlewares_test

import (
	"os"
	"testing"
)

func setupEnvironment(t *testing.T) {
	err := os.Setenv("AUTH_TIMEOUT", "3")
	if err != nil {
		t.Error(t)
		return
	}
}
