package main

import (
	"flag"

	"github.com/relnod/evo/api/websocket"
	"github.com/relnod/evo/pkg/evo"
)

var addr = flag.String("addr", ":8080", "address")

func main() {
	flag.Parse()

	server := websocket.NewServer(evo.NewSimulation(), *addr)
	server.Start()
}
