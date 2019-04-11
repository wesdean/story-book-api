package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/middlewares"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAuthenticationController_CreateToken(t *testing.T) {
	t.Run("Return token on success", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).ThenFunc(controllers.AuthenticationController{}.CreateToken)

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/authentication")

		req, err := http.NewRequest("POST", u.String(), strings.NewReader(`{"username":"owner","password":"ownerpassword"}`))

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

	t.Run("Return EOF error when body is empty", func(t *testing.T) {
		handler := http.HandlerFunc(controllers.AuthenticationController{}.CreateToken)

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/authentication")

		req, err := http.NewRequest("POST", u.String(), nil)

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

		expected := `{"error":"EOF"}`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})

	t.Run("Return 401 for bad username or password", func(t *testing.T) {
		handler := alice.New(
			middlewares.DatabaseMiddleware,
		).ThenFunc(controllers.AuthenticationController{}.CreateToken)

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/authentication")

		req, err := http.NewRequest("POST", u.String(), strings.NewReader(`{"username":"test_user","password":"password"}`))

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

		expected := `{"error":"Incorrect username or password"}`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})
}

func TestAuthenticationController_ValidateToken(t *testing.T) {
	setupEnvironment(t)

	t.Run("Return true on valid token", func(t *testing.T) {
		handler := alice.New(
			middlewares.ConfigMiddleware,
			middlewares.LoggingMiddleware,
			middlewares.DatabaseMiddleware,
		).ThenFunc(controllers.AuthenticationController{}.ValidateToken)

		token, err := utils.CreateJWTToken(
			jwt.MapClaims{"user_id": 2},
			[]byte(os.Getenv("AUTH_SECRET")),
		)
		if err != nil {
			t.Error(err)
			return
		}

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/authentication")

		req, err := http.NewRequest("GET", u.String(), nil)
		if req == nil {
			t.Error("expected not nil, got nil")
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

		expected := `{"is_valid":true}`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})

	t.Run("Return error on missing authorization header", func(t *testing.T) {
		handler := http.HandlerFunc(controllers.AuthenticationController{}.ValidateToken)

		testServer := httptest.NewServer(handler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/authentication")

		req, err := http.NewRequest("GET", u.String(), nil)

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

		expected := `{"error":"Missing authorization header"}`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})
}
