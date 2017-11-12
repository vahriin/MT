package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
	"log"
)

func (cdb CacheDB) GetUserById(id model.Id) (*model.User, error) {
	/* Add search in Cache */

	return cdb.adb.getUserById(id)
}

func (cdb CacheDB) GetUserByGoogleId(googleId []byte) (*model.User, error) {
	/* Add search in Cache */

	return cdb.adb.getUserByGoogleId(googleId)
}

func (cdb CacheDB) AddUser(user *model.GoogleUser) error {

	return cdb.adb.addUser(user)
}

func (cdb CacheDB) DeleteUser(userId model.Id) error {
	return cdb.adb.deleteUser(userId)
}

func (adb AppDB) addUser(user *model.GoogleUser) error {
	if _, err := adb.db.Exec("INSERT INTO app_user (google_id, nick) VALUES ($1, $2);",
		user.GoogleId, user.Nick); err != nil {
		log.Println("addUser returned this message: " + err.Error())
		return ErrInternal
	}
	return nil
}

func (adb AppDB) deleteUser(id model.Id) error {
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

func (adb AppDB) getUserById(id model.Id) (*model.User, error) {
	row := adb.db.QueryRow("SELECT nick FROM app_user WHERE id=$1;", id)
	var user model.User
	user.Id = id
	if err := row.Scan(&user.Nick); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return &user, ErrNotFound
		case err != nil:
			log.Println("getUserById returned this message: " + err.Error())
			return &user, ErrInternal
		}
	}
	return &user, nil
}

func (adb AppDB) getUserByGoogleId(googleId []byte) (*model.User, error) {
	row := adb.db.QueryRow("SELECT id, nick FROM app_user WHERE email=$1", googleId)
	var user model.User
	if err := row.Scan(&user.Id, &user.Nick); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return &user, ErrNotFound
		case err != nil:
			log.Println("getUserByGoogleId returned this message: " + err.Error())
			return &user, ErrInternal
		}
	}
	return &user, nil
}
