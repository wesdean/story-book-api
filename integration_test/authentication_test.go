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
	var baseUrl = config.IntegrationTest.ApiUrl + "/authentication"

	t.Run("Successful authentication", func(t *testing.T) {
		seedDb()

		resp, err := http.Post(baseUrl, "application/json", strings.NewReader(`{"username":"owner","password":"ownerpassword"}`))
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
}
