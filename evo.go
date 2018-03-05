package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/relnod/evo/collision"
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

	s := system.NewSystem(width, height, 100, 200)
	app := &App{
		system:          s,
		renderSystem:    system.NewRender(s, canvas),
		collisionSystem: system.NewCollision(s, 36),
		entitySystem:    system.NewEntity(s),
	}

	app.collisionSystem.SetCreatureCreatureCB(func(e1 *entity.Creature, e2 *entity.Creature) {
		// e1.Die()
		// e2.Die()
	})

	app.collisionSystem.SetCreatureEyeFoodCB(func(e *entity.Creature, f *entity.Food) {
		// fmt.Println("sawfood = true")
		e.SawFood = true
	})

	app.collisionSystem.SetCreatureFoodCB(func(e *entity.Creature, f *entity.Food) {
		e.Saturation += 1.0
		app.entitySystem.ResetFood(f)
	})

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

	return app
}

func (app *App) Run() {
	for {
		app.update()
		time.Sleep(time.Second / 60)

		// time.Sleep(time.Second)
	}
}

func (app *App) update() {
	app.collisionSystem.Update()
	app.entitySystem.Update()
	app.renderSystem.Update()
}
