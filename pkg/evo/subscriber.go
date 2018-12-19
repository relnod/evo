package evo

import (
	"github.com/google/uuid"
)

type Subscriber struct {
	worldSubscriptions map[uuid.UUID]WorldStream
}

func (s *Subscriber) SubscribeWorld(stream WorldStream) uuid.UUID {
	u := uuid.New()
	s.worldSubscriptions[u] = stream

	return u
}

func (s *Subscriber) UnsubscribeWorld(id uuid.UUID) {
	delete(s.worldSubscriptions, id)
}
