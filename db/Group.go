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

	defer addTx.Commit()

	if _, err := addTx.Exec(`INSERT INTO app_group (name, creator) VALUES ($1, $2);`,
		group.Name, group.Creator); err != nil {
		log.Println("addGroup returned this message: " + err.Error())
		addTx.Rollback()
		return ErrInternal
	}

	if _, err := addTx.Exec(`INSERT INTO app_user_group (user_id, group_id) VALUES ($1, $2);`,
		group.Creator, group.Id); err != nil {
		log.Println("addGroup returned this message: " + err.Error())
		addTx.Rollback()
		return ErrInternal
	}

	return nil
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

	var creatorGroups []model.Group
	for rows.Next() {
		var currentGroup model.Group

		if err := rows.Scan(&currentGroup.Id, &currentGroup.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getMainTransactionsPack returned this message: " + err.Error())
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

	defer delTx.Commit()

	if err := checkCreatorOnly(delTx, id); err == nil {
		if _, err := delTx.Exec("DELETE FROM app_user_group WHERE group_id=$1;", id); err != nil {
			delTx.Rollback()
			return ErrInternal
		}
	} else {
		delTx.Rollback()
		return err
	}

	result, err := delTx.Exec("DELETE FROM app_group WHERE group_id=$1;", id)
	if err != nil {
		log.Println("deleteGroupById returned this message: " + err.Error())
		delTx.Rollback()
		return ErrInternal
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		delTx.Rollback()
		return ErrNotFound
	}

	return err
}

func checkCreatorOnly(tx *sql.Tx, id model.Id) error {
	row := tx.QueryRow("SELECT COUNT(user_id) FROM app_user_group WHERE group_id=$1;", id)
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
