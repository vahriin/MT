package db

import (
	"github.com/vahriin/MT/model"
	"log"
)

func (cdb CacheDB) AddGoogleUser(user *model.GoogleUser) error {

	return cdb.adb.addGoogleUser(user)
}

func (cdb CacheDB) DeleteGoogleUser(userId model.Id) error {
	return cdb.adb.deleteGoogleUser(userId)
}

func (adb AppDB) addGoogleUser(user *model.GoogleUser) error {
	if _, err := adb.db.Exec("INSERT INTO app_user (google_id, nick) VALUES ($1, $2);",
		user.GoogleId, user.Nick); err != nil {
		log.Println("addUser returned this message: " + err.Error())
		return ErrInternal
	}
	return nil
}

func (adb AppDB) deleteGoogleUser(id model.Id) error {
	result, err := adb.db.Exec(`DELETE FROM app_user WHERE user_id=$1;`, id)
	if err != nil {
		log.Println("deleteUser returned this message: " + err.Error())
		return ErrInternal
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		return ErrNotFound
	}
	return nil
}