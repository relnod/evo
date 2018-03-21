package main

import (
	"flag"

	"github.com/relnod/evo"
	"github.com/relnod/evo/api"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	app := evo.NewRenderClient(api.NewWebsocketClient(*addr))
	app.Init()
	app.Start()
}
