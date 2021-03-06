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

// EntityUpdater handles the entitiy poplation.
type EntityUpdater interface {
	// UpdatePopulation updates the entitiy population.
	UpdatePopulation(creatures []*entity.Creature) []*entity.Creature

	AnimalStats() *entity.DeathStats
	PlantStats() *entity.DeathStats
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
	seed              int64
	width             int
	height            int
	initialPopulation int

	creatures []*entity.Creature

	ticker              *Ticker
	entityUpdater       EntityUpdater
	collisionDetector   world.CollisionDetector
	subscriptionHandler SubscriptionHandler
	statsCollector      StatsCollector
}

// NewSimulation creates a new simulation.
func NewSimulation(width, height, population int) *Simulation {
	return NewSimulationFromSeed(width, height, population, time.Now().Unix())
}

// NewSimulationFromSeed creates a new simulation with a given seed. Therefore
// the siumulation should be 100% reproducable.
func NewSimulationFromSeed(width, height, population int, seed int64) *Simulation {
	entityUpdater := entity.NewPopulationUpdater()
	collisionDetector := world.NewSimpleCollisionDetector(width, height)
	statsCollector := stats.NewIntervalCollector(entityUpdater, seed, 5)

	s := &Simulation{
		seed:              seed,
		width:             width,
		height:            height,
		initialPopulation: population,

		creatures: nil,

		entityUpdater:       entityUpdater,
		collisionDetector:   collisionDetector,
		statsCollector:      statsCollector,
		subscriptionHandler: api.NewSubscriptionHandler(),
	}
	s.ticker = NewTicker(time.Second / 60)
	s.init()

	return s
}

func (s *Simulation) init() {
	s.ticker.Resume()

	rand.Seed(s.seed)

	s.creatures = entity.InitPopulation(s.initialPopulation, s.width, s.height)
}

// Update updates the simulation logic
func (s *Simulation) Update() {
	collisions := s.collisionDetector.DetectCollisions(s.creatures)
	world.ResolveAllCollisions(collisions)
	s.creatures = s.entityUpdater.UpdatePopulation(s.creatures)
}

// Start starts the simulation.
func (s *Simulation) Start() error {
	for tick := range s.ticker.C {
		s.Update()
		s.subscriptionHandler.Update(s.creatures)
		s.statsCollector.Update(tick, s.creatures)
	}
	return nil
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
	// TODO
	return int(s.ticker.Interval()), nil
}

// SetTicks sets the ticks per second.
func (s *Simulation) SetTicks(ticks int) error {
	// TODO
	s.ticker.SetInterval(time.Duration(ticks))
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
