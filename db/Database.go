package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type AppDB struct {
	db *sql.DB
	//TODO: add cache (user, transactions, subtransactions) in the future
}

func InitDB(cfg string) (AppDB, error) {
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
	sum integer NOT NULL
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
