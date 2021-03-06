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
	"os"
	"strings"
	"testing"
)

func authenticationTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		utils.EncodeJSON(rw, "Authentication successful")
	}
	return http.HandlerFunc(fn)
}

func TestAuthenticationtMiddleware(t *testing.T) {
	setupEnvironment(t)

	t.Run("Successful authentication", func(t *testing.T) {
		authHandler := alice.New(
			middlewares.DatabaseMiddleware,
			middlewares.AuthenticationtMiddleware,
		).Then(middlewares.RunAPI(authenticationTestHandler()))

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

		expected := `"Authentication successful"`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})

	t.Run("Missing authorization header return 401", func(t *testing.T) {
		authHandler := alice.New(middlewares.DatabaseMiddleware).Then(middlewares.AuthenticationtMiddleware(authenticationTestHandler()))

		testServer := httptest.NewServer(authHandler)
		defer testServer.Close()

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/")

		resp, err := http.Get(u.String())
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

		expected := `{"error":"invalid authentication token"}`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})

	t.Run("Authorization timeout returns 401", func(t *testing.T) {
		err := os.Setenv("AUTH_TIMEOUT", "-3")
		if err != nil {
			t.Error(t)
			return
		}

		authHandler := middlewares.AuthenticationtMiddleware(authenticationTestHandler())

		testServer := httptest.NewServer(authHandler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/")

		authToken, err := utils.CreateJWTToken(jwt.MapClaims{}, []byte(""))

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

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
			return
		}

		expected := `{"error":"authentication has expired"}`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})
}
