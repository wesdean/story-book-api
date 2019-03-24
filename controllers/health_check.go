package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"os"
)

type HealthCheckController struct{}

func (controller HealthCheckController) Index(w http.ResponseWriter, r *http.Request) {
	var err error

	var authTokenCheck bool
	emptySecretToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.aGTWgif4pwMnjF8My859yqoueBN9ueg95F58WNFt1ps"
	token, err := utils.CreateJWTToken(jwt.MapClaims{}, []byte(os.Getenv("AUTH_SECRET")))
	if err == nil && token != emptySecretToken {
		authTokenCheck = true
	}

	var healthCheck bool
	if err == nil {
		healthCheck = true
	}

	var output = map[string]interface{}{
		"healthCheck":    healthCheck,
		"authTokenCheck": authTokenCheck,
	}

	utils.EncodeJSON(w, output)
}
