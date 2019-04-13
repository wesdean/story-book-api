package integration_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(config.IntegrationTest.ApiUrl)
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

	//todo fix logging Unix syslog delivery error
	expected := `{"authTokenCheck":true,"dbCheck":true,"healthCheck":true}`
	if bodyStr != expected {
		t.Errorf("expected %v, got %v", expected, bodyStr)
		return
	}
}
