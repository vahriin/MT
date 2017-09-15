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

func (adb *AppDB) GetUserByNick(nick string) (*model.User, error) {
	row := adb.db.QueryRow("SELECT id FROM users WHERE nick=$1", nick)
	user := new(model.User)
	user.Nick = nick
	if err := row.Scan(&user.Id); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, errors.New("no user in DB")
		case err != nil:
			return nil, err
		}
	}
	return user, nil
}
