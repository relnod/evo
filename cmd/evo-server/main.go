package main

import (
	"flag"

	"github.com/relnod/evo"
	"github.com/relnod/evo/api"
)

var addr = flag.String("addr", ":8080", "address")

func main() {
	flag.Parse()

	app := api.NewWebsocketServer(evo.NewSimulation(), *addr)
	app.Init()
	app.Start()
}
