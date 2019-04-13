package integration_test

import (
	"encoding/json"
	"github.com/wesdean/story-book-api/controllers"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestAuthentication(t *testing.T) {
	t.Run("POST /authentication", func(t *testing.T) {
		var baseUrl = config.IntegrationTest.ApiUrl + "/authentication"

		t.Run("Successful authentication", func(t *testing.T) {
			seedDb()

			resp, err := netClient.Post(baseUrl, "application/json", strings.NewReader(`{"username":"owner","password":"ownerpassword"}`))
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

			var authResp controllers.AuthenticationCreateTokenResponse
			err = json.Unmarshal(body, &authResp)
			if err != nil {
				t.Error(err)
				return
			}

			if authResp.Token == "" {
				t.Error("expected non-empty string, got empty string")
				return
			}
		})

		t.Run("Return EOF when body is empty", func(t *testing.T) {
			seedDb()

			resp, err := http.Post(baseUrl, "application/json", nil)
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

			expected := `{"error":"EOF"}`
			if bodyStr != expected {
				t.Errorf("expected %v, got %v", expected, bodyStr)
				return
			}
		})

		t.Run("Return 401 for bad username or password", func(t *testing.T) {
			seedDb()

			resp, err := http.Post(baseUrl, "application/json", strings.NewReader(`{"username":"test_user","password":"password"}`))
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

			expected := `{"error":"Incorrect username or password"}`
			if bodyStr != expected {
				t.Errorf("expected %v, got %v", expected, bodyStr)
				return
			}
		})
	})
}
