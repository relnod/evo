package evo

import (
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/relnod/evo/pkg/system"
	"github.com/relnod/evo/pkg/world"
)

// Simulation holds all simulation data.
type Simulation struct {
	width          int
	height         int
	ticksPerSecond int
	running        bool
	pause          bool

	stats           *Stats
	world           *world.World
	collisionSystem *system.Collision
	entitySystem    *system.Entity

	// Stores the subscriptions to a world change.
	worldSubscriptions map[uuid.UUID]WorldFn

	m *sync.Mutex
}

// NewSimulation creates a new simulation.
func NewSimulation(width, height, ticksPerSecond int) *Simulation {
	s := &Simulation{
		width:          width,
		height:         height,
		ticksPerSecond: ticksPerSecond,

		worldSubscriptions: make(map[uuid.UUID]WorldFn),

		m: &sync.Mutex{},
	}
	s.init()

	return s
}

func (s *Simulation) init() {
	s.running = true
	s.pause = false

	s.stats = &Stats{
		start: time.Now(),
		Seed:  time.Now().Unix(),
	}

	rand.Seed(s.stats.Seed)

	s.world = world.NewWorld(s.width, s.height)

	s.collisionSystem = system.NewCollision(s.world)
	s.entitySystem = system.NewEntity(s.world)

	s.entitySystem.Init()
}

// Start starts the simulation.
func (s *Simulation) Start() error {
	s.running = true
	for s.running {
		s.m.Lock()
		start := time.Now()
		if !s.pause {

			s.world.Update()
			s.collisionSystem.Update()
			s.entitySystem.Update()
		}
		s.handleSubscriptions()
		s.m.Unlock()

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

// Pause pauses the simulation.
func (s *Simulation) Pause() error {
	s.pause = true
	return nil
}

// Resume resumes the simulation.
func (s *Simulation) Resume() error {
	s.pause = false
	return nil
}

// PauseResume toggles pause/resume.
func (s *Simulation) PauseResume() error {
	s.m.Lock()
	if s.pause {
		s.Resume()
	} else {
		s.Pause()
	}
	s.m.Unlock()
	return nil
}

// Restart restarts the simulation
func (s *Simulation) Restart() error {
	s.m.Lock()
	s.init()
	s.m.Unlock()
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

// Ticks returns the ticks per second.
func (s *Simulation) Ticks() (int, error) {
	return s.ticksPerSecond, nil
}

// SetTicks sets the ticks per second.
func (s *Simulation) SetTicks(ticks int) error {
	if ticks <= 0 {
		ticks = 1
	}
	s.ticksPerSecond = ticks
	return nil
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
