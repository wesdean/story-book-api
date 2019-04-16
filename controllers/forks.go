package controllers

import (
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/middlewares"
	"github.com/wesdean/story-book-api/utils"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"strconv"
	"time"
)

type ForksController struct{}

type ForksControllerForksResponse struct {
	Forks []*models.Fork
}

func (controller ForksController) Index(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := middlewares.GetAuthenticatedUserFromRequest(r)
	if authenticatedUser == nil {
		utils.EncodeJSONErrorWithLogging(r, w, "access denied", http.StatusUnauthorized)
		return
	}

	queryParams := r.URL.Query()
	paramId, err := strconv.Atoi(queryParams.Get("id"))
	if err != nil {
		paramId = 0
	}
	paramParentId, err := strconv.Atoi(queryParams.Get("parent_id"))
	if err != nil {
		paramParentId = 0
	}
	paramCreatorId, err := strconv.Atoi(queryParams.Get("creator_id"))
	if err != nil {
		paramCreatorId = 0
	}

	var paramIsPublished null.Bool
	paramIsPublishedBool, err := strconv.ParseBool(queryParams.Get("is_published"))
	if err == nil {
		paramIsPublished = null.BoolFrom(paramIsPublishedBool)
	} else {
		paramIsPublished = null.NewBool(false, false)
	}

	var paramPublishedStart null.Time
	paramsPublishedStartTime, err := time.Parse("2006-01-02 15:04:05-0700", queryParams.Get("published_start"))
	if err == nil {
		paramPublishedStart = null.TimeFrom(paramsPublishedStartTime)
	} else {
		paramPublishedStart = null.NewTime(time.Now(), false)
	}

	var paramPublishedEnd null.Time
	paramsPublishedEndTime, err := time.Parse("2006-01-02 15:04:05-0700", queryParams.Get("published_end"))
	if err == nil {
		paramPublishedEnd = null.TimeFrom(paramsPublishedEndTime)
	} else {
		paramPublishedEnd = null.NewTime(time.Now(), false)
	}

	stores, err := models.GetStoresFromRequest(r)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := struct {
		Id             int
		ParentId       int
		CreatorId      int
		Title          string
		Description    string
		IsPublished    null.Bool
		PublishedStart null.Time
		PublishedEnd   null.Time
	}{
		Id:             paramId,
		ParentId:       paramParentId,
		CreatorId:      paramCreatorId,
		Title:          queryParams.Get("title"),
		Description:    queryParams.Get("description"),
		IsPublished:    paramIsPublished,
		PublishedStart: paramPublishedStart,
		PublishedEnd:   paramPublishedEnd,
	}

	options := models.NewForkQueryOptions()
	if params.Id > 0 {
		options.Id(params.Id)
	} else {
		options.ParentId(params.ParentId)
	}
	if params.CreatorId > 0 {
		options.CreatorId(params.CreatorId)
	}
	if params.Title != "" {
		options.Title(params.Title)
	}
	if params.Description != "" {
		options.Description(params.Description)
	}
	if !params.IsPublished.IsZero() {
		options.IsPublished(params.IsPublished.ValueOrZero())
	}
	if !params.PublishedStart.IsZero() {
		options.PublishedStart(params.PublishedStart.ValueOrZero())
	}
	if !params.PublishedEnd.IsZero() {
		options.PublishedEnd(params.PublishedEnd.ValueOrZero())
	}

	forks, err := stores.ForkStore.GetForks(options)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ForksControllerForksResponse{Forks: forks}
	utils.EncodeJSON(w, response)
}
