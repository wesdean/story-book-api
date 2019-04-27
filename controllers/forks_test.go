package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/database/models"
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

func TestForksController_Index(t *testing.T) {
	setupEnvironment(t)
	seedDb()

	t.Run("Without body", func(t *testing.T) {
		forksController_Index := func(t *testing.T, token string) {
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

				expected := 6
				if len(forksResp.Forks) != expected {
					t.Errorf("expected %v, got %v", expected, len(forksResp.Forks))
					return
				}
			})

			//t.Run("Get forks by id", func(t *testing.T) {
			//	handler := alice.New(
			//		middlewares.DatabaseMiddleware,
			//		middlewares.AuthenticationtMiddleware,
			//	).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//	testServer := httptest.NewServer(handler)
			//	defer testServer.Close()
			//
			//	client := &http.Client{}
			//
			//	var u bytes.Buffer
			//	u.WriteString(string(testServer.URL))
			//	u.WriteString("/forks?id=3")
			//
			//	req, err := http.NewRequest("GET", u.String(), nil)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//	req.Header.Set("Authorization", token)
			//
			//	resp, err := client.Do(req)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	body, _ := ioutil.ReadAll(resp.Body)
			//	bodyStr := strings.Trim(string(body), "\n")
			//
			//	if resp.StatusCode != http.StatusOK {
			//		t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//		return
			//	}
			//
			//	var forksResp controllers.ForksControllerForksResponse
			//	err = json.Unmarshal(body, &forksResp)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	expectedCount := 1
			//	if len(forksResp.Forks) != expectedCount {
			//		t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//		return
			//	}
			//
			//	expected := "Test Fork 2"
			//	if forksResp.Forks[0].Title != expected {
			//		t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
			//		return
			//	}
			//})
			//
			//t.Run("Get forks by parent", func(t *testing.T) {
			//	handler := alice.New(
			//		middlewares.DatabaseMiddleware,
			//		middlewares.AuthenticationtMiddleware,
			//	).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//	testServer := httptest.NewServer(handler)
			//	defer testServer.Close()
			//
			//	client := &http.Client{}
			//
			//	var u bytes.Buffer
			//	u.WriteString(string(testServer.URL))
			//	u.WriteString("/forks?parent_id=1")
			//
			//	req, err := http.NewRequest("GET", u.String(), nil)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//	req.Header.Set("Authorization", token)
			//
			//	resp, err := client.Do(req)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	body, _ := ioutil.ReadAll(resp.Body)
			//	bodyStr := strings.Trim(string(body), "\n")
			//
			//	if resp.StatusCode != http.StatusOK {
			//		t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//		return
			//	}
			//
			//	var forksResp controllers.ForksControllerForksResponse
			//	err = json.Unmarshal(body, &forksResp)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	expectedCount := 2
			//	if len(forksResp.Forks) != expectedCount {
			//		t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//		return
			//	}
			//
			//	expected := "Test Fork 1"
			//	if forksResp.Forks[0].Title != expected {
			//		t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
			//		return
			//	}
			//})
			//
			//t.Run("Get forks by creator", func(t *testing.T) {
			//	handler := alice.New(
			//		middlewares.DatabaseMiddleware,
			//		middlewares.AuthenticationtMiddleware,
			//	).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//	testServer := httptest.NewServer(handler)
			//	defer testServer.Close()
			//
			//	client := &http.Client{}
			//
			//	var u bytes.Buffer
			//	u.WriteString(string(testServer.URL))
			//	u.WriteString("/forks?creator_id=2")
			//
			//	req, err := http.NewRequest("GET", u.String(), nil)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//	req.Header.Set("Authorization", token)
			//
			//	resp, err := client.Do(req)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	body, _ := ioutil.ReadAll(resp.Body)
			//	bodyStr := strings.Trim(string(body), "\n")
			//
			//	if resp.StatusCode != http.StatusOK {
			//		t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//		return
			//	}
			//
			//	var forksResp controllers.ForksControllerForksResponse
			//	err = json.Unmarshal(body, &forksResp)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	expectedCount := 4
			//	if len(forksResp.Forks) != expectedCount {
			//		t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//		return
			//	}
			//
			//	expected := "Test Story 2"
			//	if forksResp.Forks[0].Title != expected {
			//		t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
			//		return
			//	}
			//})
			//
			//t.Run("Get forks by title", func(t *testing.T) {
			//	handler := alice.New(
			//		middlewares.DatabaseMiddleware,
			//		middlewares.AuthenticationtMiddleware,
			//	).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//	testServer := httptest.NewServer(handler)
			//	defer testServer.Close()
			//
			//	client := &http.Client{}
			//
			//	var u bytes.Buffer
			//	u.WriteString(string(testServer.URL))
			//	u.WriteString(fmt.Sprintf("/forks?title=%s", url.QueryEscape("story")))
			//
			//	req, err := http.NewRequest("GET", u.String(), nil)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//	req.Header.Set("Authorization", token)
			//
			//	resp, err := client.Do(req)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	body, _ := ioutil.ReadAll(resp.Body)
			//	bodyStr := strings.Trim(string(body), "\n")
			//
			//	if resp.StatusCode != http.StatusOK {
			//		t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//		return
			//	}
			//
			//	var forksResp controllers.ForksControllerForksResponse
			//	err = json.Unmarshal(body, &forksResp)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	expectedCount := 6
			//	if len(forksResp.Forks) != expectedCount {
			//		t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//		return
			//	}
			//
			//	expected := 1
			//	if forksResp.Forks[0].Id != expected {
			//		t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			//		return
			//	}
			//})
			//
			//t.Run("Get forks by description", func(t *testing.T) {
			//	handler := alice.New(
			//		middlewares.DatabaseMiddleware,
			//		middlewares.AuthenticationtMiddleware,
			//	).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//	testServer := httptest.NewServer(handler)
			//	defer testServer.Close()
			//
			//	client := &http.Client{}
			//
			//	var u bytes.Buffer
			//	u.WriteString(string(testServer.URL))
			//	u.WriteString(fmt.Sprintf("/forks?description=%s", url.QueryEscape("girl")))
			//
			//	req, err := http.NewRequest("GET", u.String(), nil)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//	req.Header.Set("Authorization", token)
			//
			//	resp, err := client.Do(req)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	body, _ := ioutil.ReadAll(resp.Body)
			//	bodyStr := strings.Trim(string(body), "\n")
			//
			//	if resp.StatusCode != http.StatusOK {
			//		t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//		return
			//	}
			//
			//	var forksResp controllers.ForksControllerForksResponse
			//	err = json.Unmarshal(body, &forksResp)
			//	if err != nil {
			//		t.Error(err)
			//		return
			//	}
			//
			//	expectedCount := 1
			//	if len(forksResp.Forks) != expectedCount {
			//		t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//		return
			//	}
			//
			//	expected := 1
			//	if forksResp.Forks[0].Id != expected {
			//		t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			//		return
			//	}
			//})
			//
			//t.Run("Get forks by whether they are published", func(t *testing.T) {
			//	t.Run("Is published", func(t *testing.T) {
			//		handler := alice.New(
			//			middlewares.DatabaseMiddleware,
			//			middlewares.AuthenticationtMiddleware,
			//		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//		testServer := httptest.NewServer(handler)
			//		defer testServer.Close()
			//
			//		client := &http.Client{}
			//
			//		var u bytes.Buffer
			//		u.WriteString(string(testServer.URL))
			//		u.WriteString("/forks?is_published=true")
			//
			//		req, err := http.NewRequest("GET", u.String(), nil)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//		req.Header.Set("Authorization", token)
			//
			//		resp, err := client.Do(req)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		body, _ := ioutil.ReadAll(resp.Body)
			//		bodyStr := strings.Trim(string(body), "\n")
			//
			//		if resp.StatusCode != http.StatusOK {
			//			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//			return
			//		}
			//
			//		var forksResp controllers.ForksControllerForksResponse
			//		err = json.Unmarshal(body, &forksResp)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		expectedCount := 2
			//		if len(forksResp.Forks) != expectedCount {
			//			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//			return
			//		}
			//
			//		expected := 4
			//		if forksResp.Forks[0].Id != expected {
			//			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			//			return
			//		}
			//	})
			//
			//	t.Run("Is not published", func(t *testing.T) {
			//		handler := alice.New(
			//			middlewares.DatabaseMiddleware,
			//			middlewares.AuthenticationtMiddleware,
			//		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//		testServer := httptest.NewServer(handler)
			//		defer testServer.Close()
			//
			//		client := &http.Client{}
			//
			//		var u bytes.Buffer
			//		u.WriteString(string(testServer.URL))
			//		u.WriteString("/forks?is_published=false")
			//
			//		req, err := http.NewRequest("GET", u.String(), nil)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//		req.Header.Set("Authorization", token)
			//
			//		resp, err := client.Do(req)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		body, _ := ioutil.ReadAll(resp.Body)
			//		bodyStr := strings.Trim(string(body), "\n")
			//
			//		if resp.StatusCode != http.StatusOK {
			//			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//			return
			//		}
			//
			//		var forksResp controllers.ForksControllerForksResponse
			//		err = json.Unmarshal(body, &forksResp)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		expectedCount := 4
			//		if len(forksResp.Forks) != expectedCount {
			//			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//			return
			//		}
			//
			//		expected := 1
			//		if forksResp.Forks[0].Id != expected {
			//			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			//			return
			//		}
			//	})
			//})
			//
			//t.Run("Get forks by when they were published", func(t *testing.T) {
			//	t.Run("Start date only", func(t *testing.T) {
			//		handler := alice.New(
			//			middlewares.DatabaseMiddleware,
			//			middlewares.AuthenticationtMiddleware,
			//		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//		testServer := httptest.NewServer(handler)
			//		defer testServer.Close()
			//
			//		client := &http.Client{}
			//
			//		var u bytes.Buffer
			//		u.WriteString(string(testServer.URL))
			//		u.WriteString(fmt.Sprintf("/forks?published_start=%s", url.QueryEscape("2019-03-01 00:00:00-0600")))
			//
			//		req, err := http.NewRequest("GET", u.String(), nil)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//		req.Header.Set("Authorization", token)
			//
			//		resp, err := client.Do(req)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		body, _ := ioutil.ReadAll(resp.Body)
			//		bodyStr := strings.Trim(string(body), "\n")
			//
			//		if resp.StatusCode != http.StatusOK {
			//			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//			return
			//		}
			//
			//		var forksResp controllers.ForksControllerForksResponse
			//		err = json.Unmarshal(body, &forksResp)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		expectedCount := 1
			//		if len(forksResp.Forks) != expectedCount {
			//			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//			return
			//		}
			//
			//		expected := 4
			//		if forksResp.Forks[0].Id != expected {
			//			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			//			return
			//		}
			//	})
			//
			//	t.Run("End date only", func(t *testing.T) {
			//		handler := alice.New(
			//			middlewares.DatabaseMiddleware,
			//			middlewares.AuthenticationtMiddleware,
			//		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//		testServer := httptest.NewServer(handler)
			//		defer testServer.Close()
			//
			//		client := &http.Client{}
			//
			//		var u bytes.Buffer
			//		u.WriteString(string(testServer.URL))
			//		u.WriteString(fmt.Sprintf(
			//			"/forks?published_end=%s",
			//			url.QueryEscape("2019-04-27 00:00:00-0600"),
			//		))
			//
			//		req, err := http.NewRequest("GET", u.String(), nil)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//		req.Header.Set("Authorization", token)
			//
			//		resp, err := client.Do(req)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		body, _ := ioutil.ReadAll(resp.Body)
			//		bodyStr := strings.Trim(string(body), "\n")
			//
			//		if resp.StatusCode != http.StatusOK {
			//			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//			return
			//		}
			//
			//		var forksResp controllers.ForksControllerForksResponse
			//		err = json.Unmarshal(body, &forksResp)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		expectedCount := 2
			//		if len(forksResp.Forks) != expectedCount {
			//			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//			return
			//		}
			//
			//		expected := 4
			//		if forksResp.Forks[0].Id != expected {
			//			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			//			return
			//		}
			//	})
			//
			//	t.Run("Start and end date only", func(t *testing.T) {
			//		handler := alice.New(
			//			middlewares.DatabaseMiddleware,
			//			middlewares.AuthenticationtMiddleware,
			//		).Then(middlewares.RunAPI(controllers.ForksController{}.Index))
			//
			//		testServer := httptest.NewServer(handler)
			//		defer testServer.Close()
			//
			//		client := &http.Client{}
			//
			//		var u bytes.Buffer
			//		u.WriteString(string(testServer.URL))
			//		u.WriteString(fmt.Sprintf(
			//			"/forks?published_start=%s&published_end=%s",
			//			url.QueryEscape("2019-03-01 00:00:00-0600"),
			//			url.QueryEscape("2019-04-27 00:00:00-0600"),
			//		))
			//
			//		req, err := http.NewRequest("GET", u.String(), nil)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//		req.Header.Set("Authorization", token)
			//
			//		resp, err := client.Do(req)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		body, _ := ioutil.ReadAll(resp.Body)
			//		bodyStr := strings.Trim(string(body), "\n")
			//
			//		if resp.StatusCode != http.StatusOK {
			//			t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
			//			return
			//		}
			//
			//		var forksResp controllers.ForksControllerForksResponse
			//		err = json.Unmarshal(body, &forksResp)
			//		if err != nil {
			//			t.Error(err)
			//			return
			//		}
			//
			//		expectedCount := 1
			//		if len(forksResp.Forks) != expectedCount {
			//			t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
			//			return
			//		}
			//
			//		expected := 4
			//		if forksResp.Forks[0].Id != expected {
			//			t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
			//			return
			//		}
			//	})
			//})
		}

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
				jwt.MapClaims{"user_id": 6},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

			forksController_Index(t, token)
		})
	})

	t.Run("With body", func(t *testing.T) {
		t.Run("Unauthenticated user denied access", func(t *testing.T) {
			handler := alice.New(
				middlewares.DatabaseMiddleware,
			).Then(middlewares.RunAPI(controllers.ForksController{}.Index))

			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			client := &http.Client{}

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString("/forks?with_body=true")

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

		t.Run("As superuser, can get all forks", func(t *testing.T) {
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": 1},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

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
				u.WriteString("/forks?with_body=true")

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

				expected := 6
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
				u.WriteString("/forks?with_body=true&id=3")

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
				u.WriteString("/forks?with_body=true&parent_id=1")

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
				u.WriteString("/forks?with_body=true&creator_id=2")

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

				expectedCount := 4
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&title=%s", url.QueryEscape("story")))

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

				expectedCount := 6
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&description=%s", url.QueryEscape("girl")))

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
					u.WriteString("/forks?with_body=true&is_published=true")

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
					u.WriteString("/forks?with_body=true&is_published=false")

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

					expectedCount := 4
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
					u.WriteString(fmt.Sprintf("/forks?with_body=true&published_start=%s", url.QueryEscape("2019-03-01 00:00:00-0600")))

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
					u.WriteString(fmt.Sprintf("/forks?with_body=true&published_end=%s", url.QueryEscape("2019-04-27 00:00:00-0600")))

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
						"/forks?with_body=true&published_start=%s&published_end=%s",
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
		})

		t.Run("As owner, can get forks with ownership or below", func(t *testing.T) {
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": 2},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

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
				u.WriteString("/forks?with_body=true")

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

				expectedCount := 3
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
				u.WriteString("/forks?with_body=true&id=5")

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

				expected := "Test Fork 1"
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
				u.WriteString("/forks?with_body=true&parent_id=4")

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
				u.WriteString("/forks?with_body=true&creator_id=2")

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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&title=%s", url.QueryEscape("story")))

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

				expectedCount := 3
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&description=%s", url.QueryEscape("boy")))

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
					u.WriteString("/forks?with_body=true&is_published=true")

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
					u.WriteString("/forks?with_body=true&is_published=false")

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

					expected := 6
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
					u.WriteString(fmt.Sprintf("/forks?with_body=true&published_start=%s", url.QueryEscape("2019-03-01 00:00:00-0600")))

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
					u.WriteString(fmt.Sprintf("/forks?with_body=true&published_end=%s", url.QueryEscape("2019-04-27 00:00:00-0600")))

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
						"/forks?with_body=true&published_start=%s&published_end=%s",
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
		})

		t.Run("As author, can get forks with authorship or below", func(t *testing.T) {
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": 3},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

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
				u.WriteString("/forks?with_body=true")

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
				u.WriteString("/forks?with_body=true&id=2")

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

				expected := "Test Fork 1"
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
				u.WriteString("/forks?with_body=true&parent_id=1")

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
				u.WriteString("/forks?with_body=true&creator_id=1")

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

				expected := "Test Story 1"
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
				u.WriteString(fmt.Sprintf(
					"/forks?with_body=true&title=%s",
					url.QueryEscape("story"),
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&description=%s", url.QueryEscape("girl")))

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
					u.WriteString("/forks?with_body=true&is_published=true")

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

					expected := 8
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
					u.WriteString("/forks?with_body=true&is_published=false")

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
					u.WriteString(fmt.Sprintf(
						"/forks?with_body=true&published_start=%s",
						url.QueryEscape("2019-01-01 00:00:00-0600"),
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

					expected := 8
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
					u.WriteString(fmt.Sprintf("/forks?with_body=true&published_end=%s", url.QueryEscape("2019-04-27 00:00:00-0600")))

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

					expected := 8
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
						"/forks?with_body=true&published_start=%s&published_end=%s",
						url.QueryEscape("2019-01-01 00:00:00-0600"),
						url.QueryEscape("2019-01-27 00:00:00-0600"),
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

					expected := 8
					if forksResp.Forks[0].Id != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
						return
					}
				})
			})
		})

		t.Run("As editor, can get forks with editorship or below", func(t *testing.T) {
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": 4},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

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
				u.WriteString("/forks?with_body=true")

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
				u.WriteString("/forks?with_body=true&id=4")

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
				u.WriteString("/forks?with_body=true&parent_id=1")

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
				u.WriteString("/forks?with_body=true&creator_id=2")

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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&title=%s", url.QueryEscape("story")))

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

				expected := 4
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&description=%s", url.QueryEscape("boy")))

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
					u.WriteString("/forks?with_body=true&is_published=true")

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
					u.WriteString("/forks?with_body=true&is_published=false")

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

					expected := 9
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
					u.WriteString(fmt.Sprintf(
						"/forks?with_body=true&published_start=%s",
						url.QueryEscape("2019-03-01 00:00:00-0600"),
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
					u.WriteString(fmt.Sprintf(
						"/forks?with_body=true&published_end=%s",
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
						"/forks?with_body=true&published_start=%s&published_end=%s",
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
		})

		t.Run("As proofreader, can get forks with proofreadership and below", func(t *testing.T) {
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": 5},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

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
				u.WriteString("/forks?with_body=true")

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
				u.WriteString("/forks?with_body=true&id=3")

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
				u.WriteString("/forks?with_body=true&parent_id=1")

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
				u.WriteString("/forks?with_body=true&creator_id=3")

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

				expected := "Test story 4"
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&title=%s", url.QueryEscape("story")))

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

				expected := 7
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&description=%s", url.QueryEscape("chin")))

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

				expected := 7
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
					u.WriteString("/forks?with_body=true&is_published=true")

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

					expected := 8
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
					u.WriteString("/forks?with_body=true&is_published=false")

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

					expected := 7
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
					u.WriteString(fmt.Sprintf(
						"/forks?with_body=true&published_start=%s",
						url.QueryEscape("2019-01-01 00:00:00-0600"),
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

					expected := 8
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
					u.WriteString(fmt.Sprintf(
						"/forks?with_body=true&published_end=%s",
						url.QueryEscape("2019-01-27 00:00:00-0600"),
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

					expected := 8
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
						"/forks?with_body=true&published_start=%s&published_end=%s",
						url.QueryEscape("2019-01-01 00:00:00-0600"),
						url.QueryEscape("2019-01-27 00:00:00-0600"),
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

					expected := 8
					if forksResp.Forks[0].Id != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
						return
					}
				})
			})
		})

		t.Run("As reader, can get forks with readership", func(t *testing.T) {
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": 6},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

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
				u.WriteString("/forks?with_body=true")

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
				u.WriteString("/forks?with_body=true&id=3")

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
				u.WriteString("/forks?with_body=true&parent_id=1")

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
				u.WriteString("/forks?with_body=true&creator_id=1")

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

				expected := "Test Story 1"
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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&title=%s", url.QueryEscape("story")))

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
				u.WriteString(fmt.Sprintf("/forks?with_body=true&description=%s", url.QueryEscape("girl")))

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
					u.WriteString("/forks?with_body=true&is_published=true")

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

					expected := 8
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
					u.WriteString("/forks?with_body=true&is_published=false")

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
					u.WriteString(fmt.Sprintf(
						"/forks?with_body=true&published_start=%s",
						url.QueryEscape("2019-01-01 00:00:00-0600"),
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

					expected := 8
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
					u.WriteString(fmt.Sprintf(
						"/forks?with_body=true&published_end=%s",
						url.QueryEscape("2019-01-27 00:00:00-0600"),
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

					expected := 8
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
						"/forks?with_body=true&published_start=%s&published_end=%s",
						url.QueryEscape("2019-01-01 00:00:00-0600"),
						url.QueryEscape("2019-01-27 00:00:00-0600"),
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

					expected := 8
					if forksResp.Forks[0].Id != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
						return
					}
				})
			})
		})
	})
}

func TestForksController_Create(t *testing.T) {
	setupEnvironment(t)

	handler := alice.New(
		middlewares.DatabaseMiddleware,
		middlewares.AuthenticationtMiddleware,
	).Then(middlewares.RunAPI(controllers.ForksController{}.Create))

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	client := &http.Client{}

	var u bytes.Buffer
	u.WriteString(string(testServer.URL))
	u.WriteString("/forks")

	t.Run("Unauthenticated user denied access", func(t *testing.T) {
		req, err := http.NewRequest("POST", u.String(), nil)

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
		userId := 1
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Empty request body returns error", func(t *testing.T) {
			req, err := http.NewRequest("POST", u.String(), nil)
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

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected %v, got %v\n%v", http.StatusBadRequest, resp.StatusCode, bodyStr)
				return
			}

			expected := `{"error":"invalid request body"}`
			if bodyStr != expected {
				t.Errorf("expected %v, got %v", expected, bodyStr)
				return
			}
		})

		t.Run("Can create top-level fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
				return
			}

			var forkResp controllers.ForksControllerForkResponse
			err = json.Unmarshal(body, &forkResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if forkResp.Fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, forkResp.Fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", forkResp.Fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Can create fork from fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 8, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
				return
			}

			var forkResp controllers.ForksControllerForkResponse
			err = json.Unmarshal(body, &forkResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if forkResp.Fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, forkResp.Fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", forkResp.Fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 3
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})
	})

	t.Run("As owner", func(t *testing.T) {
		userId := 2
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can create top-level fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
				return
			}

			var forkResp controllers.ForksControllerForkResponse
			err = json.Unmarshal(body, &forkResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if forkResp.Fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, forkResp.Fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", forkResp.Fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Can create fork from owned fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
				return
			}

			var forkResp controllers.ForksControllerForkResponse
			err = json.Unmarshal(body, &forkResp)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if forkResp.Fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, forkResp.Fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", forkResp.Fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Cannot create fork from unowned fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 8, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As author", func(t *testing.T) {
		userId := 3
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can create top-level fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Can create fork from authored fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 2, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 2
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Cannot create fork from unauthored fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 7, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As editor", func(t *testing.T) {
		userId := 4
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can create top-level fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Cannot create fork from edited fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})

		t.Run("Cannot create fork from unedited fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 7, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As proofreader", func(t *testing.T) {
		userId := 5
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can create top-level fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Cannot create fork from proofread fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 7, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})

		t.Run("Cannot create fork from unproofread fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As reader", func(t *testing.T) {
		userId := 6
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can create top-level fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			expectedId := 10
			if fork.Id != expectedId {
				t.Errorf("expected %v, got %v", expectedId, fork.Id)
				return
			}

			db := openDB()
			defer closeDB(db)
			userRoleLinksStore := models.NewUserRoleLinkStore(db, logger)
			links, err := userRoleLinksStore.GetLinksForResource("fork", fork.Id)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(links.Links) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(links.Links))
				return
			}
		})

		t.Run("Cannot create fork from read fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 2, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})

		t.Run("Cannot create fork from unread fork", func(t *testing.T) {
			seedDb()

			forkStr := `{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})
}

func TestForksController_Update(t *testing.T) {
	setupEnvironment(t)

	handler := alice.New(
		middlewares.DatabaseMiddleware,
		middlewares.AuthenticationtMiddleware,
	).Then(middlewares.RunAPI(controllers.ForksController{}.Update))

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	client := &http.Client{}

	var u bytes.Buffer
	u.WriteString(string(testServer.URL))
	u.WriteString("/forks")

	t.Run("Unauthenticated user denied access", func(t *testing.T) {
		req, err := http.NewRequest("PUT", u.String(), nil)

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
		userId := 1
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Empty request body returns error", func(t *testing.T) {
			req, err := http.NewRequest("PUT", u.String(), nil)
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

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected %v, got %v\n%v", http.StatusBadRequest, resp.StatusCode, bodyStr)
				return
			}

			expected := `{"error":"invalid request body"}`
			if bodyStr != expected {
				t.Errorf("expected %v, got %v", expected, bodyStr)
				return
			}
		})

		t.Run("Can update top-level fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 1, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("PUT", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			if fork.Title != title {
				t.Errorf("expected %v, got %v", title, fork.Title)
				return
			}
		})

		t.Run("Can update fork from fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 3, "ParentId": 8, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("PUT", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			if fork.Title != title {
				t.Errorf("expected %v, got %v", title, fork.Title)
				return
			}
		})
	})

	t.Run("As owner", func(t *testing.T) {
		userId := 2
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can update owned fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 4, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("PUT", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			if fork.Title != title {
				t.Errorf("expected %v, got %v", title, fork.Title)
				return
			}
		})

		t.Run("Cannot update unowned fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 2, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("PUT", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As author", func(t *testing.T) {
		userId := 3
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can update fork authored fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 2, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			if fork.Title != title {
				t.Errorf("expected %v, got %v", title, fork.Title)
				return
			}
		})

		t.Run("Cannot update fork from unauthored fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 7, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As editor", func(t *testing.T) {
		userId := 4
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can update fork edited fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 4, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			var fork models.Fork
			err = json.Unmarshal(body, &fork)
			if err != nil {
				t.Error(err)
				return
			}

			if fork.Title != title {
				t.Errorf("expected %v, got %v", title, fork.Title)
				return
			}
		})

		t.Run("Cannot update fork from edited fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 7, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As proofreader", func(t *testing.T) {
		userId := 5
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Cannot update fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 7, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As reader", func(t *testing.T) {
		userId := 6
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Cannot update fork", func(t *testing.T) {
			seedDb()

			title := "Updated Fork"
			forkStr := fmt.Sprintf(`{"Id": 7, "Title": "%s", "Description": "This is a fork!"}`, title)
			req, err := http.NewRequest("POST", u.String(), strings.NewReader(forkStr))
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})
}

func TestForksController_Delete(t *testing.T) {
	setupEnvironment(t)

	r := mux.NewRouter()
	r.Handle("/forks/{id}", alice.New(
		middlewares.DatabaseMiddleware,
		middlewares.AuthenticationtMiddleware,
	).Then(middlewares.RunAPI(controllers.ForksController{}.Delete)))

	testServer := httptest.NewServer(r)
	defer testServer.Close()

	client := &http.Client{}

	t.Run("Unauthenticated user denied access", func(t *testing.T) {
		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/forks/1")

		req, err := http.NewRequest("DELETE", u.String(), nil)

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
		userId := 1
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can delete top-level fork", func(t *testing.T) {
			seedDb()

			forkId := 1

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString("/forks/1")
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("expected %v, got %v\n%v", http.StatusNoContent, resp.StatusCode, bodyStr)
				return
			}

			db := openDB()
			defer closeDB(db)
			forkStore := models.NewForkStore(db, logger)
			_, err = forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
			if err != nil {
				t.Error(err)
				return
			}
		})

		t.Run("Can delete fork from fork", func(t *testing.T) {
			seedDb()

			forkId := 8

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString("/forks/1")
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("expected %v, got %v\n%v", http.StatusNoContent, resp.StatusCode, bodyStr)
				return
			}

			db := openDB()
			defer closeDB(db)
			forkStore := models.NewForkStore(db, logger)
			_, err = forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
			if err != nil {
				t.Error(err)
				return
			}
		})
	})

	t.Run("As owner", func(t *testing.T) {
		userId := 2
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can delete fork from owned fork", func(t *testing.T) {
			seedDb()

			forkId := 4

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("expected %v, got %v\n%v", http.StatusNoContent, resp.StatusCode, bodyStr)
				return
			}

			db := openDB()
			defer closeDB(db)
			forkStore := models.NewForkStore(db, logger)
			_, err = forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
			if err != nil {
				t.Error(err)
				return
			}
		})

		t.Run("Cannot delete fork from unowned fork", func(t *testing.T) {
			seedDb()

			forkId := 8

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As author", func(t *testing.T) {
		userId := 3
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Can delete fork from authored fork", func(t *testing.T) {
			seedDb()

			forkId := 2

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("expected %v, got %v\n%v", http.StatusNoContent, resp.StatusCode, bodyStr)
				return
			}

			db := openDB()
			defer closeDB(db)
			forkStore := models.NewForkStore(db, logger)
			_, err = forkStore.GetFork(models.NewForkQueryOptions().Id(forkId))
			if err != nil {
				t.Error(err)
				return
			}
		})

		t.Run("Cannot delete fork from unauthored fork", func(t *testing.T) {
			seedDb()

			forkId := 7

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As editor", func(t *testing.T) {
		userId := 4
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Cannot delete fork from edited fork", func(t *testing.T) {
			seedDb()

			forkId := 4

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})

		t.Run("Cannot delete fork from unedited fork", func(t *testing.T) {
			seedDb()

			forkId := 7

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As proofreader", func(t *testing.T) {
		userId := 5
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Cannot delete fork from proofread fork", func(t *testing.T) {
			seedDb()

			forkId := 7

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})

		t.Run("Cannot delete fork from unproofread fork", func(t *testing.T) {
			seedDb()

			forkId := 4

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})

	t.Run("As reader", func(t *testing.T) {
		userId := 6
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": userId},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Cannot delete fork from read fork", func(t *testing.T) {
			seedDb()

			forkId := 2

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})

		t.Run("Cannot delete fork from unread fork", func(t *testing.T) {
			seedDb()

			forkId := 4

			var u bytes.Buffer
			u.WriteString(string(testServer.URL))
			u.WriteString(fmt.Sprintf("/forks/%d", forkId))
			req, err := http.NewRequest("DELETE", u.String(), nil)
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

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})
	})
}
