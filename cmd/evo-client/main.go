package main

import (
	"flag"

	"github.com/relnod/evo"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	app := evo.NewRenderClient(evo.NewWebsocketClient(*addr))
	app.Init()
	app.Start()
}
