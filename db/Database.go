package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"errors"
	"github.com/vahriin/MT/config"
)

var ErrNotFound = errors.New("db: entry not found")
var ErrInternal = errors.New("db: internal db error")

type AppDB struct {
	db *sql.DB
}

type CacheDB struct {
	adb AppDB
	//TODO: add cache (user, transactions, subtransactions) pointer in the future
}



func InitDB(cfg *config.AppDbConfig) (CacheDB, error) {
	var cdb CacheDB
	var err error

	var dbConfig string =
		"dbname=" + cfg.Db.Name + " " +
		"user=" + cfg.Db.User + " " +
		"password='" + cfg.Db.Password + "'" + " " +
		"host=" + cfg.Db.Host + " " +
		"port=" + cfg.Db.Port + " " + //TODO: add connect_timeout
		"sslmode=" + cfg.Db.Sslmode

	//TODO: add support of these verification types
	switch cfg.Db.Sslmode {
	case "require":
		fallthrough
	case "verify-ca":
		fallthrough
	case "verify-full":
		panic("This verification type is not supported")
	}

	cdb.adb, err = initAppDB(dbConfig)
	if err != nil {
		return cdb, err
	}

	/*init cache*/

	return cdb, nil
}

func initAppDB(cfg string) (AppDB, error) {
	db, err := sql.Open("postgres", cfg)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	createTables(db)
	return AppDB{db: db}, nil
}

func createTables(db *sql.DB) {
	//TODO: check exist and change these requests
	createUsers := `

	CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	email varchar(255) UNIQUE NOT NULL,
	nick varchar(255) NOT NULL,
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
