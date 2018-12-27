package evo

import (
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/world"
)

// CollisionHandler handles collision detection.
type CollisionHandler interface {
	DetectCollisions(creatures []*entity.Creature)
}

// EntityHandler handles the entitiy pouplation.
type EntityHandler interface {
	// InitPopulation initializes the population with a given count.
	InitPopulation(count int) []*entity.Creature

	// UpdatePopulation updates the entitiy population.
	UpdatePopulation(creatures []*entity.Creature) []*entity.Creature
}

// Simulation holds all simulation data.
type Simulation struct {
	width          int
	height         int
	ticksPerSecond int
	running        bool
	pause          bool
	creatures      []*entity.Creature

	stats *Stats

	collisionHandler CollisionHandler
	entityHandler    EntityHandler

	// Stores the subscriptions to a entities changed event.
	entitiesChangedSubscriptions map[uuid.UUID]EntitiesChangedFn

	m *sync.Mutex
}

// NewSimulation creates a new simulation.
func NewSimulation(width, height, ticksPerSecond int) *Simulation {
	s := &Simulation{
		width:          width,
		height:         height,
		ticksPerSecond: ticksPerSecond,
		creatures:      nil,

		collisionHandler:             world.NewSimpleCollisionHandler(width, height),
		entityHandler:                entity.NewHandler(width, height),
		entitiesChangedSubscriptions: make(map[uuid.UUID]EntitiesChangedFn),

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

	s.creatures = s.entityHandler.InitPopulation(1000)
}

// Start starts the simulation.
func (s *Simulation) Start() error {
	s.running = true
	for s.running {
		s.m.Lock()
		start := time.Now()
		if !s.pause {
			s.collisionHandler.DetectCollisions(s.creatures)
			s.creatures = s.entityHandler.UpdatePopulation(s.creatures)
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

// Size returns the size of the simulation
func (s *Simulation) Size() (int, int, error) {
	return s.width, s.height, nil
}

// Creatures returns all creatures.
func (s *Simulation) Creatures() ([]*entity.Creature, error) {
	return s.creatures, nil
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

// SubscribeEntitiesChanged implements the entities changed subscription.
func (s *Simulation) SubscribeEntitiesChanged(fn EntitiesChangedFn) uuid.UUID {
	u := uuid.New()
	s.entitiesChangedSubscriptions[u] = fn

	return u
}

// UnsubscribeEntitiesChanged implments the unsubscription for the entities
// changed event.
func (s *Simulation) UnsubscribeEntitiesChanged(id uuid.UUID) {
	delete(s.entitiesChangedSubscriptions, id)
}

func (s *Simulation) handleSubscriptions() {
	for _, fn := range s.entitiesChangedSubscriptions {
		fn(s.creatures)
	}
}

func (s *Simulation) updateStats() {
	s.stats.Running = time.Since(s.stats.start) / (time.Millisecond * 1000)
	s.stats.Population = len(s.creatures)
}
