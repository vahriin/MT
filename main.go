package main

import (
	"net/http"
	"github.com/vahriin/MT/api"
	"github.com/vahriin/MT/db"
)

func main() {
	var config string = "user=vahriin dbname=MT_DB sslmode=disable"

	server := http.Server{Addr: "127.0.0.1:4000"}
	appDb, _ := db.InitDB(config)
	http.Handle("/transactions", api.TransactionsHandler(&appDb))
	http.Handle("/transaction", api.TransactionHandler(&appDb))
	server.ListenAndServe()
}
