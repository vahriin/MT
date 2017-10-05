package daemon

import (
	"github.com/vahriin/MT/config"
	"os"
	"log"
	"github.com/vahriin/MT/db"
	"net/http"
	"github.com/vahriin/MT/api"
)

/* this function use for test run of server */
/* delete later */
func TempRun() {
	var configstr string = "user=vahriin dbname=MT_DB sslmode=disable"

	server := http.Server{Addr: "127.0.0.1:4000"}
	appDb, _ := db.InitDB(configstr)
	http.Handle("/transactions", api.TransactionsHandler(&appDb))
	http.Handle("/transaction", api.TransactionHandler(&appDb))
	http.Handle("/user", api.UserHandler(&appDb))
	server.ListenAndServe()
}

func Run(config *config.AppConfig) {
	/* log init */

	file, err := os.OpenFile(config.System.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	log.SetOutput(file)

	/* end of log init */

	/*-----------------------------------------------------*/

	/* db init */

	//AppDatabase, err := db.InitDB(config.DataModule)
	if err != nil {
		panic(err)
	}

	/* end of db init */

	/*-----------------------------------------------------*/

	/* handler init */

	//http.Handle("/transactions", api.TransactionsHandler(&AppDatabase))
	//http.Handle("/transaction", api.TransactionHandler(&AppDatabase))
	//http.Handle("/user", api.UserHandler(&AppDatabase))


	/* end of hanler init */

	/*-----------------------------------------------------*/

	/* server init */

	Server := initServer(config.Server)
	Server.ListenAndServe()

	/* end of server init */

	/* End of app initialization */
}
