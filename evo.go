package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/relnod/evo/config"
	"github.com/relnod/evo/platform"
	"github.com/relnod/evo/system"
	"github.com/relnod/evo/world"
)

func main() {
	app := NewApp()
	app.Run()
}

type App struct {
	ticksPerSecond int

	window          *platform.Window
	world           *world.World
	renderSystem    *system.Render
	collisionSystem *system.Collision
	entitySystem    *system.Entity
}

func NewApp() *App {
	seed := time.Now().Unix()
	log.Println("Seed: ", seed)
	rand.Seed(seed)

	window := platform.NewWindow()
	world := world.NewWorld(
		float32(window.Width),
		float32(window.Height),
		world.EdgeModeLoop,
		world.StartModeFixed,
	)

	app := &App{
		ticksPerSecond:  60,
		window:          window,
		world:           world,
		renderSystem:    system.NewRender(world, window),
		collisionSystem: system.NewCollision(world),
		entitySystem:    system.NewEntity(world),
	}

	app.renderSystem.Init()
	app.entitySystem.Init()

	window.AddKeyListener(func(e *platform.KeyEvent) {
		switch e.KeyCode {
		case 38: // UP
			config.WorldSpeed += 5
		case 40: // UP
			config.WorldSpeed -= 5
		}
	})

	window.AddMouseListener(func(e *platform.MouseEvent) {
		creature := app.collisionSystem.FindCreature(&e.Pos)
		if creature != nil {
			fmt.Printf("Creature:\n")
			fmt.Printf("Generation: %d\n", creature.Consts.Generation)
			fmt.Printf("Radius: %f\n", creature.Radius)
			fmt.Printf("\n")
		}
	})

	return app
}

func (app *App) Run() {
	for {
		app.update()
		time.Sleep(time.Second / time.Duration(app.ticksPerSecond))
	}
}

func (app *App) update() {
	app.world.UpdateCells()
	app.collisionSystem.Update()
	app.entitySystem.Update()
	app.renderSystem.Update()
}
