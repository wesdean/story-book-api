package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"strings"
	"testing"
	"time"
)

func TestForkStore_CreateFork(t *testing.T) {
	forkStore := models.NewForkStore(db, logger)

	t.Run("Successfully create fork", func(t *testing.T) {
		seedDb()

		fork := &models.Fork{
			CreatorId:   1,
			Title:       "--Test Fork--",
			Description: "--Description--",
		}
		err := forkStore.CreateFork(fork)
		if err != nil {
			t.Error(err)
			return
		}
		if fork.Id <= 0 {
			t.Errorf("expected > 0, got %v", fork.Id)
			return
		}
	})

	t.Run("Fail to create fork with duplicate creator, parent, title", func(t *testing.T) {
		seedDb()

		fork := &models.Fork{
			CreatorId:   1,
			Title:       "Test Story 1",
			Description: "--Description--",
		}
		err := forkStore.CreateFork(fork)
		if err == nil {
			t.Error("expected an duplicate key error, got none")
			return
		}

		expected := `duplicate key forks_parent_id_creator_id_title_unique`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Fail to create fork with missing creator", func(t *testing.T) {
		seedDb()

		fork := &models.Fork{
			Title:       "--Test Fork--",
			Description: "--Description--",
		}
		err := forkStore.CreateFork(fork)
		if err == nil {
			t.Error("expected an error, got none")
			return
		}

		expected := `invalid creator id`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Fail to create fork with missing title", func(t *testing.T) {
		seedDb()

		fork := &models.Fork{
			CreatorId:   1,
			Description: "--Description--",
		}
		err := forkStore.CreateFork(fork)
		if err == nil {
			t.Error("expected an error, got none")
			return
		}

		expected := `invalid title`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Fail to create fork with missing description", func(t *testing.T) {
		seedDb()

		fork := &models.Fork{
			CreatorId: 1,
			Title:     "--Test Fork--",
		}
		err := forkStore.CreateFork(fork)
		if err == nil {
			t.Error("expected an error, got none")
			return
		}

		expected := `invalid description`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})
}

func TestForkStore_UpdateFork(t *testing.T) {
	forkStore := models.NewForkStore(db, logger)

	t.Run("Successfully update fork", func(t *testing.T) {
		seedDb()

		forkId := 1
		title := "--Test Fork--"
		fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
		if err != nil {
			t.Error(err)
			return
		}
		fork.Title = title
		err = forkStore.UpdateFork(fork)
		if err != nil {
			t.Error(err)
			return
		}

		fork, err = forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
		if err != nil {
			t.Error(err)
			return
		}
		if fork.Title != title {
			t.Errorf("expected %v, got %v", title, fork.Title)
			return
		}
	})

	t.Run("Fail to update fork with duplicate creator, parent, title", func(t *testing.T) {
		seedDb()

		updateFork := models.Fork{
			Id:        1,
			ParentId:  0,
			CreatorId: 2,
			Title:     "Test Story 2",
		}
		fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(updateFork.Id))
		if err != nil {
			t.Error(err)
			return
		}

		fork.ParentId = updateFork.ParentId
		fork.CreatorId = updateFork.CreatorId
		fork.Title = updateFork.Title
		err = forkStore.UpdateFork(fork)
		if err == nil {
			t.Error("expected an duplicate key error, got none")
			return
		}

		expected := `duplicate key forks_parent_id_creator_id_title_unique`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Fail to update fork with missing creator", func(t *testing.T) {
		seedDb()

		forkId := 1
		fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
		if err != nil {
			t.Error(err)
			return
		}

		fork.CreatorId = 0
		err = forkStore.UpdateFork(fork)
		if err == nil {
			t.Error("expected an duplicate key error, got none")
			return
		}

		expected := `invalid creator id`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Fail to update fork with missing title", func(t *testing.T) {
		seedDb()

		forkId := 1
		fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
		if err != nil {
			t.Error(err)
			return
		}

		fork.Title = ""
		err = forkStore.UpdateFork(fork)
		if err == nil {
			t.Error("expected an duplicate key error, got none")
			return
		}

		expected := `invalid title`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Fail to update fork with missing description", func(t *testing.T) {
		seedDb()

		forkId := 1
		fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
		if err != nil {
			t.Error(err)
			return
		}

		fork.Description = ""
		err = forkStore.UpdateFork(fork)
		if err == nil {
			t.Error("expected an duplicate key error, got none")
			return
		}

		expected := `invalid description`
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("expected to contain %v, got %v", expected, err.Error())
			return
		}
	})
}

