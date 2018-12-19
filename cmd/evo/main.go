package main

import (
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/graphics"
)

func main() {
	app := graphics.NewClient(evo.NewSimulation())
	app.Init()
	app.Start()
}
