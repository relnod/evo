package main

import (
	"flag"

	"github.com/relnod/evo/api/websocket"
	"github.com/relnod/evo/graphics"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	app := graphics.NewClient(websocket.NewClient(*addr))
	app.Init()
	app.Start()
}
