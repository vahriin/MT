package main

import (
	"github.com/vahriin/MT/daemon"
	"github.com/vahriin/MT/config"
)

func main() {
	cfg := config.ReadConfig()
	daemon.Run(cfg)
}


