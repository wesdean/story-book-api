package models

import (
	"fmt"
	"github.com/wesdean/story-book-api/database"
	"github.com/wesdean/story-book-api/logging"
	"strings"
)

type UserRoleLinkStore struct {
	*Store
}

type UserRoleLinks struct {
	Links []UserRoleLink
}

type UserRoleLink struct {
	UserId       int
	RoleId       int
	ResourceType string
	ResourceId   int
}

func NewUserRoleLinkStore(db *database.Database, logger *logging.Logger) *UserRoleLinkStore {
	return &UserRoleLinkStore{Store: NewStore(db, logger)}
}

func (store *UserRoleLinkStore) CreateLinks(roleLinks []UserRoleLink) error {
	var valuesStr []string
	var params []interface{}
	for i := 0; i < len(roleLinks); i++ {
		valuesStr = append(valuesStr, fmt.Sprintf("($%d, $%d, $%d, $%d)", i+1, i+2, i+3, i+4))
		params = append(params, roleLinks[i].UserId, roleLinks[i].RoleId, roleLinks[i].ResourceType, roleLinks[i].ResourceId)
	}

	sqlQuery := `insert into user_role_links (user_id, user_role_id, resource_type, resource_id)
		values ` + strings.Join(valuesStr, ",")
	_, err := store.db.Tx.Exec(sqlQuery, params...)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to add roles to user: %s", err.Error())
		return err
	}

	return nil
}

func (store *UserRoleLinkStore) DeleteLinks(roleLinks []UserRoleLink) error {
	var whereStr []string
	var params []interface{}
	for i := 0; i < len(roleLinks); i++ {
		whereStr = append(whereStr,
			fmt.Sprintf(
				"(user_id = $%d and user_role_id = $%d and resource_type = $%d and resource_id = $%d)",
				i+1, i+2, i+3, i+4))
		params = append(params, roleLinks[i].UserId, roleLinks[i].RoleId, roleLinks[i].ResourceType, roleLinks[i].ResourceId)
	}

	sqlQuery := `delete from user_role_links where ` + strings.Join(whereStr, " or ")
	_, err := store.db.Tx.Exec(sqlQuery, params...)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to add roles to user: %s", err.Error())
		return err
	}

	return nil
}

func (store *UserRoleLinkStore) GetLinksForUser(userId int) (*UserRoleLinks, error) {
	sqlQuery := `select user_id, user_role_id, resource_type, resource_id from user_role_links where user_id = $1`
	rows, err := store.db.Tx.Query(sqlQuery, userId)
	if err != nil {
		logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to retrieve user role links: %s", err.Error())
		return nil, err
	}

	userRoleLinks := &UserRoleLinks{}
	for rows.Next() {
		link := UserRoleLink{}
		err = rows.Scan(
			&link.UserId,
			&link.RoleId,
			&link.ResourceType,
			&link.ResourceId)
		if err != nil {
			logging.Logf(store.logger, logging.LOGLEVEL_ERROR, "failed to scan user role link: %s", err.Error())
			return nil, err
		}
		userRoleLinks.Links = append(userRoleLinks.Links, link)
	}
	return userRoleLinks, err
}
