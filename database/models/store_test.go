package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"testing"
)

func TestNewModel(t *testing.T) {
	model := models.NewStore(db, logger)
	err := model.Ping()
	if err != nil {
		t.Error(err)
	}
}
