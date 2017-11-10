package server

import (
	"github.com/vahriin/MT/config"
	"net/http"
)

func InitServer(config *config.ServerConfig) http.Server {
	http.Handle("/", http.FileServer(http.Dir(config.HtmlRoot)))
	return createServer(config)
}

func createServer(config *config.ServerConfig) http.Server {
	return http.Server {
		Addr: config.Address + ":" + config.Port}
}
