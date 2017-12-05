package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
	"log"
)

func (cdb CacheDB) AddGroup(group *model.Group) error {
	if err := cdb.adb.addGroup(group); err != nil {
		return err
	}

	/* push group to cache */
	return nil
}

func (cdb CacheDB) AddUserToGroup(user model.Id, group model.Id) error {
	return cdb.adb.addUserToGroup(user, group)
}

func (cdb CacheDB) GetGroupMember(groupId model.Id) (*[]model.User, error) {
	return cdb.adb.getGroupMembers(groupId)
}

func (cdb CacheDB) GetGroupsByUser(userId model.Id) (*[]model.Group, error) {
	/* cache*/

	return cdb.adb.getGroupsByUser(userId)
}

func (cdb CacheDB) GetGroupById(groupId model.Id) (*model.Group, error) {
	/* get from cache */

	return cdb.adb.getGroupById(groupId)
}

func (cdb CacheDB) GetGroupsByCreator(creatorId model.Id) (*[]model.Group, error) {

	return cdb.adb.getGroupsByCreator(creatorId)
}

func (cdb CacheDB) DeleteGroupById(groupId model.Id) error {

	return cdb.adb.deleteGroupById(groupId)
}

func (adb AppDB) addGroup(group *model.Group) error {
	addTx, err := adb.db.Begin()
	if err != nil {
		log.Println("addGroup returned this message: " + err.Error())
		addTx.Rollback()
		return ErrInternal
	}

	if _, err := addTx.Exec(`INSERT INTO app_group (name, creator_id) VALUES ($1, $2);`,
		group.Name, group.Creator); err != nil {
		log.Println("addGroup returned this message: " + err.Error())
		addTx.Rollback()
		return ErrInternal
	}

	row := addTx.QueryRow("SELECT MAX(id) FROM app_group;")
	if err := row.Scan(&group.Id); err != nil {
		log.Println("addGroup returned this message: " + err.Error())
		addTx.Rollback()
		return ErrInternal
	}

	if _, err := addTx.Exec(`INSERT INTO app_user_group (user_id, gr_id) VALUES ($1, $2);`,
		group.Creator, group.Id); err != nil {
		log.Println("addGroup returned this message: " + err.Error())
		addTx.Rollback()
		return ErrInternal
	}

	addTx.Commit()

	return nil
}

func (adb AppDB) addUserToGroup(user model.Id, group model.Id) error {
	row := adb.db.QueryRow("SELECT user_id, gr_id FROM app_user_group WHERE user_id=$1 AND gr_id=$2",
		user, group)

	var tempUser, tempGroup model.Id

	err := row.Scan(&tempUser, &tempGroup)
	if err == sql.ErrNoRows {
		_, err = adb.db.Exec("INSERT INTO app_user_group (user_id, gr_id) VALUES($1, $2);", user, group)
		if err != nil {
			log.Println("addUserToGroup returned this message: " + err.Error())
			return ErrInternal
		}
	} else {
		return ErrForbidden
	}


	return nil
}

func (adb AppDB) getGroupsByUser(id model.Id) (*[]model.Group, error) {
	rows, err := adb.db.Query(`SELECT id, name, creator_id
		FROM app_group
		JOIN app_user_group
		ON app_group.id = app_user_group.gr_id
		WHERE app_user_group.user_id = $1`, id)

	if err != nil {
		return nil, ErrInternal
	}

	defer rows.Close()

	var userGroups []model.Group
	for rows.Next() {
		var currentGroup model.Group

		if err := rows.Scan(&currentGroup.Id, &currentGroup.Name, &currentGroup.Creator); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getGroupsByUser returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}

		userGroups = append(userGroups, currentGroup)
	}
	return &userGroups, nil
}

func (adb AppDB) getGroupMembers(id model.Id) (*[]model.User, error) {
	rows, err := adb.db.Query(`SELECT id, nick
		FROM app_user_group
		JOIN app_user
		ON app_user_group.user_id = app_user.id
		WHERE app_user_group.gr_id = $1;`, id)
	if err != nil {
		log.Println("getGroupsByCreator returned this message: " + err.Error())
		return nil, ErrInternal
	}

	defer rows.Close()

	var groupMembers []model.User
	for rows.Next() {
		var currentMember model.User

		if err := rows.Scan(&currentMember.Id, &currentMember.Nick); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getGroupsByCreator returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}

		groupMembers = append(groupMembers, currentMember)
	}

	return &groupMembers, nil
}

func (adb AppDB) getGroupById(id model.Id) (*model.Group, error) {
	row := adb.db.QueryRow(`SELECT name, creator_id FROM app_group WHERE id=$1;`, id)
	var group model.Group
	if err := row.Scan(&group.Name, &group.Creator); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		} else {
			log.Println("getGroupById returned this message: " + err.Error())
			return nil, ErrInternal
		}
	}

	group.Id = id

	return &group, nil
}

func (adb AppDB) getGroupsByCreator(creatorId model.Id) (*[]model.Group, error) {
	rows, err := adb.db.Query(`SELECT id, name FROM app_group WHERE creator_id=$1;`, creatorId)
	if err != nil {
		return nil, ErrInternal
	}

	defer rows.Close()

	var creatorGroups []model.Group
	for rows.Next() {
		var currentGroup model.Group

		if err := rows.Scan(&currentGroup.Id, &currentGroup.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getGroupsByCreator returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}

		currentGroup.Creator = creatorId

		creatorGroups = append(creatorGroups, currentGroup)
	}

	return &creatorGroups, nil
}

func (adb AppDB) deleteGroupById(id model.Id) error {
	delTx, err := adb.db.Begin()
	if err != nil {
		log.Println("deleteGroupById returned this message: " + err.Error())
		delTx.Rollback()
		return ErrInternal
	}

	if err := checkCreatorOnly(delTx, id); err == nil {
		if _, err := delTx.Exec("DELETE FROM app_user_group WHERE gr_id=$1;", id); err != nil {
			delTx.Rollback()
			return ErrInternal
		}
	} else {
		delTx.Rollback()
		return err
	}

	result, err := delTx.Exec("DELETE FROM app_group WHERE gr_id=$1;", id)
	if err != nil {
		log.Println("deleteGroupById returned this message: " + err.Error())
		delTx.Rollback()
		return ErrInternal
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		delTx.Rollback()
		return ErrNotFound
	}

	delTx.Commit()

	return err
}

func checkCreatorOnly(tx *sql.Tx, id model.Id) error {
	row := tx.QueryRow("SELECT COUNT(user_id) FROM app_user_group WHERE gr_id=$1;", id)
	var amountMembers int
	if err := row.Scan(&amountMembers); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		} else {
			log.Println("checkCreatorOnly returned this message: " + err.Error())
			return ErrInternal
		}
	}

	if amountMembers == 1 {
		return nil
	} else {
		return ErrForbidden
	}
}
