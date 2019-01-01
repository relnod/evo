package evo

import (
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/relnod/evo/api"
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/stats"
	"github.com/relnod/evo/pkg/world"
)

// CollisionHandler handles collision detection.
type CollisionHandler interface {
	DetectCollisions(creatures []*entity.Creature)
}

// EntityHandler handles the entitiy poplation.
type EntityHandler interface {
	// InitPopulation initializes the population with a given count.
	InitPopulation(count int) []*entity.Creature

	// UpdatePopulation updates the entitiy population.
	UpdatePopulation(creatures []*entity.Creature) []*entity.Creature

	AnimalStats() *entity.Stats
	PlantStats() *entity.Stats
}

// StatsCollector collects stats.
type StatsCollector interface {
	Update(tick int, creatures []*entity.Creature)
	Stats() *stats.Stats
}

// SubscriptionHandler defines an evnet subscriber.
type SubscriptionHandler interface {
	// SubscribeEntitiesChanged subscribes to changes of entities.
	// Each time the entities get updated, the provided function gets called.
	// The returned unique id can be used to unsubscribe later.
	SubscribeEntitiesChanged(fn api.EntitiesChangedFn) uuid.UUID

	// UnsubscribeWorldChange ends a subscription to the world change.
	UnsubscribeEntitiesChanged(id uuid.UUID)

	// Update triggers the subscriptions.
	Update(creatures []*entity.Creature)
}

// Simulation holds all simulation data.
type Simulation struct {
	seed   int64
	width  int
	height int

	creatures []*entity.Creature

	ticker              *ticker
	collisionHandler    CollisionHandler
	entityHandler       EntityHandler
	subscriptionHandler SubscriptionHandler
	statsCollector      StatsCollector
}

// NewSimulation creates a new simulation.
func NewSimulation(width, height, ticksPerSecond int) *Simulation {
	seed := time.Now().Unix()
	s := &Simulation{
		seed:   seed,
		width:  width,
		height: height,

		creatures: nil,

		collisionHandler:    world.NewSimpleCollisionHandler(width, height),
		entityHandler:       entity.NewHandler(width, height),
		statsCollector:      stats.NewIntervalCollector(seed, 5),
		subscriptionHandler: api.NewSubscriptionHandler(),
	}
	s.ticker = newTicker(ticksPerSecond, func(tick int) error {
		s.collisionHandler.DetectCollisions(s.creatures)
		s.creatures = s.entityHandler.UpdatePopulation(s.creatures)
		return nil
	})
	s.ticker.SetAlwaysUpdate(func(tick int) error {
		s.subscriptionHandler.Update(s.creatures)
		s.statsCollector.Update(tick, s.creatures)
		return nil
	})
	s.init()

	return s
}

func (s *Simulation) init() {
	s.ticker.Resume()

	rand.Seed(s.seed)

	s.creatures = s.entityHandler.InitPopulation(1000)
}

// Start starts the simulation.
func (s *Simulation) Start() error {
	return s.ticker.Start()
}

// Stop stops the simulation.
func (s *Simulation) Stop() error {
	s.ticker.Stop()
	// TODO: cleanup
	// TODO: maybe wait until ticker is stoped
	return nil
}

// Pause pauses the simulation.
func (s *Simulation) Pause() error {
	s.ticker.Pause()
	return nil
}

// Resume resumes the simulation.
func (s *Simulation) Resume() error {
	s.ticker.Resume()
	return nil
}

// PauseResume toggles pause/resume.
func (s *Simulation) PauseResume() error {
	s.ticker.TogglePauseResume()
	return nil
}

// Restart restarts the simulation
func (s *Simulation) Restart() error {
	s.ticker.Lock()
	s.init()
	s.ticker.Unlock()
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
func (s *Simulation) Stats() (*stats.Stats, error) {
	return s.statsCollector.Stats(), nil
}

// Ticks returns the ticks per second.
func (s *Simulation) Ticks() (int, error) {
	return s.ticker.TicksPerSecond(), nil
}

// SetTicks sets the ticks per second.
func (s *Simulation) SetTicks(ticks int) error {
	if ticks <= 0 {
		ticks = 1
	}
	s.ticker.SetTicksPerSecond(ticks)
	return nil
}

// SubscribeEntitiesChanged implements the entities changed subscription.
func (s *Simulation) SubscribeEntitiesChanged(fn api.EntitiesChangedFn) uuid.UUID {
	return s.subscriptionHandler.SubscribeEntitiesChanged(fn)
}

// UnsubscribeEntitiesChanged implments the unsubscription for the entities
// changed event.
func (s *Simulation) UnsubscribeEntitiesChanged(id uuid.UUID) {
	s.subscriptionHandler.UnsubscribeEntitiesChanged(id)
}
