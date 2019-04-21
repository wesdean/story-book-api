package models_test

import (
	"github.com/wesdean/story-book-api/database/models"
	"testing"
)

func TestUserRoleLinkStore_GetLinksForUserByResourceType(t *testing.T) {
	seedDb()

	userRoleLinkStore := models.NewUserRoleLinkStore(db, logger)
	links, err := userRoleLinkStore.GetLinksForUserByResourceType(2, "fork")
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

func TestUserRoleLinkStore_GetLinksForResource(t *testing.T) {
	seedDb()

	userRoleLinkStore := models.NewUserRoleLinkStore(db, logger)
	links, err := userRoleLinkStore.GetLinksForResource("fork", 3)
	if err != nil {
		t.Error(err)
		return
	}

	expected := 2
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
	links, err = userRoleLinkStore.GetLinksForUserByResourceType(userId, "fork")
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
				ResourceType: "fork",
				ResourceId:   1,
			},
		})
	if err != nil {
		t.Error(err)
		return
	}

	links, err = userRoleLinkStore.GetLinksForUserByResourceType(userId, "fork")
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

func TestUserRoleLinkStore_CopyLinksForResource(t *testing.T) {
	seedDb()

	userRoleLinkStore := models.NewUserRoleLinkStore(db, logger)
	err := userRoleLinkStore.CopyLinksForResource("fork", 3, 5)
	if err != nil {
		t.Error(err)
		return
	}

	links, err := userRoleLinkStore.GetLinksForResource("fork", 5)
	if err != nil {
		t.Error(err)
		return
	}

	expected := 2
	if len(links.Links) != expected {
		t.Errorf("expected %v, got %v", expected, len(links.Links))
		return
	}
}

func TestUserRoleLinkStore_DeleteLinks(t *testing.T) {
	var err error
	seedDb()

	userId := 3

	userRoleLinkStore := models.NewUserRoleLinkStore(db, logger)
	var links *models.UserRoleLinks
	links, err = userRoleLinkStore.GetLinksForUserByResourceType(userId, "fork")
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
				ResourceId:   2,
			},
		})
	if err != nil {
		t.Error(err)
		return
	}

	links, err = userRoleLinkStore.GetLinksForUserByResourceType(userId, "fork")
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
