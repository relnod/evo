package main

import (
	"flag"

	"github.com/relnod/evo/api/websocket"
	"github.com/relnod/evo/pkg/graphics"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	client := graphics.NewClient(websocket.NewClient(*addr))
	client.Init()
	client.Start()
}
