package daemon

import (
	"github.com/vahriin/MT/api"
	"github.com/vahriin/MT/config"
	"github.com/vahriin/MT/db"
	"github.com/vahriin/MT/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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

	AppDatabase, err := db.InitDB(&config.AppDbConfig)
	if err != nil {
		panic(err)
	}

	/* end of db init */

	/*-----------------------------------------------------*/

	/* handler init */

	http.Handle("/transactions", api.TransactionsHandler(&AppDatabase))
	http.Handle("/transaction-id", api.TransactionIdHandler(&AppDatabase))
	http.Handle("/user", api.UserHandler(&AppDatabase))

	/* end of handler init */

	/*-----------------------------------------------------*/

	/* server init */

	Server := server.InitServer(&config.Server)
	go Server.ListenAndServeTLS(config.Server.CertFile, config.Server.KeyFile)
	//go Server.ListenAndServe()

	/* end of server init */

	waitForSignal()

	/* End of app initialization */
}

/*No test*/
func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
