package models_test

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/utils"
	"os"
	"testing"
)

func TestUserStore_GetUsers(t *testing.T) {
	t.Run("Get all users", func(t *testing.T) {
		seedDb()

		userStore := models.NewUserStore(db, logger)
		users, err := userStore.GetUsers(nil)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 7
		if len(users) != expected {
			t.Errorf("expected %v, got %v", expected, len(users))
			return
		}
	})

	t.Run("Get user by ID", func(t *testing.T) {
		userStore := models.NewUserStore(db, logger)
		users, err := userStore.GetUsers(models.NewUserQueryOptions().Id(2))
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(users) != expectedCount {
			t.Errorf("expectedCount %v, got %v", expectedCount, len(users))
			return
		}

		expectedUsername := "owner"
		if users[0].Username != expectedUsername {
			t.Errorf("expected %v, got %v", expectedUsername, users[0].Username)
			return
		}
	})

	t.Run("Get user by Username", func(t *testing.T) {
		userStore := models.NewUserStore(db, logger)
		users, err := userStore.GetUsers(models.NewUserQueryOptions().Username("owner"))
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(users) != expectedCount {
			t.Errorf("expectedCount %v, got %v", expectedCount, len(users))
			return
		}

		expectedId := 2
		if users[0].Id != expectedId {
			t.Errorf("expected %v, got %v", expectedId, users[0].Id)
			return
		}
	})

	t.Run("Get user by Username and Password", func(t *testing.T) {
		userStore := models.NewUserStore(db, logger)
		options := models.NewUserQueryOptions().
			Username("owner").
			Password("ownerpassword")
		users, err := userStore.GetUsers(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(users) != expectedCount {
			t.Errorf("expectedCount %v, got %v", expectedCount, len(users))
			return
		}

		expectedId := 2
		if users[0].Id != expectedId {
			t.Errorf("expected %v, got %v", expectedId, users[0].Id)
			return
		}
	})

	t.Run("Return empty array when no users found", func(t *testing.T) {
		userStore := models.NewUserStore(db, logger)
		options := models.NewUserQueryOptions().Id(-1)
		users, err := userStore.GetUsers(options)
		if err != nil {
			t.Error(err)
			return
		}

		if len(users) != 0 {
			t.Errorf("expected 0, got %v", len(users))
			return
		}
	})
}

func TestUserStore_GetUser(t *testing.T) {
	t.Run("Return single user", func(t *testing.T) {
		userStore := models.NewUserStore(db, logger)
		user, err := userStore.GetUser(nil)
		if err != nil {
			t.Error(err)
			return
		}

		if user.Id != 1 {
			t.Errorf("expected 1, got %v", user.Id)
			return
		}
	})

	t.Run("Return nil when no user found", func(t *testing.T) {
		userStore := models.NewUserStore(db, logger)
		options := models.NewUserQueryOptions().Id(-1)
		user, err := userStore.GetUser(options)
		if err != nil {
			t.Error(err)
			return
		}

		if user != nil {
			t.Errorf("expected nil, got %v", user)
			return
		}
	})
}

func TestUserStore_AuthenticateUser(t *testing.T) {
	setupEnvironment(t)

	t.Run("Successful authentication", func(t *testing.T) {
		user := &models.User{
			Id:       2,
			Username: "owner",
		}

		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": user.Id},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		userStore := models.NewUserStore(db, logger)
		authUser, err := userStore.AuthenticateUser(token)
		if err != nil {
			t.Error(err)
			return
		}

		if authUser.User.Id != user.Id {
			t.Errorf("expected %v, got %v", user.Id, authUser.User.Id)
			return
		}

		if authUser.Timestamp <= 0 {
			t.Errorf("expected >0, got %v", authUser.Timestamp)
			return
		}
	})
}

func TestUserStore_DisableUser(t *testing.T) {
	seedDb()

	userId := 4

	userStore := models.NewUserStore(db, logger)
	err := userStore.DisableUser(userId)
	if err != nil {
		t.Error(err)
		return
	}
	user, err := userStore.GetUser(models.NewUserQueryOptions().Id(userId))
	if err != nil {
		t.Error(err)
		return
	}

	if !user.Disabled {
		t.Error("expected true, got false")
		return
	}
}

func TestUserStore_EnableUser(t *testing.T) {
	seedDb()

	userId := 6

	userStore := models.NewUserStore(db, logger)
	err := userStore.EnableUser(userId)
	if err != nil {
		t.Error(err)
		return
	}
	user, err := userStore.GetUser(models.NewUserQueryOptions().Id(userId))
	if err != nil {
		t.Error(err)
		return
	}

	if user.Disabled {
		t.Error("expected false, got true")
		return
	}
}

func TestUserStore_ArchiveUser(t *testing.T) {
	seedDb()

	userId := 4

	userStore := models.NewUserStore(db, logger)
	err := userStore.ArchiveUser(userId)
	if err != nil {
		t.Error(err)
		return
	}
	user, err := userStore.GetUser(models.NewUserQueryOptions().Id(userId))
	if err != nil {
		t.Error(err)
		return
	}

	if !user.Archived {
		t.Error("expected true, got false")
		return
	}
}

func TestUserStore_UnarchiveUser(t *testing.T) {
	seedDb()

	userId := 7

	userStore := models.NewUserStore(db, logger)
	err := userStore.UnarchiveUser(userId)
	if err != nil {
		t.Error(err)
		return
	}
	user, err := userStore.GetUser(models.NewUserQueryOptions().Id(userId))
	if err != nil {
		t.Error(err)
		return
	}

	if user.Archived {
		t.Error("expected false, got true")
		return
	}
}
