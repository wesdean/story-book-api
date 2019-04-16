package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/middlewares"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func forksController_Index(t *testing.T, token string) {
	t.Run("Default returns all top-level forks", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/forks")

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := strings.Trim(string(body), "\n")

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			return
		}

		var forksResp controllers.ForksControllerForksResponse
		err = json.Unmarshal(body, &forksResp)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 2
		if len(forksResp.Forks) != expected {
			t.Errorf("expected %v, got %v", expected, len(forksResp.Forks))
			return
		}
	})

	t.Run("Get forks by id", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/forks?id=3")

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := strings.Trim(string(body), "\n")

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			return
		}

		var forksResp controllers.ForksControllerForksResponse
		err = json.Unmarshal(body, &forksResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forksResp.Forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			return
		}

		expected := "Test Fork 2"
		if forksResp.Forks[0].Title != expected {
			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
			return
		}
	})

	t.Run("Get forks by parent", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/forks?parent_id=1")

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := strings.Trim(string(body), "\n")

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			return
		}

		var forksResp controllers.ForksControllerForksResponse
		err = json.Unmarshal(body, &forksResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 2
		if len(forksResp.Forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			return
		}

		expected := "Test Fork 1"
		if forksResp.Forks[0].Title != expected {
			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
			return
		}
	})

	t.Run("Get forks by creator", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/forks?creator_id=2")

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := strings.Trim(string(body), "\n")

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			return
		}

		var forksResp controllers.ForksControllerForksResponse
		err = json.Unmarshal(body, &forksResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forksResp.Forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			return
		}

		expected := "Test Story 2"
		if forksResp.Forks[0].Title != expected {
			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
			return
		}
	})

	t.Run("Get forks by title", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString(fmt.Sprintf("/forks?title=%s", url.QueryEscape("story")))

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := strings.Trim(string(body), "\n")

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			return
		}

		var forksResp controllers.ForksControllerForksResponse
		err = json.Unmarshal(body, &forksResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 2
		if len(forksResp.Forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			return
		}

		expected := 1
		if forksResp.Forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			return
		}
	})

	t.Run("Get forks by description", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString(fmt.Sprintf("/forks?description=%s", url.QueryEscape("girl")))

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := strings.Trim(string(body), "\n")

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			return
		}

		var forksResp controllers.ForksControllerForksResponse
		err = json.Unmarshal(body, &forksResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(forksResp.Forks) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			return
		}

		expected := 1
		if forksResp.Forks[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			return
		}
	})

	t.Run("Get forks by whether they are published", func(t *testing.T) {
		t.Run("Is published", func(t *testing.T) {
			handler := alice.New(
				middlewares.DatabaseMiddleware,
				middlewares.AuthenticationtMiddleware,
			).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			client := &http.Client{}

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString("/forks?is_published=true")

			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			bodyStr := strings.Trim(string(body), "\n")

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
				return
			}

			var forksResp controllers.ForksControllerForksResponse
			err = json.Unmarshal(body, &forksResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(forksResp.Forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
				return
			}

			expected := 4
			if forksResp.Forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
				return
			}
		})

		t.Run("Is not published", func(t *testing.T) {
			handler := alice.New(
				middlewares.DatabaseMiddleware,
				middlewares.AuthenticationtMiddleware,
			).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			client := &http.Client{}

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString("/forks?is_published=false")

			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			bodyStr := strings.Trim(string(body), "\n")

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
				return
			}

			var forksResp controllers.ForksControllerForksResponse
			err = json.Unmarshal(body, &forksResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(forksResp.Forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
				return
			}

			expected := 1
			if forksResp.Forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
				return
			}
		})
	})

	t.Run("Get forks by when they were published", func(t *testing.T) {
		t.Run("Start date only", func(t *testing.T) {
			handler := alice.New(
				middlewares.DatabaseMiddleware,
				middlewares.AuthenticationtMiddleware,
			).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			client := &http.Client{}

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks?published_start=%s", url.QueryEscape("2019-03-01 00:00:00-0600")))

			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			bodyStr := strings.Trim(string(body), "\n")

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
				return
			}

			var forksResp controllers.ForksControllerForksResponse
			err = json.Unmarshal(body, &forksResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(forksResp.Forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
				return
			}

			expected := 4
			if forksResp.Forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
				return
			}
		})

		t.Run("End date only", func(t *testing.T) {
			handler := alice.New(
				middlewares.DatabaseMiddleware,
				middlewares.AuthenticationtMiddleware,
			).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			client := &http.Client{}

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks?published_end=%s", url.QueryEscape("2019-04-27 00:00:00-0600")))

			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			bodyStr := strings.Trim(string(body), "\n")

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
				return
			}

			var forksResp controllers.ForksControllerForksResponse
			err = json.Unmarshal(body, &forksResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(forksResp.Forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
				return
			}

			expected := 4
			if forksResp.Forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
				return
			}
		})

		t.Run("Start and end date only", func(t *testing.T) {
			handler := alice.New(
				middlewares.DatabaseMiddleware,
				middlewares.AuthenticationtMiddleware,
			).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			client := &http.Client{}

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf(
				"/forks?published_start=%s&published_end=%s",
				url.QueryEscape("2019-03-01 00:00:00-0600"),
				url.QueryEscape("2019-04-27 00:00:00-0600"),
			))

			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			bodyStr := strings.Trim(string(body), "\n")

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
				return
			}

			var forksResp controllers.ForksControllerForksResponse
			err = json.Unmarshal(body, &forksResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(forksResp.Forks) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
				return
			}

			expected := 4
			if forksResp.Forks[0].Id != expected {
				t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
				return
			}
		})
	})
}

func TestForksController_Index(t *testing.T) {
	setupEnvironment(t)
	seedDb()

	t.Run("Unauthenticated user denied access", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/forks")

		req, err := http.NewRequest("GET", u.String(), nil)

		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := strings.Trim(string(body), "\n")

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
			return
		}
	})

	t.Run("As superuser", func(t *testing.T) {
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 1},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		forksController_Index(t, token)
	})

	t.Run("As owner", func(t *testing.T) {
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 2},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		forksController_Index(t, token)
	})

	t.Run("As author", func(t *testing.T) {
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 3},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		forksController_Index(t, token)
	})

	t.Run("As editor", func(t *testing.T) {
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 3},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		forksController_Index(t, token)
	})

	t.Run("As proofreader", func(t *testing.T) {
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 3},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		forksController_Index(t, token)
	})

	t.Run("As reader", func(t *testing.T) {
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 5},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		forksController_Index(t, token)
	})
}
