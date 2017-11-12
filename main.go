package main

import (
	"github.com/vahriin/MT/config"
	"github.com/vahriin/MT/daemon"
)

func main() {
	cfg := config.ReadConfig()
	daemon.Run(cfg)
}
