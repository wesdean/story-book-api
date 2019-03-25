package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/wesdean/story-book-api/database"
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
		utils.EncodeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var db *database.Database
	dbContext, ok := context.GetOk(r, "DB")
	if ok {
		db = dbContext.(*database.Database)
	} else {
		utils.EncodeJSONError(w, "Missing database connection", http.StatusInternalServerError)
		return
	}

	err = db.Begin()
	if err != nil {
		utils.EncodeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Rollback()

	userStore := models.NewUserStore(db)
	options := models.NewUserQueryOptions().
		Username(reqData.Username).
		Password(reqData.Password)
	user, err := userStore.GetUser(options)
	if err != nil {
		utils.EncodeJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if user == nil {
		utils.EncodeJSONError(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.CreateJWTToken(jwt.MapClaims{"user_id": user.Id}, []byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		utils.EncodeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AuthenticationCreateTokenResponse{Token: token}
	utils.EncodeJSON(w, response)
}

func (controller AuthenticationController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var err error

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		utils.EncodeJSONError(w, "Missing authorization header", http.StatusBadRequest)
		return
	}

	var db *database.Database
	dbContext, ok := context.GetOk(r, "DB")
	if ok {
		db = dbContext.(*database.Database)
	} else {
		utils.EncodeJSONError(w, "Missing database connection", http.StatusInternalServerError)
		return
	}

	err = db.Begin()
	if err != nil {
		utils.EncodeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Rollback()

	userStore := models.NewUserStore(db)
	_, err = userStore.AuthenticateUser(authHeader)
	if err != nil {
		utils.EncodeJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := AuthenticationValidateTokenResponse{IsValid: true}

	utils.EncodeJSON(w, response)
}
