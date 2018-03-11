package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/config"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/platform"
	"github.com/relnod/evo/system"
)

func main() {
	app := NewApp()
	app.Run()
}

type App struct {
	ticksPerSecond int

	window          *platform.Window
	system          *system.System
	renderSystem    *system.Render
	collisionSystem *system.Collision
	entitySystem    *system.Entity
}

func NewApp() *App {
	seed := time.Now().Unix()
	log.Println("Seed: ", seed)
	rand.Seed(seed)

	w := platform.NewWindow()
	s := system.NewSystem(float32(w.Width), float32(w.Height), 1)

	app := &App{
		ticksPerSecond:  60,
		window:          w,
		system:          s,
		renderSystem:    system.NewRender(s, w),
		collisionSystem: system.NewCollision(s, 4),
		entitySystem:    system.NewEntity(s, system.ModeFixed),
	}

	app.collisionSystem.SetCreatureBorderCB(func(e *entity.Creature, border int) {
		switch border {
		case collision.LEFT:
			e.Pos.X += float32(app.system.Width)
		case collision.RIGHT:
			e.Pos.X -= float32(app.system.Width)
		case collision.TOP:
			e.Pos.Y += float32(app.system.Height)
		case collision.BOT:
			e.Pos.Y -= float32(app.system.Height)
		}
	})

	app.renderSystem.Init()
	app.entitySystem.Init()

	w.AddKeyListener(func(e *platform.KeyEvent) {
		switch e.KeyCode {
		case 38: // UP
			config.WorldSpeed += 5
		case 40: // UP
			config.WorldSpeed -= 5
		}
	})

	w.AddMouseListener(func(e *platform.MouseEvent) {
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
	app.collisionSystem.Update()
	app.entitySystem.Update()
	app.renderSystem.Update()
}