func TestForkStore_DeleteForks(t *testing.T) {
	forkStore := models.NewForkStore(db, logger)

	t.Run("Delete all forks", func(t *testing.T) {
		seedDb()

		options := models.NewForkQueryOptions()

		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 5
		if len(forks) != expected {
			t.Errorf("expected %v, got %v", expected, len(forks))
			return
		}

		deleteCount, err := forkStore.DeleteForks(options)
		if err != nil {
			t.Error(err)
			return
		}
		if deleteCount != expected {
			t.Errorf("expected %v, got %v", expected, deleteCount)
			return
		}
	})

	t.Run("Delete forks by id", func(t *testing.T) {
		seedDb()

		forkId := 3
		options := models.NewForkQueryOptions().Id(forkId)
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 1
		if len(forks) != expected {
			t.Errorf("expected %v, got %v", expected, len(forks))
			return
		}

		deleteCount, err := forkStore.DeleteForks(options)
		if err != nil {
			t.Error(err)
			return
		}
		if deleteCount != expected {
			t.Errorf("expected %v, got %v", expected, deleteCount)
			return
		}
	})

	t.Run("Delete forks by title", func(t *testing.T) {
		seedDb()

		options := models.NewForkQueryOptions().Title("Fork 1")
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 2
		if len(forks) != expected {
			t.Errorf("expected %v, got %v", expected, len(forks))
			return
		}

		deleteCount, err := forkStore.DeleteForks(options)
		if err != nil {
			t.Error(err)
			return
		}
		if deleteCount != expected {
			t.Errorf("expected %v, got %v", expected, deleteCount)
			return
		}
	})

	t.Run("Delete forks by description", func(t *testing.T) {
		seedDb()

		options := models.NewForkQueryOptions().Description("me")
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 1
		if len(forks) != expected {
			t.Errorf("expected %v, got %v", expected, len(forks))
			return
		}

		deleteCount, err := forkStore.DeleteForks(options)
		if err != nil {
			t.Error(err)
			return
		}
		if deleteCount != expected {
			t.Errorf("expected %v, got %v", expected, deleteCount)
			return
		}
	})

	t.Run("Delete forks by whether they are published", func(t *testing.T) {
		t.Run("Is published", func(t *testing.T) {
			seedDb()

			options := models.NewForkQueryOptions().IsPublished(true)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expected := 2
			if len(forks) != expected {
				t.Errorf("expected %v, got %v", expected, len(forks))
				return
			}

			deleteCount, err := forkStore.DeleteForks(options)
			if err != nil {
				t.Error(err)
				return
			}
			if deleteCount != expected {
				t.Errorf("expected %v, got %v", expected, deleteCount)
				return
			}
		})

		t.Run("Is not published", func(t *testing.T) {
			seedDb()

			options := models.NewForkQueryOptions().IsPublished(false)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expected := 3
			if len(forks) != expected {
				t.Errorf("expected %v, got %v", expected, len(forks))
				return
			}

			deleteCount, err := forkStore.DeleteForks(options)
			if err != nil {
				t.Error(err)
				return
			}
			if deleteCount != expected {
				t.Errorf("expected %v, got %v", expected, deleteCount)
				return
			}
		})
	})

	t.Run("Delete forks by when they were published", func(t *testing.T) {
		t.Run("Start date only", func(t *testing.T) {
			seedDb()

			publishedStart, err := time.Parse("2006-01-02", "2019-03-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().PublishedStart(publishedStart)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expected := 2
			if len(forks) != expected {
				t.Errorf("expected %v, got %v", expected, len(forks))
				return
			}

			deleteCount, err := forkStore.DeleteForks(options)
			if err != nil {
				t.Error(err)
				return
			}
			if deleteCount != expected {
				t.Errorf("expected %v, got %v", expected, deleteCount)
				return
			}
		})

		t.Run("End date only", func(t *testing.T) {
			seedDb()

			publishedEnd, err := time.Parse("2006-01-02", "2019-05-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().PublishedEnd(publishedEnd)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expected := 2
			if len(forks) != expected {
				t.Errorf("expected %v, got %v", expected, len(forks))
				return
			}

			deleteCount, err := forkStore.DeleteForks(options)
			if err != nil {
				t.Error(err)
				return
			}
			if deleteCount != expected {
				t.Errorf("expected %v, got %v", expected, deleteCount)
				return
			}
		})

		t.Run("Start and end date", func(t *testing.T) {
			seedDb()

			publishedStart, err := time.Parse("2006-01-02", "2019-03-01")
			if err != nil {
				t.Error(err)
				return
			}
			publishedEnd, err := time.Parse("2006-01-02", "2019-04-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().Published(publishedStart, publishedEnd)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expected := 1
			if len(forks) != expected {
				t.Errorf("expected %v, got %v", expected, len(forks))
				return
			}

			deleteCount, err := forkStore.DeleteForks(options)
			if err != nil {
				t.Error(err)
				return
			}
			if deleteCount != expected {
				t.Errorf("expected %v, got %v", expected, deleteCount)
				return
			}
		})

	})
}

func TestForkStore_GetForks(t *testing.T) {
	seedDb()

	forkStore := models.NewForkStore(db, logger)

	t.Run("Get all forks", func(t *testing.T) {
		forks, err := forkStore.GetForks(nil)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 5
		if len(forks) != expected {
			t.Errorf("expected %v, got %v", expected, len(forks))
			return
		}
	})

	t.Run("Get forks by id", func(t *testing.T) {
		forkId := 3
		options := models.NewForkQueryOptions().Id(forkId)
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := "Test Fork 2"
		if forks[0].Title != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Title)
			return
		}
	})

	t.Run("Get forks by title", func(t *testing.T) {
		options := models.NewForkQueryOptions().Title("Fork 1")
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 2
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 2
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by description", func(t *testing.T) {
		options := models.NewForkQueryOptions().Description("me")
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 3
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by whether they are published", func(t *testing.T) {
		t.Run("Is published", func(t *testing.T) {
			options := models.NewForkQueryOptions().IsPublished(true)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			expected := 2
			if forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forks[0].Id)
				return
			}
		})

		t.Run("Is not published", func(t *testing.T) {
			options := models.NewForkQueryOptions().IsPublished(false)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 3
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			expected := 1
			if forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forks[0].Id)
				return
			}
		})
	})

	t.Run("Get forks by when they were published", func(t *testing.T) {
		t.Run("Start date only", func(t *testing.T) {
			publishedStart, err := time.Parse("2006-01-02", "2019-03-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().PublishedStart(publishedStart)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			publishedStart, err = time.Parse("2006-01-02", "2019-04-01")
			if err != nil {
				t.Error(err)
				return
			}
			options = models.NewForkQueryOptions().PublishedStart(publishedStart)
			forks, err = forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount = 1
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}
		})

		t.Run("End date only", func(t *testing.T) {
			publishedEnd, err := time.Parse("2006-01-02", "2019-05-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().PublishedEnd(publishedEnd)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			publishedStart, err := time.Parse("2006-01-02", "2019-04-01")
			if err != nil {
				t.Error(err)
				return
			}
			options = models.NewForkQueryOptions().PublishedStart(publishedStart)
			forks, err = forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount = 1
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}
		})

		t.Run("Start and end date", func(t *testing.T) {
			publishedStart, err := time.Parse("2006-01-02", "2019-03-01")
			if err != nil {
				t.Error(err)
				return
			}
			publishedEnd, err := time.Parse("2006-01-02", "2019-04-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().Published(publishedStart, publishedEnd)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}
		})

	})

	t.Run("Get forks by owner", func(t *testing.T) {
		userId := 2

		options := models.NewForkQueryOptions().Owner(userId)
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 4
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by author", func(t *testing.T) {
		userId := 3

		options := models.NewForkQueryOptions().Author(userId)
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 2
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by editor", func(t *testing.T) {
		userId := 4

		options := models.NewForkQueryOptions().Editor(userId)
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 4
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by proofreader", func(t *testing.T) {
		userId := 5

		options := models.NewForkQueryOptions().Proofreader(userId)
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 3
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by reader", func(t *testing.T) {
		userId := 6

		options := models.NewForkQueryOptions().Reader(userId)
		forks, err := forkStore.GetForks(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 3
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 1
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by whether user can read", func(t *testing.T) {
		t.Run("Get for super user", func(t *testing.T) {
			userId := 1

			options := models.NewForkQueryOptions().UserCanRead(userId)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 5
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			expected := 1
			if forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forks[0].Id)
				return
			}
		})

		t.Run("Get for reader", func(t *testing.T) {
			userId := 6

			options := models.NewForkQueryOptions().UserCanRead(userId)
			forks, err := forkStore.GetForks(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 3
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			expected := 1
			if forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forks[0].Id)
				return
			}
		})
	})
}

