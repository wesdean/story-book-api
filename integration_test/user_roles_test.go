package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestUserRoles(t *testing.T) {
	t.Run("GET /user_roles", func(t *testing.T) {
		seedDb()

		var baseUrl = config.IntegrationTest.ApiUrl + "/user_roles"
		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 2},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("Get all roles", func(t *testing.T) {
			req, err := http.NewRequest("GET", baseUrl, nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := netClient.Do(req)
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
			req, err := http.NewRequest("GET", baseUrl+"?id=3", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := netClient.Do(req)
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
			req, err := http.NewRequest("GET", baseUrl+"?name=owner", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := netClient.Do(req)
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
			req, err := http.NewRequest("GET", baseUrl+"?label=Owner", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := netClient.Do(req)
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
			req, err := http.NewRequest("GET", fmt.Sprintf("%s?description=%s", baseUrl, url.QueryEscape("Superman Only")), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := netClient.Do(req)
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
