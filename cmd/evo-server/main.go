package main

import (
	"flag"

	"github.com/relnod/evo"
	"github.com/relnod/evo/api/websocket"
)

var addr = flag.String("addr", ":8080", "address")

func main() {
	flag.Parse()

	app := websocket.NewServer(evo.NewSimulation(), *addr)
	app.Init()
	app.Start()
}
