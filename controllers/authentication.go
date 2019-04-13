package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"os"
)

type AuthenticationController struct{}

type AuthenticationCreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticationCreateTokenResponse struct {
	Token string `json:"token"`
}

type AuthenticationValidateTokenResponse struct {
	IsValid bool `json:"is_valid"`
}

func (controller AuthenticationController) CreateToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reqData AuthenticationCreateTokenRequest
	err := decoder.Decode(&reqData)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusBadRequest)
		return
	}

	stores, err := models.GetStoresFromRequest(r)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	options := models.NewUserQueryOptions().
		Username(reqData.Username).
		Password(reqData.Password)
	user, err := stores.UserStore.GetUser(options)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusUnauthorized)
		return
	}
	if user == nil {
		utils.EncodeJSONErrorWithLogging(r, w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.CreateJWTToken(jwt.MapClaims{"user_id": user.Id}, []byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AuthenticationCreateTokenResponse{Token: token}
	utils.EncodeJSON(w, response)
}

func (controller AuthenticationController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var err error

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		utils.EncodeJSONErrorWithLogging(r, w, "missing authorization header", http.StatusBadRequest)
		return
	}

	stores, err := models.GetStoresFromRequest(r)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stores.UserStore.AuthenticateUser(authHeader)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := AuthenticationValidateTokenResponse{IsValid: true}

	utils.EncodeJSON(w, response)
}
