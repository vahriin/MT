package daemon

import (
	"github.com/vahriin/MT/config"
	"net/http"
)

func initServer(config config.ServerConfig) http.Server {
	/* parse config and create server */

	return http.Server{Addr: ":4000"}
}
