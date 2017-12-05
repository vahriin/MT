package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/vahriin/MT/config"
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
	return AppDB{db: db}, nil
}


