package middlewares_test

import (
	"bytes"
	"github.com/wesdean/story-book-api/middlewares"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func configTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		utils.EncodeJSON(rw, "Configuration successful")
	}
	return http.HandlerFunc(fn)
}

func TestConfigMiddleware(t *testing.T) {
	setupEnvironment(t)

	t.Run("Successful configuration", func(t *testing.T) {
		loggingHandler := middlewares.ConfigMiddleware(configTestHandler())

		testServer := httptest.NewServer(loggingHandler)
		defer testServer.Close()

		client := &http.Client{}

		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString("/")

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
			return
		}

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

		expected := `"Configuration successful"`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})
}
