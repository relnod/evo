package evo

import (
	"math/rand"
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

	AnimalStats() *entity.Stats
	PlantStats() *entity.Stats
}

// Simulation holds all simulation data.
type Simulation struct {
	width                   int
	height                  int
	statsCollectionInterval int

	creatures                    []*entity.Creature
	stats                        *Stats
	entitiesChangedSubscriptions map[uuid.UUID]EntitiesChangedFn

	ticker           *ticker
	collisionHandler CollisionHandler
	entityHandler    EntityHandler
}

// NewSimulation creates a new simulation.
func NewSimulation(width, height, ticksPerSecond int) *Simulation {
	s := &Simulation{
		width:     width,
		height:    height,
		creatures: nil,

		statsCollectionInterval:      5,
		collisionHandler:             world.NewSimpleCollisionHandler(width, height),
		entityHandler:                entity.NewHandler(width, height),
		entitiesChangedSubscriptions: make(map[uuid.UUID]EntitiesChangedFn),
	}
	s.ticker = newTicker(ticksPerSecond, func(tick int) error {
		s.collisionHandler.DetectCollisions(s.creatures)
		s.creatures = s.entityHandler.UpdatePopulation(s.creatures)
		return nil
	})
	s.ticker.SetAlwaysUpdate(func(tick int) error {
		s.handleSubscriptions()
		if tick%s.statsCollectionInterval == 0 {
			s.collectTimeStats()
		}
		return nil
	})
	s.init()

	return s
}

func (s *Simulation) init() {
	s.ticker.Resume()

	s.stats = NewStats()
	rand.Seed(s.stats.Seed)

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
func (s *Simulation) Stats() (*Stats, error) {
	s.updateStats()
	return s.stats, nil
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

func (s *Simulation) collectTimeStats() {
	s.stats.OverTime.Add(s.currentTimeStat())
}

func (s *Simulation) currentTimeStat() *TimeStat {
	t := &TimeStat{
		Population: len(s.creatures),
		Animal:     &EntityTimeStat{},
		Plant:      &EntityTimeStat{},
	}

	for _, c := range s.creatures {
		if c.Brain == nil {
			t.Plant.Add(c)
		} else {
			t.Animal.Add(c)
		}

	}
	return t
}

func (s *Simulation) updateStats() {
	s.stats.Running = time.Since(s.stats.start) / (time.Millisecond * 1000)
	s.stats.Current = s.currentTimeStat()
	// s.stats.Animal = s.entityHandler.AnimalStats()
	// s.stats.Plant = s.entityHandler.PlantStats()
}