func TestForkStore_GetFork(t *testing.T) {
	forkStore := models.NewForkStore(db, logger)
	seedDb()

	forkId := 3
	options := models.NewForkQueryOptions().Id(forkId)
	forks, err := forkStore.GetFork(options)
	if err != nil {
		t.Error(err)
		return
	}

	expected := "Test Fork 2"
	if forks.Title != expected {
		t.Errorf("expected %v, got %v", expected, forks.Title)
		return
	}
}

func TestForkStore_GetForksWithBody(t *testing.T) {
	seedDb()

	forkStore := models.NewForkStore(db, logger)

	t.Run("Get all forks", func(t *testing.T) {
		forks, err := forkStore.GetForksWithBody(nil)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 5
		if len(forks) != expected {
			t.Errorf("expected %v, got %v", expected, len(forks))
			return
		}

		if forks[0].Body == "" {
			t.Error("expected non-empty string, got empty string")
			return
		}
	})

	t.Run("Get forks by id", func(t *testing.T) {
		forkId := 3
		options := models.NewForkQueryOptions().Id(forkId)
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := "Test Fork 2"
		if forks[0].Title != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Title)
			return
		}
	})

	t.Run("Get forks by title", func(t *testing.T) {
		options := models.NewForkQueryOptions().Title("Fork 1")
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 2
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 2
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by description", func(t *testing.T) {
		options := models.NewForkQueryOptions().Description("me")
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 3
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by whether they are published", func(t *testing.T) {
		t.Run("Is published", func(t *testing.T) {
			options := models.NewForkQueryOptions().IsPublished(true)
			forks, err := forkStore.GetForksWithBody(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			expected := 2
			if forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forks[0].Id)
				return
			}
		})

		t.Run("Is not published", func(t *testing.T) {
			options := models.NewForkQueryOptions().IsPublished(false)
			forks, err := forkStore.GetForksWithBody(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 3
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			expected := 1
			if forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forks[0].Id)
				return
			}
		})
	})

	t.Run("Get forks by when they were published", func(t *testing.T) {
		t.Run("Start date only", func(t *testing.T) {
			publishedStart, err := time.Parse("2006-01-02", "2019-03-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().PublishedStart(publishedStart)
			forks, err := forkStore.GetForksWithBody(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			publishedStart, err = time.Parse("2006-01-02", "2019-04-01")
			if err != nil {
				t.Error(err)
				return
			}
			options = models.NewForkQueryOptions().PublishedStart(publishedStart)
			forks, err = forkStore.GetForksWithBody(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount = 1
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}
		})

		t.Run("End date only", func(t *testing.T) {
			publishedEnd, err := time.Parse("2006-01-02", "2019-05-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().PublishedEnd(publishedEnd)
			forks, err := forkStore.GetForksWithBody(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}

			publishedStart, err := time.Parse("2006-01-02", "2019-04-01")
			if err != nil {
				t.Error(err)
				return
			}
			options = models.NewForkQueryOptions().PublishedStart(publishedStart)
			forks, err = forkStore.GetForksWithBody(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount = 1
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}
		})

		t.Run("Start and end date", func(t *testing.T) {
			publishedStart, err := time.Parse("2006-01-02", "2019-03-01")
			if err != nil {
				t.Error(err)
				return
			}
			publishedEnd, err := time.Parse("2006-01-02", "2019-04-01")
			if err != nil {
				t.Error(err)
				return
			}
			options := models.NewForkQueryOptions().Published(publishedStart, publishedEnd)
			forks, err := forkStore.GetForksWithBody(options)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forks))
				return
			}
		})

	})

	t.Run("Get forks by owner", func(t *testing.T) {
		userId := 2

		options := models.NewForkQueryOptions().Owner(userId)
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 4
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by author", func(t *testing.T) {
		userId := 3

		options := models.NewForkQueryOptions().Author(userId)
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 2
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by editor", func(t *testing.T) {
		userId := 4

		options := models.NewForkQueryOptions().Editor(userId)
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 4
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by proofreader", func(t *testing.T) {
		userId := 5

		options := models.NewForkQueryOptions().Proofreader(userId)
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 3
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})

	t.Run("Get forks by reader", func(t *testing.T) {
		userId := 6

		options := models.NewForkQueryOptions().Reader(userId)
		forks, err := forkStore.GetForksWithBody(options)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 3
		if len(forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forks))
			return
		}

		expected := 1
		if forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forks[0].Id)
			return
		}
	})
}

