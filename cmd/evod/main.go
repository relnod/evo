package main

import (
	"flag"

	"github.com/relnod/evo/api/server"
	"github.com/relnod/evo/pkg/evo"
)

var addr = flag.String("addr", ":8080", "address")

func main() {
	flag.Parse()

	server := server.New(evo.NewSimulation(), *addr)
	server.Start()
}
