package middlewares_test

import (
	"bytes"
	"github.com/dgrijalva/jwt-go"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/middlewares"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func authorizationTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		utils.EncodeJSON(rw, "Authorization successful")
	}
	return http.HandlerFunc(fn)
}

func TestAuthorizationMiddleware(t *testing.T) {
	setupEnvironment(t)

	t.Run("Successful authorization", func(t *testing.T) {
		authHandler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
			middlewares.AuthorizationMiddleware,
		).Then(middlewares.RunAPI(authorizationTestHandler()))

		testServer := httptest.NewServer(authHandler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/")

		authToken, err := utils.CreateJWTToken(jwt.MapClaims{"user_id": 2}, []byte(""))

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", authToken)

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

		expected := `"Authorization successful"`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})
}
