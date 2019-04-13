package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"gopkg.in/guregu/null.v3"
	"testing"
)

func TestUserRoleStore_CreateRole(t *testing.T) {
	seedDb()

	userRoleStore := models.NewUserRoleStore(db, logger)
	newRole, err := userRoleStore.CreateRole("testRole", "Test Role", null.StringFrom("A test role"))
	if err != nil {
		t.Error(err)
		return
	}

	expected := 7
	if newRole.Id != expected {
		t.Errorf("expected %v, got %v", expected, newRole.Id)
		return
	}
}

func TestUserRoleStore_DeleteRole(t *testing.T) {
	seedDb()

	userRoleStore := models.NewUserRoleStore(db, logger)
	err := userRoleStore.DeleteRole(3)
	if err != nil {
		t.Error(err)
		return
	}

	roles, err := userRoleStore.GetRoles(nil)
	if err != nil {
		t.Error(err)
		return
	}
	expected := 5
	if len(roles) != expected {
		t.Errorf("expected %v, got %v", expected, len(roles))
		return
	}
}

func TestUserRoleStore_GetRoles(t *testing.T) {
	seedDb()

	userRoleStore := models.NewUserRoleStore(db, logger)
	roles, err := userRoleStore.GetRoles(nil)
	if err != nil {
		t.Error(err)
		return
	}

	expected := 6
	if len(roles) != expected {
		t.Errorf("expected %v, got %v", expected, len(roles))
		return
	}
}

func TestUserRoleStore_GetRole(t *testing.T) {
	seedDb()
	userRoleStore := models.NewUserRoleStore(db, logger)

	t.Run("Get first user role", func(t *testing.T) {
		role, err := userRoleStore.GetRole(nil)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 1
		if role.Id != expected {
			t.Errorf("expected %v, got %v", expected, role.Id)
			return
		}
	})

	t.Run("Get user role by ID", func(t *testing.T) {
		options := models.NewUserRoleQueryOptions().Id(3)
		role, err := userRoleStore.GetRole(options)
		if err != nil {
			t.Error(err)
			return
		}

		expected := "author"
		if role.Name != expected {
			t.Errorf("expected %v, got %v", expected, role.Name)
			return
		}
	})

	t.Run("Get user role by name", func(t *testing.T) {
		options := models.NewUserRoleQueryOptions().Name("owner")
		role, err := userRoleStore.GetRole(options)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 2
		if role.Id != expected {
			t.Errorf("expected %v, got %v", expected, role.Id)
			return
		}
	})

	t.Run("Get user role by description", func(t *testing.T) {
		options := models.NewUserRoleQueryOptions().Description("Superman Only")
		role, err := userRoleStore.GetRole(options)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 1
		if role.Id != expected {
			t.Errorf("expected %v, got %v", expected, role.Id)
			return
		}
	})
}
