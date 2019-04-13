package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"testing"
)

func TestUserRoleLinkStore_GetLinksForUser(t *testing.T) {
	seedDb()

	userRoleLinkStore := models.NewUserRoleLinkStore(db, logger)
	links, err := userRoleLinkStore.GetLinksForUser(2)
	if err != nil {
		t.Error(err)
		return
	}

	expected := 1
	if len(links.Links) != expected {
		t.Errorf("expected %v, got %v", expected, len(links.Links))
		return
	}
}

func TestUserRoleLinkStore_CreateLinks(t *testing.T) {
	var err error
	seedDb()

	userId := 3

	userRoleLinkStore := models.NewUserRoleLinkStore(db, logger)
	var links *models.UserRoleLinks
	links, err = userRoleLinkStore.GetLinksForUser(userId)
	if err != nil {
		t.Error(err)
		return
	}

	expectedLinkCount := 1
	if len(links.Links) != expectedLinkCount {
		t.Errorf("expected %v, got %v", expectedLinkCount, len(links.Links))
		return
	}

	err = userRoleLinkStore.CreateLinks(
		[]models.UserRoleLink{
			{
				UserId:       userId,
				RoleId:       1,
				ResourceType: "application",
				ResourceId:   0,
			},
		})
	if err != nil {
		t.Error(err)
		return
	}

	links, err = userRoleLinkStore.GetLinksForUser(userId)
	if err != nil {
		t.Error(err)
		return
	}

	expectedLinkCount = 2
	if len(links.Links) != expectedLinkCount {
		t.Errorf("expected %v, got %v", expectedLinkCount, len(links.Links))
		return
	}
}

func TestUserRoleLinkStore_DeleteLinks(t *testing.T) {
	var err error
	seedDb()

	userId := 3

	userRoleLinkStore := models.NewUserRoleLinkStore(db, logger)
	var links *models.UserRoleLinks
	links, err = userRoleLinkStore.GetLinksForUser(userId)
	if err != nil {
		t.Error(err)
		return
	}

	expectedLinkCount := 1
	if len(links.Links) != expectedLinkCount {
		t.Errorf("expected %v, got %v", expectedLinkCount, len(links.Links))
		return
	}

	err = userRoleLinkStore.DeleteLinks(
		[]models.UserRoleLink{
			{
				UserId:       userId,
				RoleId:       3,
				ResourceType: "fork",
				ResourceId:   0,
			},
		})
	if err != nil {
		t.Error(err)
		return
	}

	links, err = userRoleLinkStore.GetLinksForUser(userId)
	if err != nil {
		t.Error(err)
		return
	}

	expectedLinkCount = 0
	if len(links.Links) != expectedLinkCount {
		t.Errorf("expected %v, got %v", expectedLinkCount, len(links.Links))
		return
	}
}
