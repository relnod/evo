package main

import (
	"flag"

	"github.com/relnod/evo/api/client"
	"github.com/relnod/evo/pkg/graphics"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	client := graphics.NewClient(client.New(*addr))
	client.Init()
	client.Start()
}
