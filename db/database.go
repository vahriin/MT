package db

import (
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/vahriin/MT/model"
	"errors"
	"log"
)


type AppDB struct {
	db *sql.DB
}


func InitDB(cfg string) (AppDB, error) {
	db, err := sql.Open("postgres", cfg)
	if err != nil {
		log.Fatal(err)
	} else {
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
	}
	return AppDB{db}, nil
}


func (adb AppDB) CreateTables() {
	createUsers := `

	CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	nick varchar(255) NOT NULL,
	passhash varchar(255) NOT NULL
	);`


	createTransactions := `

	CREATE TABLE IF NOT EXISTS transactions (
	tr_id serial PRIMARY KEY,
	date timestamp without time zone NOT NULL,
	source integer NOT NULL,
	targets integer[] NOT NULL,
	sum integer NOT NULL,
	matter text NOT NULL,
	comment text NOT NULL
	);`

	createSubtransactions := `

	CREATE TABLE IF NOT EXISTS subtransactions (
	tr_id integer NOT NULL,
	source integer NOT NULL,
	target integer NOT NULL,
	sum integer NOT NULL
	);`

	_, err := adb.db.Exec(createUsers)
	if err != nil {
		log.Fatal(err)
	}

	_, err = adb.db.Exec(createTransactions)
	if err != nil {
		log.Fatal(err)
	}

	_, err = adb.db.Exec(createSubtransactions)
	if err != nil {
		log.Fatal(err)
	}
}


func (adb AppDB) GetPassUser(username string) (*model.PassUser, error) {
	row := adb.db.QueryRow("SELECT * FROM clients WHERE nick=$1", username)
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



