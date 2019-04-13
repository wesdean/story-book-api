package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/middlewares"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUserRolesController_Index(t *testing.T) {
	t.Run("Get all roles", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).Then(middlewares.RunAPI(controllers.UserRolesController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/user_roles")

		req, err := http.NewRequest("GET", u.String(), nil)

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

		var rolesResp controllers.UserRolesResponse
		err = json.Unmarshal(body, &rolesResp)
		if err != nil {
			t.Error(err)
			return
		}

		expected := 6
		if len(rolesResp.Roles) != expected {
			t.Errorf("expected %v, got %v", expected, len(rolesResp.Roles))
			return
		}
	})

	t.Run("Get roles by ID", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).Then(middlewares.RunAPI(controllers.UserRolesController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/user_roles?id=3")

		req, err := http.NewRequest("GET", u.String(), nil)

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

		var rolesResp controllers.UserRolesResponse
		err = json.Unmarshal(body, &rolesResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(rolesResp.Roles) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(rolesResp.Roles))
			return
		}

		expected := "author"
		if rolesResp.Roles[0].Name != expected {
			t.Errorf("expected %v, got %v", expected, rolesResp.Roles[0].Name)
			return
		}
	})

	t.Run("Get roles by name", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).Then(middlewares.RunAPI(controllers.UserRolesController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/user_roles?name=owner")

		req, err := http.NewRequest("GET", u.String(), nil)

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

		var rolesResp controllers.UserRolesResponse
		err = json.Unmarshal(body, &rolesResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(rolesResp.Roles) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(rolesResp.Roles))
			return
		}

		expected := 2
		if rolesResp.Roles[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, rolesResp.Roles[0].Id)
			return
		}
	})

	t.Run("Get roles by label", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).Then(middlewares.RunAPI(controllers.UserRolesController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/user_roles?label=Owner")

		req, err := http.NewRequest("GET", u.String(), nil)

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

		var rolesResp controllers.UserRolesResponse
		err = json.Unmarshal(body, &rolesResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(rolesResp.Roles) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(rolesResp.Roles))
			return
		}

		expected := 2
		if rolesResp.Roles[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, rolesResp.Roles[0].Id)
			return
		}
	})

	t.Run("Get roles by description", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).Then(middlewares.RunAPI(controllers.UserRolesController{}.Index))

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString(fmt.Sprintf("/user_roles?description=%s", url.QueryEscape("Superman Only")))

		req, err := http.NewRequest("GET", u.String(), nil)

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

		var rolesResp controllers.UserRolesResponse
		err = json.Unmarshal(body, &rolesResp)
		if err != nil {
			t.Error(err)
			return
		}

		expectedCount := 1
		if len(rolesResp.Roles) != expectedCount {
			t.Errorf("expected %v, got %v", expectedCount, len(rolesResp.Roles))
			return
		}

		expected := 1
		if rolesResp.Roles[0].Id != expected {
			t.Errorf("expected %v, got %v", expected, rolesResp.Roles[0].Id)
			return
		}
	})
}
