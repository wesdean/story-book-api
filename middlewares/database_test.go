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

func databaseTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		utils.EncodeJSON(rw, "Database test successful")
	}
	return http.HandlerFunc(fn)
}

func TestDatabaseMiddleware(t *testing.T) {
	setupEnvironment(t)

	dbHandler := middlewares.DatabaseMiddleware(databaseTestHandler())

	testServer := httptest.NewServer(dbHandler)
	defer testServer.Close()

	client := &http.Client{}

	var u bytes.Buffer
	u.WriteString(string(testServer.URL))
	u.WriteString("/")

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

	expected := `"Database test successful"`
	if bodyStr != expected {
		t.Errorf("expected %v, got %v", expected, bodyStr)
		return
	}
}