func TestForkStore_GetForkWithBody(t *testing.T) {
	forkStore := models.NewForkStore(db, logger)
	seedDb()

	forkId := 3
	options := models.NewForkQueryOptions().Id(forkId)
	fork, err := forkStore.GetForkWithBody(options)
	if err != nil {
		t.Error(err)
		return
	}

	expected := "Test Fork 2"
	if fork.Title != expected {
		t.Errorf("expected %v, got %v", expected, fork.Title)
		return
	}

	if fork.Body == "" {
		t.Error("expected non-empty string, got empty string")
		return
	}
}

func TestFork_Validate(t *testing.T) {
	fork := models.Fork{
		Id:          0,
		CreatorId:   1,
		Title:       "Title",
		Description: "Description",
	}

	t.Run("Is valid fork", func(t *testing.T) {
		tFork := fork
		err := tFork.Validate()
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("Missing creator id", func(t *testing.T) {
		tFork := fork
		tFork.CreatorId = 0
		err := tFork.Validate()
		if err == nil {
			t.Error("expected not nil, got nil")
			return
		}

		expected := "invalid creator id"
		if err.Error() != expected {
			t.Errorf("expected %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Missing title", func(t *testing.T) {
		tFork := fork
		tFork.Title = ""
		err := tFork.Validate()
		if err == nil {
			t.Error("expected not nil, got nil")
			return
		}

		expected := "invalid title"
		if err.Error() != expected {
			t.Errorf("expected %v, got %v", expected, err.Error())
			return
		}
	})

	t.Run("Missing description", func(t *testing.T) {
		tFork := fork
		tFork.Description = ""
		err := tFork.Validate()
		if err == nil {
			t.Error("expected not nil, got nil")
			return
		}

		expected := "invalid description"
		if err.Error() != expected {
			t.Errorf("expected %v, got %v", expected, err.Error())
			return
		}
	})
}

func TestForkStore_UserCanCreate(t *testing.T) {
	seedDb()

	forkStore := models.NewForkStore(db, logger)

	superuser, err := forkStore.UserCanCreate(1, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if !superuser {
		t.Error("Superuser should be able to create")
		return
	}

	owner, err := forkStore.UserCanCreate(2, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if !owner {
		t.Error("Owner should be able to create")
		return
	}

	author, err := forkStore.UserCanCreate(3, 2)
	if err != nil {
		t.Error(err)
		return
	}
	if !author {
		t.Error("Author should be able to create")
		return
	}

	editor, err := forkStore.UserCanCreate(4, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if editor {
		t.Error("Editor should not be able to create")
		return
	}

	proofreader, err := forkStore.UserCanCreate(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if proofreader {
		t.Error("Proofreader should not be able to create")
		return
	}

	reader, err := forkStore.UserCanCreate(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if reader {
		t.Error("Reader should not be able to create")
		return
	}

	noPerms, err := forkStore.UserCanCreate(-1, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if noPerms {
		t.Error("User without permissions should not be able to create")
		return
	}
}

func TestForkStore_UserCanUpdate(t *testing.T) {
	seedDb()

	forkStore := models.NewForkStore(db, logger)

	superuser, err := forkStore.UserCanCreate(1, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if !superuser {
		t.Error("Superuser should be able to update")
		return
	}

	owner, err := forkStore.UserCanUpdate(2, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if !owner {
		t.Error("Owner should be able to update")
		return
	}

	author, err := forkStore.UserCanUpdate(3, 2)
	if err != nil {
		t.Error(err)
		return
	}
	if !author {
		t.Error("Author should be able to update")
		return
	}

	editor, err := forkStore.UserCanUpdate(4, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if !editor {
		t.Error("Editor should be able to update")
		return
	}

	proofreader, err := forkStore.UserCanUpdate(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if proofreader {
		t.Error("Proofreader should not be able to update")
		return
	}

	reader, err := forkStore.UserCanUpdate(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if reader {
		t.Error("Reader should not be able to update")
		return
	}

	noPerms, err := forkStore.UserCanUpdate(-1, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if noPerms {
		t.Error("User without permissions should not be able to create")
		return
	}
}

func TestForkStore_UserCanDelete(t *testing.T) {
	seedDb()

	forkStore := models.NewForkStore(db, logger)

	superuser, err := forkStore.UserCanDelete(1, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if !superuser {
		t.Error("Superuser should be able to delete")
		return
	}

	owner, err := forkStore.UserCanDelete(2, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if !owner {
		t.Error("Owner should be able to delete")
		return
	}

	author, err := forkStore.UserCanDelete(3, 2)
	if err != nil {
		t.Error(err)
		return
	}
	if !author {
		t.Error("Author should be able to delete")
		return
	}

	editor, err := forkStore.UserCanDelete(4, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if editor {
		t.Error("Editor should not be able to delete")
		return
	}

	proofreader, err := forkStore.UserCanDelete(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if proofreader {
		t.Error("Proofreader should not be able to delete")
		return
	}

	reader, err := forkStore.UserCanDelete(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if reader {
		t.Error("Reader should not be able to delete")
		return
	}

	noPerms, err := forkStore.UserCanDelete(-1, 4)
	if err != nil {
		t.Error(err)
		return
	}
	if noPerms {
		t.Error("User without permissions should not be able to delete")
		return
	}
}
