package models_test

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/utils"
	"gopkg.in/guregu/null.v3"
	"os"
	"testing"
)

func TestUserStore_GetUsers(t *testing.T) {
	setupTest(t)
	defer tearDown(t)

	t.Run("Get all users", func(t *testing.T) {
		userStore := models.NewUserStore(db)
		users, err := userStore.GetUsers(nil)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 5
		if len(users) != expected {
			t.Errorf("expected %v, got %v", expected, len(users))
			return
		}
	})

	t.Run("Get user by ID", func(t *testing.T) {
		userStore := models.NewUserStore(db)
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
		if users[0].Username.ValueOrZero() != expectedUsername {
			t.Errorf("expected %v, got %v", expectedUsername, users[0].Username.ValueOrZero())
			return
		}
	})

	t.Run("Get user by Username", func(t *testing.T) {
		userStore := models.NewUserStore(db)
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

		expectedId := int64(2)
		if users[0].Id.ValueOrZero() != expectedId {
			t.Errorf("expected %v, got %v", expectedId, users[0].Id.ValueOrZero())
			return
		}
	})

	t.Run("Get user by Username and Password", func(t *testing.T) {
		userStore := models.NewUserStore(db)
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

		expectedId := int64(2)
		if users[0].Id.ValueOrZero() != expectedId {
			t.Errorf("expected %v, got %v", expectedId, users[0].Id.ValueOrZero())
			return
		}
	})

	t.Run("Return empty array when no users found", func(t *testing.T) {
		userStore := models.NewUserStore(db)
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
		userStore := models.NewUserStore(db)
		user, err := userStore.GetUser(nil)
		if err != nil {
			t.Error(err)
			return
		}

		if user.Id.ValueOrZero() != 1 {
			t.Errorf("expected 1, got %v", user.Id)
			return
		}
	})

	t.Run("Return nil when no user found", func(t *testing.T) {
		userStore := models.NewUserStore(db)
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
	setupTest(t)
	defer tearDown(t)

	t.Run("Successful authentication", func(t *testing.T) {
		user := &models.User{
			Id:       null.IntFrom(2),
			Username: null.StringFrom("owner"),
		}

		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": user.Id},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		userStore := models.NewUserStore(db)
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
