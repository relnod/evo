package main

import (
	"flag"

	"github.com/relnod/evo/api/client"
	"github.com/relnod/evo/api/server"
)

var listenAddr = flag.String("listenAddr", "localhost:8080", "listen address")
var serveAddr = flag.String("serveAddr", ":8081", "serve address")
var debug = flag.Bool("debug", false, "enable debugging")

func main() {
	flag.Parse()

	server := server.New(client.New(*listenAddr), *serveAddr, *debug)
	server.Start()
}
