package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"testing"
	"time"
)

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
}

func TestForkStore_GetFork(t *testing.T) {
	forkStore := models.NewForkStore(db, logger)
	seedDb()

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
}
