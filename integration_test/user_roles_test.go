package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/wesdean/story-book-api/controllers"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestUserRoles(t *testing.T) {
	t.Run("GET /user_roles", func(t *testing.T) {
		seedDb()

		var baseUrl = config.IntegrationTest.ApiUrl + "/user_roles"

		t.Run("Get all roles", func(t *testing.T) {
			resp, err := netClient.Get(baseUrl)
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

			var roles controllers.UserRolesResponse
			err = json.Unmarshal(body, &roles)
			if err != nil {
				t.Error(err)
				return
			}

			expected := 6
			if len(roles.Roles) != expected {
				t.Errorf("expected %v, got %v", expected, len(roles.Roles))
				return
			}
		})

		t.Run("Get roles by id", func(t *testing.T) {
			resp, err := netClient.Get(baseUrl + "?id=3")
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

			var roles controllers.UserRolesResponse
			err = json.Unmarshal(body, &roles)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(roles.Roles) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(roles.Roles))
				return
			}

			expectedId := 3
			if roles.Roles[0].Id != 3 {
				t.Errorf("expected %v, got %v", expectedId, roles.Roles[0].Id)
				return
			}
		})

		t.Run("Get roles by name", func(t *testing.T) {
			resp, err := netClient.Get(baseUrl + "?name=owner")
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

			var roles controllers.UserRolesResponse
			err = json.Unmarshal(body, &roles)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(roles.Roles) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(roles.Roles))
				return
			}

			expectedId := 2
			if roles.Roles[0].Id != 2 {
				t.Errorf("expected %v, got %v", expectedId, roles.Roles[0].Id)
				return
			}
		})

		t.Run("Get roles by label", func(t *testing.T) {
			resp, err := netClient.Get(baseUrl + "?label=Owner")
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

			var roles controllers.UserRolesResponse
			err = json.Unmarshal(body, &roles)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(roles.Roles) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(roles.Roles))
				return
			}

			expectedId := 2
			if roles.Roles[0].Id != 2 {
				t.Errorf("expected %v, got %v", expectedId, roles.Roles[0].Id)
				return
			}
		})

		t.Run("Get roles by description", func(t *testing.T) {
			resp, err := netClient.Get(fmt.Sprintf("%s?description=%s", baseUrl, url.QueryEscape("Superman Only")))
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

			var roles controllers.UserRolesResponse
			err = json.Unmarshal(body, &roles)
			if err != nil {
				t.Error(err)
				return
			}

			expectedCount := 1
			if len(roles.Roles) != expectedCount {
				t.Errorf("expected %v, got %v", expectedCount, len(roles.Roles))
				return
			}

			expectedId := 1
			if roles.Roles[0].Id != 1 {
				t.Errorf("expected %v, got %v", expectedId, roles.Roles[0].Id)
				return
			}
		})
	})
}
