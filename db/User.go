package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
	"log"
)

func (cdb CacheDB) GetUserById(id model.Id) (model.User, error) {
	/* Add search in Cache */

	return cdb.adb.getUserById(id)
}

func (cdb CacheDB) GetUserByEmail(email string) (model.User, error) {
	/* Add search in Cache */

	return cdb.adb.getUserByEmail(email)
}

func (adb AppDB) getUserById(id model.Id) (model.User, error) {
	row := adb.db.QueryRow("SELECT nick FROM users WHERE id=$1", id)
	var user model.User
	user.Id = id
	if err := row.Scan(&user.Nick); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return user, ErrNotFound
		case err != nil:
			log.Println("getUserById returned this message: " + err.Error())
			return user, ErrInternal
		}
	}
	return user, nil
}

func (adb AppDB) getUserByEmail(email string) (model.User, error) {
	row := adb.db.QueryRow("SELECT id, nick FROM users WHERE email=$1", email)
	var user model.User
	if err := row.Scan(&user.Id, &user.Nick); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return user, ErrNotFound
		case err != nil:
			log.Println("getUserByEmail returned this message: " + err.Error())
			return user, ErrInternal
		}
	}
	return user, nil
}
