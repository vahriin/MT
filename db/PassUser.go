package db

import (
	"database/sql"
	"errors"
	"github.com/vahriin/MT/model"
)

func (adb AppDB) AddPassUser(pu *model.PassUser) error {
	row := adb.db.QueryRow("SELECT id FROM users WHERE nick=$1", pu.Nick)
	var err error
	if row.Scan(&pu.Id) == sql.ErrNoRows {
		_, err = adb.db.Exec("INSERT INTO users(nick, passhash) VALUES($1, $2);",
			pu.Nick, pu.PassHash)
	} else {
		err = errors.New("User " + pu.Nick + " already exists")
	}

	return err
}


func (adb AppDB) GetPassUser(nick string) (*model.PassUser, error) {
	row := adb.db.QueryRow("SELECT * FROM users WHERE nick=$1", nick)
	passuser := new(model.PassUser)
	if err := row.Scan(&passuser.Id, &passuser.Nick, &passuser.PassHash); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, errors.New("no user in DB")
		case err != nil:
			return nil, err
		}
	}
	return passuser, nil
}


func (adb AppDB) DelPassUser(pu *model.PassUser) error {
	_, err := adb.db.Exec("DELETE FROM users WHERE id=$1", pu.Id)
	return err
}