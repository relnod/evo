package evo

import (
	"log"
	"math/rand"
	"time"

	"github.com/relnod/evo/system"
	"github.com/relnod/evo/world"
	uuid "github.com/satori/go.uuid"
)

// Simulation holds all simulation data.
type Simulation struct {
	ticksPerSecond int

	world           *world.World
	collisionSystem *system.Collision
	entitySystem    *system.Entity

	// events  []Event
	streams map[uuid.UUID]Stream
}

// NewSimulation creates a new simulation.
func NewSimulation() *Simulation {
	seed := time.Now().Unix()
	rand.Seed(seed)
	log.Println("Seed: ", seed)

	world := world.NewWorld(
		500.0, // @todo
		500.0, // @todo
		world.EdgeModeLoop,
		world.StartModeRandom,
	)

	collisionSystem := system.NewCollision(world)
	entitySystem := system.NewEntity(world)

	entitySystem.Init()

	return &Simulation{
		ticksPerSecond:  60,
		world:           world,
		collisionSystem: collisionSystem,
		entitySystem:    entitySystem,
		streams:         make(map[uuid.UUID]Stream),
	}
}

// Start starts the simulation in a new gorotinue.
func (s *Simulation) Start() {
	for {
		s.world.UpdateCells()
		s.collisionSystem.Update()
		s.entitySystem.Update()
		s.handleStreams()
		// s.handleEvents()

		time.Sleep(time.Second / time.Duration(s.ticksPerSecond))
	}
}

func (s *Simulation) GetWorld() *world.World {
	return s.world
}

func (s *Simulation) RegisterStream(stream Stream) uuid.UUID {
	u := uuid.NewV4()
	s.streams[u] = stream

	return u
}

func (s *Simulation) UnRegisterStream(id uuid.UUID) {
	delete(s.streams, id)
}

func (s *Simulation) handleStreams() {
	for _, stream := range s.streams {
		stream(s.world)
	}
}

// // AddEvent pushes an event to the internal event queue.
// func (s *Simulation) AddEvent(event Event) {
// 	s.events = append(s.events, event)
// }

// func (s *Simulation) handleEvents() {
// 	for _, event := range s.events {
// 		switch event.Type() {
// 		case EventRecieveWorld:
// 			event.SetData(*s.world)
// 		}
// 	}

// 	s.events = s.events[:0]
// }
