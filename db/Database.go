package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"errors"
)

var ErrNotFound = errors.New("db: entry not found")
var ErrInternal = errors.New("db: interal db error")

type AppDB struct {
	db *sql.DB
}

type CacheDB struct {
	adb AppDB
	//TODO: add cache (user, transactions, subtransactions) pointer in the future
}

func InitDB(cfg string) (CacheDB, error) {
	var cdb CacheDB
	var err error

	/* parse cfg string, split to dbCfg and CacheCfg */
	dbCfg := cfg //temp
	cdb.adb, err = InitAppDB(dbCfg)
	if err != nil {
		return cdb, err
	}

	/*init cache*/

	return cdb, nil
}

func InitAppDB(cfg string) (AppDB, error) {
	db, err := sql.Open("postgres", cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	createTables(db)
	return AppDB{db: db}, nil
}

func createTables(db *sql.DB) {

	createUsers := `

	CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	email varchar(255) UNIQUE NOT NULL,
	nick varchar(255) UNIQUE NOT NULL,
	passhash varchar(64) NOT NULL
	);`

	createTransactions := `

	CREATE TABLE IF NOT EXISTS transactions (
	tr_id serial PRIMARY KEY,
	date timestamp(0) without time zone NOT NULL,
	source integer NOT NULL,
	sum integer NOT NULL,
	matter text NOT NULL,
	comment text NOT NULL
	);`

	createSubtransactions := `

	CREATE TABLE IF NOT EXISTS subtransactions (
	tr_id integer NOT NULL,
	source integer NOT NULL,
	target integer NOT NULL,
	sum integer NOT NULL,
	proportion integer NOT NULL
	);`

	_, err := db.Exec(createUsers)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTransactions)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createSubtransactions)
	if err != nil {
		log.Fatal(err)
	}
}
