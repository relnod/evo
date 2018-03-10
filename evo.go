package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/config"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/system"
)

func main() {
	seed := time.Now().Unix()
	log.Println("Seed: ", seed)
	rand.Seed(seed)

	app := NewApp()
	app.Run()
}

type App struct {
	ticksPerSecond  int
	system          *system.System
	renderSystem    *system.Render
	collisionSystem *system.Collision
	entitySystem    *system.Entity
}

func NewApp() *App {
	document := js.Global.Get("document")

	// @todo: instead wait till dom is ready
	time.Sleep(time.Millisecond)

	window := js.Global.Get("window")
	w, err := strconv.Atoi(window.Get("innerWidth").String())
	if err != nil {
		// @todo
	}
	h, err := strconv.Atoi(window.Get("innerHeight").String())
	if err != nil {
		// @todo
	}

	width := float32(w)
	height := float32(h)

	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", width)
	canvas.Set("height", height)

	body := document.Get("body")
	body.Get("style").Set("margin", 0)
	body.Call("appendChild", canvas)

	s := system.NewSystem(width, height, 1)
	app := &App{
		ticksPerSecond:  60,
		system:          s,
		renderSystem:    system.NewRender(s, canvas),
		collisionSystem: system.NewCollision(s, 36),
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

	js.Global.Call("addEventListener", "keyup", func(event *js.Object) {
		keycode := event.Get("keyCode").Int()
		if keycode == 38 { // UP
			config.WorldSpeed += 5
		}
		if keycode == 40 { // DOWN
			config.WorldSpeed -= 5
		}
	}, false)

	return app
}

func (app *App) Run() {
	for {
		app.update()
		time.Sleep(time.Second / time.Duration(app.ticksPerSecond))

		// time.Sleep(time.Second)
	}
}

func (app *App) update() {
	app.collisionSystem.Update()
	app.entitySystem.Update()
	app.renderSystem.Update()
}
