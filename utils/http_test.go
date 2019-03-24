package utils_test

import (
	"github.com/andreyvit/diff"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncodeJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		utils.EncodeJSON(w, map[string]bool{"test": true})
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := strings.Trim(string(body), "\n")

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %v", resp.StatusCode)
		return
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %v", resp.Header.Get("Content-Type"))
		return
	}

	expected := `{"test":true}`
	if bodyStr != expected {
		t.Errorf("result not as expected:\n%v", diff.CharacterDiff(expected, bodyStr))
		return
	}
}

func TestEncodeJSONError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		utils.EncodeJSONError(w, "Test error", http.StatusInternalServerError)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := strings.Trim(string(body), "\n")

	if resp.StatusCode != 500 {
		t.Errorf("expected 200, got %v", resp.StatusCode)
		return
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %v", resp.Header.Get("Content-Type"))
		return
	}

	expected := `{"error":"Test error"}`
	if bodyStr != expected {
		t.Errorf("result not as expected:\n%v", diff.CharacterDiff(expected, bodyStr))
		return
	}
}
