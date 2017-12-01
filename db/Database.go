package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/vahriin/MT/config"
	"log"
)

var ErrNotFound = errors.New("db: entry not found")
var ErrInternal = errors.New("db: internal db error")
var ErrForbidden = errors.New("db: this operation are forbidden")

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

	var dbConfig = "dbname=" + cfg.Db.Name + " " +
		"user='" + cfg.Db.User + "' " +
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
	createUser := `

	CREATE TABLE IF NOT EXISTS app_user (
	id serial PRIMARY KEY,
	google_id varchar(255) UNIQUE NOT NULL,
	nick varchar(255) NOT NULL
	);`

	createGroup := `

	CREATE TABLE IF NOT EXISTS app_group (
	id serial PRIMARY KEY,
	name varchar(255) UNIQUE NOT NULL,
	creator_id integer NOT NULL
	);`

	createUserGroup := `

	CREATE TABLE IF NOT EXISTS app_user_group (
	user_id integer NOT NULL,
	gr_id integer NOT NULL
	);`

	createTransaction := `

	CREATE TABLE IF NOT EXISTS app_transaction (
	tr_id serial PRIMARY KEY,
	gr_id integer NOT NULL,
	date timestamp(0) without time zone NOT NULL,
	source integer NOT NULL,
	sum integer NOT NULL,
	matter text NOT NULL,
	comment text NOT NULL
	);`

	createSubtransaction := `

	CREATE TABLE IF NOT EXISTS app_subtransaction (
	tr_id integer NOT NULL,
	source integer NOT NULL,
	target integer NOT NULL,
	sum integer NOT NULL,
	proportion integer NOT NULL
	);`

	var err error
	if _, err = db.Exec(createUser); err != nil {
		log.Fatal("CreateUser returned this message: " + err.Error())
	}

	if _, err = db.Exec(createGroup); err != nil {
		log.Fatal("CreateGroup returned this message: " + err.Error())
	}

	if _, err = db.Exec(createUserGroup); err != nil {
		log.Fatal("CreateUserGroup returned this message: " + err.Error())
	}

	if _, err = db.Exec(createTransaction); err != nil {
		log.Fatal("CreateTransaction returned this message: " + err.Error())
	}

	if _, err = db.Exec(createSubtransaction); err != nil {
		log.Fatal("CreateSubtransaction returned this message: " + err.Error())
	}
}
