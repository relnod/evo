package main

import (
	"flag"
	"fmt"

	"github.com/relnod/evo/api"
	"github.com/relnod/evo/graphics"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	fmt.Println(*addr)

	app := graphics.NewClient(api.NewWebsocketClient(*addr))
	app.Init()
	app.Start()
}
