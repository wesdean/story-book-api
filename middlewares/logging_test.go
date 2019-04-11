package middlewares_test

import (
	"bytes"
	"github.com/justinas/alice"
	"github.com/wesdean/story-book-api/middlewares"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func loggingTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		utils.EncodeJSON(rw, "Logging successful")
	}
	return http.HandlerFunc(fn)
}

func TestLoggingMiddleware(t *testing.T) {
	setupEnvironment(t)

	t.Run("Successful logging", func(t *testing.T) {
		loggingHandler := alice.New(
			middlewares.ConfigMiddleware,
		).Then(middlewares.LoggingMiddleware(loggingTestHandler()))

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

		expected := `"Logging successful"`
		if bodyStr != expected {
			t.Errorf("expected %v, got %v", expected, bodyStr)
			return
		}
	})
}
