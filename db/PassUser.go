package db

import (
	"database/sql"
	"errors"
	"github.com/vahriin/MT/model"
	"log"
)

/* change this for OAuth */
func (adb AppDB) AddPassUser(pu *model.PassUser) error {
	row := adb.db.QueryRow("SELECT id FROM users WHERE nick=$1", pu.Nick)
	var err error
	if row.Scan(&pu.Id) == sql.ErrNoRows {
		_, err = adb.db.Exec("INSERT INTO users(nick, email, passhash) VALUES($1, $2, $3)",
			pu.Nick, pu.Email, pu.PassHash)
	} else {
		err = errors.New("User " + pu.Nick + " already exists")
	}

	return err
}

func (adb AppDB) GetPassUserByEmail(email string) (*model.PassUser, error) {
	row := adb.db.QueryRow("SELECT id, email, nick, passhash FROM users WHERE nick=$1", email)
	passuser := new(model.PassUser)
	if err := row.Scan(&passuser.Id, &passuser.Email, &passuser.Nick, &passuser.PassHash); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrNotFound
		case err != nil:
			log.Println("getPassUserByEmail returned this message: " + err.Error())
			return nil, ErrInternal
		}
	}
	return passuser, nil
}

func (adb AppDB) DeletePassUser(pu *model.PassUser) error {
	result, err := adb.db.Exec("DELETE FROM users WHERE id=$1", pu.Id)
	if err != nil {
		log.Println("DeletePassUser returned this message: " + err.Error())
		return ErrInternal
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		return ErrNotFound
	}
	return err
}
