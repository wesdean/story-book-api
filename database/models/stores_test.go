package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"reflect"
	"testing"
)

func TestNewStores(t *testing.T) {
	stores := models.NewStores(db, logger)
	typ := reflect.TypeOf(stores)
	if typ.Elem().Name() != "Stores" {
		t.Errorf("expected %v, got %v", "Stores", typ.Elem().Name())
		return
	}
}
