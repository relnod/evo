package evo

import (
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/relnod/evo/pkg/system"
	"github.com/relnod/evo/pkg/world"
)

// Simulation holds all simulation data.
type Simulation struct {
	ticksPerSecond int

	world           *world.World
	collisionSystem *system.Collision
	entitySystem    *system.Entity

	worldSubscriptions map[uuid.UUID]WorldStream
}

// NewSimulation creates a new simulation.
func NewSimulation() *Simulation {
	seed := time.Now().Unix()
	rand.Seed(seed)
	log.Println("Seed: ", seed)

	world := world.NewWorld(
		1000.0, // @todo
		1000.0, // @todo
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

		worldSubscriptions: make(map[uuid.UUID]WorldStream),
	}
}

// Start starts the simulation in a new gorotinue.
func (s *Simulation) Start() error {
	for {
		s.world.UpdateCells()
		s.collisionSystem.Update()
		s.entitySystem.Update()
		s.pushStreams()
		// s.handleStreams()
		// s.handleEvents()

		time.Sleep(time.Second / time.Duration(s.ticksPerSecond))
	}
}

func (s *Simulation) GetWorld() (*world.World, error) {
	return s.world, nil
}

func (s *Simulation) SubscribeWorld(stream WorldStream) uuid.UUID {
	u := uuid.New()
	s.worldSubscriptions[u] = stream

	return u
}

func (s *Simulation) UnsubscribeWorld(id uuid.UUID) {
	delete(s.worldSubscriptions, id)
}

func (s *Simulation) pushStreams() {
	for _, stream := range s.worldSubscriptions {
		stream(s.world)
	}
}
