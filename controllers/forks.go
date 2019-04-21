package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/middlewares"
	"github.com/wesdean/story-book-api/utils"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ForksController struct{}

type ForksControllerForksResponse struct {
	Forks []*models.Fork
}

type ForksControllerForksWithBodyResponse struct {
	Forks []*models.ForkWithBody
}

func (controller ForksController) Index(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := middlewares.GetAuthenticatedUserFromRequest(r)
	if authenticatedUser == nil {
		utils.EncodeJSONErrorWithLogging(r, w, "access denied", http.StatusUnauthorized)
		return
	}

	queryParams := r.URL.Query()

	paramWithBody, err := strconv.ParseBool(queryParams.Get("with_body"))
	if err != nil {
		paramWithBody = false
	}

	paramUserCanRead := 0
	if paramWithBody {
		paramUserCanRead = authenticatedUser.Id
	}

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

	paramOwner, err := strconv.Atoi(queryParams.Get("owner"))
	if err != nil {
		paramOwner = 0
	}

	paramAuthor, err := strconv.Atoi(queryParams.Get("author"))
	if err != nil {
		paramAuthor = 0
	}

	paramEditor, err := strconv.Atoi(queryParams.Get("editor"))
	if err != nil {
		paramEditor = 0
	}

	paramProofreader, err := strconv.Atoi(queryParams.Get("proofreader"))
	if err != nil {
		paramProofreader = 0
	}

	paramReader, err := strconv.Atoi(queryParams.Get("reader"))
	if err != nil {
		paramReader = 0
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
		Owner          int
		Author         int
		Editor         int
		Proofreader    int
		Reader         int
		UserCanRead    int
	}{
		Id:             paramId,
		ParentId:       paramParentId,
		CreatorId:      paramCreatorId,
		Title:          queryParams.Get("title"),
		Description:    queryParams.Get("description"),
		IsPublished:    paramIsPublished,
		PublishedStart: paramPublishedStart,
		PublishedEnd:   paramPublishedEnd,
		Owner:          paramOwner,
		Author:         paramAuthor,
		Editor:         paramEditor,
		Proofreader:    paramProofreader,
		Reader:         paramReader,
		UserCanRead:    paramUserCanRead,
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
	if params.Owner > 0 {
		options.Owner(params.Owner)
	}
	if params.Author > 0 {
		options.Author(params.Author)
	}
	if params.Editor > 0 {
		options.Editor(params.Editor)
	}
	if params.Proofreader > 0 {
		options.Proofreader(params.Proofreader)
	}
	if params.Reader > 0 {
		options.Reader(params.Reader)
	}
	if params.UserCanRead > 0 {
		options.UserCanRead(params.UserCanRead)
	}

	var response interface{}
	if paramWithBody {
		forks, err := stores.ForkStore.GetForksWithBody(options)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}

		response = ForksControllerForksWithBodyResponse{Forks: forks}
	} else {
		forks, err := stores.ForkStore.GetForks(options)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}
		response = ForksControllerForksResponse{Forks: forks}
	}

	utils.EncodeJSON(w, response)
}

func (controller ForksController) Create(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := middlewares.GetAuthenticatedUserFromRequest(r)
	if authenticatedUser == nil {
		utils.EncodeJSONErrorWithLogging(r, w, "access denied", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var fork *models.Fork
	err := decoder.Decode(&fork)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			err = errors.New("invalid request body")
		}
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusBadRequest)
		return
	}

	stores, err := models.GetStoresFromRequest(r)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	canCreate, err := stores.ForkStore.UserCanCreate(authenticatedUser.Id, fork.ParentId)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !canCreate {
		utils.EncodeJSONErrorWithLogging(r, w, "authorization denied", http.StatusUnauthorized)
		return
	}

	fork.CreatorId = authenticatedUser.Id
	err = stores.ForkStore.CreateFork(fork)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fork.ParentId == 0 {
		links := []models.UserRoleLink{
			{
				UserId:       fork.CreatorId,
				RoleId:       models.USER_ROLE_OWNER,
				ResourceType: "fork",
				ResourceId:   fork.Id,
			},
		}
		err = stores.UserRoleLinkStore.CreateLinks(links)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = stores.UserRoleLinkStore.CopyLinksForResource("fork", fork.ParentId, fork.Id)
		if err != nil {
			utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = stores.Commit()
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.EncodeJSON(w, fork)
}

func (controller ForksController) Update(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := middlewares.GetAuthenticatedUserFromRequest(r)
	if authenticatedUser == nil {
		utils.EncodeJSONErrorWithLogging(r, w, "access denied", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var fork *models.Fork
	err := decoder.Decode(&fork)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			err = errors.New("invalid request body")
		}
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusBadRequest)
		return
	}

	stores, err := models.GetStoresFromRequest(r)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	canUpdate, err := stores.ForkStore.UserCanUpdate(authenticatedUser.Id, fork.Id)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !canUpdate {
		utils.EncodeJSONErrorWithLogging(r, w, "authorization denied", http.StatusUnauthorized)
		return
	}

	fork.CreatorId = authenticatedUser.Id
	err = stores.ForkStore.UpdateFork(fork)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = stores.Commit()
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.EncodeJSON(w, fork)
}

func (controller ForksController) Delete(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := middlewares.GetAuthenticatedUserFromRequest(r)
	if authenticatedUser == nil {
		utils.EncodeJSONErrorWithLogging(r, w, "access denied", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	forkIdStr := params["id"]
	forkId, err := strconv.Atoi(forkIdStr)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, fmt.Sprintf("invalid fork id: %s", forkIdStr), http.StatusBadRequest)
		return
	}

	stores, err := models.GetStoresFromRequest(r)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	canDelete, err := stores.ForkStore.UserCanDelete(authenticatedUser.Id, forkId)
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !canDelete {
		utils.EncodeJSONErrorWithLogging(r, w, "authorization denied", http.StatusUnauthorized)
		return
	}

	_, err = stores.ForkStore.DeleteForks(models.NewForkQueryOptions().Id(forkId))
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = stores.Commit()
	if err != nil {
		utils.EncodeJSONErrorWithLogging(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
