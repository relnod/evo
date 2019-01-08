package main

import (
	"flag"

	"github.com/relnod/evo/api/server"
	"github.com/relnod/evo/pkg/evo"
)

var addr = flag.String("addr", ":8080", "address")
var debug = flag.Bool("debug", false, "enable debugging")

func main() {
	flag.Parse()

	server := server.New(evo.NewSimulationFromSeed(2000, 2000, 1000, 2), *addr, *debug)
	server.Start()
}
