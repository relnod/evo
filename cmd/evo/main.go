package main

import (
	"github.com/relnod/evo"
	"github.com/relnod/evo/graphics"
)

func main() {
	app := graphics.NewClient(evo.NewSimulation())
	app.Init()
	app.Start()
}
