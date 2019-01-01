package api

import (
	"github.com/google/uuid"
	"github.com/relnod/evo/pkg/entity"
)

// EntitiesChangedFn defnies a callback function for entities.
type EntitiesChangedFn func([]*entity.Creature)

// SubscriptionHandler handles event subscriptions.
type SubscriptionHandler struct {
	entitiesChangedSubscriptions map[uuid.UUID]EntitiesChangedFn
}

// NewSubscriptionHandler returns a new event subscriber.
func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{
		entitiesChangedSubscriptions: make(map[uuid.UUID]EntitiesChangedFn),
	}
}

// SubscribeEntitiesChanged implements the entities changed subscription.
func (s *SubscriptionHandler) SubscribeEntitiesChanged(fn EntitiesChangedFn) uuid.UUID {
	u := uuid.New()
	s.entitiesChangedSubscriptions[u] = fn

	return u
}

// UnsubscribeEntitiesChanged implments the unsubscription for the entities
// changed event.
func (s *SubscriptionHandler) UnsubscribeEntitiesChanged(id uuid.UUID) {
	delete(s.entitiesChangedSubscriptions, id)
}

// Update triggers all event subscriptions.
func (s *SubscriptionHandler) Update(creatures []*entity.Creature) {
	for _, fn := range s.entitiesChangedSubscriptions {
		fn(creatures)
	}
}
