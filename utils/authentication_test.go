package utils_test

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wesdean/story-book-api/utils"
	"testing"
)

func TestCreateJWTToken(t *testing.T) {
	token, err := utils.CreateJWTToken(jwt.MapClaims{}, []byte("testing"))
	if err != nil {
		t.Error(err)
		return
	}
	if token == "" {
		t.Error("expected non-empty string,got empty string")
		return
	}
}
