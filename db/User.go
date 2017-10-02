package db

import (
	"database/sql"
	"errors"
	"github.com/vahriin/MT/model"
)

func (adb *AppDB) GetUserById(id model.Id) (*model.User, error) {
	row := adb.db.QueryRow("SELECT nick FROM users WHERE id=$1", id)
	user := new(model.User)
	user.Id = id
	if err := row.Scan(&user.Nick); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, errors.New("no user in DB")
		case err != nil:
			return nil, err
		}
	}
	return user, nil
}

func (adb *AppDB) GetUserByEmail(email string) (*model.User, error) {
	row := adb.db.QueryRow("SELECT id, nick FROM users WHERE email=$1", email)
	user := new(model.User)
	if err := row.Scan(&user.Id, &user.Nick); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, errors.New("no user in DB")
		case err != nil:
			return nil, err
		}
	}
	return user, nil
}
