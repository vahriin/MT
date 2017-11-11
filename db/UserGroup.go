package db

import (
	"github.com/vahriin/MT/model"
	"database/sql"
	"log"
)

func (cdb CacheDB) GetGroupsByUserId(id model.Id) (*[]model.Group, error) {
	groupsId, err := cdb.adb.getGroupsIdByUserId(id)
	if err != nil {
		return nil, err
	}

	var groups []model.Group
	for groupId := range *groupsId {
		/*try to get currentGroup from cache*/

		currentGroup, err := cdb.adb.getGroupById(model.Id(groupId))
		if err != nil {
			return nil, err
		}

		/*push currentGroup to cache*/

		groups = append(groups, *currentGroup)
	}

	return &groups, nil
}

func (cdb CacheDB) GetUserByGroupId(id model.Id) (*[]model.User, error) {
	usersId, err := cdb.adb.getUsersIdByGroupId(id)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for userId := range *usersId {
		/*try to get currentUser from cache*/

		currentUser, err := cdb.adb.getUserById(model.Id(userId))
		if err != nil {
			return nil, err
		}

		/*push currentUser to cache*/

		users = append(users, *currentUser)
	}

	return &users, nil
}

func (cdb CacheDB) AddUserGroup(userId, groupId model.Id) error {
	_, err := cdb.adb.db.Exec(`INSERT INTO app_user_group (user_id, group_id)
		VALUES ($1, $2);`, userId, groupId)
	if err != nil {
		log.Println("addUserToGroup returned this message: " + err.Error())
		return ErrInternal
	}
	return nil
}

func (cdb CacheDB) DelUserGroup(userId, groupId model.Id) error {
	result, err := cdb.adb.db.Exec(`DELETE FROM app_user_group WHERE user_id=$1 AND group_id=$2;`,
		userId, groupId)
	if err != nil {
		log.Println("deleteSubtransactionPack returned this message: " + err.Error())
		return ErrInternal
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		return ErrNotFound
	}
	return nil
}



func (adb AppDB) getGroupsIdByUserId(id model.Id) (*[]model.Id, error) {
	rows, err := adb.db.Query(`SELECT group_id FROM app_user_group WHERE user_id=$1;`, id)
	if err != nil {
		return nil, ErrInternal
	}

	var groupsId []model.Id
	for rows.Next() {
		var currentGroupId model.Id

		if err:= rows.Scan(&currentGroupId); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getGroupsByUserId returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}
		groupsId = append(groupsId, currentGroupId)
	}

	return &groupsId, nil
}

func (adb AppDB) getUsersIdByGroupId(id model.Id) (*[]model.Id, error) {
	rows, err := adb.db.Query(`SELECT user_id FROM app_user_group WHERE group_id=$1;`, id)
	if err != nil {
		return nil, ErrInternal
	}

	var usersId []model.Id
	for rows.Next() {
		var currentUserId model.Id

		if err:= rows.Scan(&currentUserId); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getUsersByGrupId returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}
		usersId = append(usersId, currentUserId)
	}

	return &usersId, nil
}

