package main

import (
	"github.com/relnod/evo"
)

func main() {
	app := evo.NewRenderClient(evo.NewSimulation())
	app.Init()
	app.Start()
}
