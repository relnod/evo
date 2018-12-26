package evo

import (
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/relnod/evo/pkg/system"
	"github.com/relnod/evo/pkg/world"
)

// Simulation holds all simulation data.
type Simulation struct {
	ticksPerSecond int
	running        bool

	stats           *Stats
	world           *world.World
	collisionSystem *system.Collision
	entitySystem    *system.Entity

	// Stores the subscriptions to a world change.
	worldSubscriptions map[uuid.UUID]WorldFn
}

// NewSimulation creates a new simulation.
func NewSimulation() *Simulation {
	stats := &Stats{
		start: time.Now(),
		Seed:  time.Now().Unix(),
	}

	rand.Seed(stats.Seed)

	world := world.NewWorld(
		2000.0, // TODO: make it configurably
		2000.0, // TODO: make it configurably
	)

	collisionSystem := system.NewCollision(world)
	entitySystem := system.NewEntity(world)

	entitySystem.Init()

	return &Simulation{
		ticksPerSecond: 60,

		stats:           stats,
		world:           world,
		collisionSystem: collisionSystem,
		entitySystem:    entitySystem,

		worldSubscriptions: make(map[uuid.UUID]WorldFn),
	}
}

// Start starts the simulation.
func (s *Simulation) Start() error {
	s.running = true
	for s.running {
		start := time.Now()

		s.world.Update()
		s.collisionSystem.Update()
		s.entitySystem.Update()
		s.handleSubscriptions()

		time.Sleep(time.Second/time.Duration(s.ticksPerSecond) - time.Since(start))
	}
	return nil
}

// Stop stops the simulation.
func (s *Simulation) Stop() error {
	s.running = false
	// TODO: cleanup
	return nil
}

// World retruns the state of the current world.
func (s *Simulation) World() (*world.World, error) {
	return s.world, nil
}

// Stats returns the current statistics.
func (s *Simulation) Stats() (*Stats, error) {
	s.updateStats()
	return s.stats, nil
}

// SubscribeWorldChange implements the world change subscription.
func (s *Simulation) SubscribeWorldChange(stream WorldFn) uuid.UUID {
	u := uuid.New()
	s.worldSubscriptions[u] = stream

	return u
}

// UnsubscribeWorldChange implements the unsubscription of a world change.
func (s *Simulation) UnsubscribeWorldChange(id uuid.UUID) {
	delete(s.worldSubscriptions, id)
}

func (s *Simulation) handleSubscriptions() {
	for _, stream := range s.worldSubscriptions {
		stream(s.world)
	}
}

func (s *Simulation) updateStats() {
	s.stats.Running = time.Since(s.stats.start) / (time.Millisecond * 1000)
	s.stats.Population = len(s.world.Creatures)
}
