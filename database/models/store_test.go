package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"testing"
)

func TestNewModel(t *testing.T) {
	setupTest(t)
	defer tearDown(t)

	model := models.NewStore(db)
	err := model.Ping()
	if err != nil {
		t.Error(err)
	}
}
