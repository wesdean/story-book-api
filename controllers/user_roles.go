package controllers

import (
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/utils"
	"net/http"
	"strconv"
)

type UserRolesController struct{}

type UserRolesResponse struct {
	Roles []*models.UserRole
}

func (controller UserRolesController) Index(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	paramId, err := strconv.Atoi(queryParams.Get("id"))
	if err != nil {
		paramId = 0
	}

	stores, err := models.GetStoresFromRequest(r)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := struct {
		Id          int
		Name        string
		Label       string
		Description string
	}{
		Id:          paramId,
		Name:        queryParams.Get("name"),
		Label:       queryParams.Get("label"),
		Description: queryParams.Get("description"),
	}

	options := models.NewUserRoleQueryOptions()
	if params.Id > 0 {
		options.Id(params.Id)
	}
	if params.Name != "" {
		options.Name(params.Name)
	}
	if params.Label != "" {
		options.Label(params.Label)
	}
	if params.Description != "" {
		options.Description(params.Description)
	}

	roles, err := stores.UserRoleStore.GetRoles(options)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UserRolesResponse{Roles: roles}
	utils.EncodeJSON(w, response)
}
