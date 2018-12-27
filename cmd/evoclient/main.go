package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/relnod/evo/api/client"
	"github.com/relnod/evo/pkg/graphics"
)

var addr = flag.String("addr", "localhost:8080", "address")
var help = flag.Bool("help", false, "Print help")

func main() {
	flag.Parse()

	if *help == true {
		fmt.Println((&graphics.Client{}).Usage())
		os.Exit(1)
	}

	client := graphics.NewClient(client.New(*addr))
	client.Init()
	client.Start()
}
