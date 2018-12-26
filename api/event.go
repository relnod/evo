package api

import (
	"encoding/json"

	"github.com/google/uuid"
)

// EventType defines the type of an event.
type EventType int

// All event types.
const (
	EventWorld        EventType = iota
	EventSubscription EventType = iota
)

type SubscriptionType int

const (
	SubscriptionWorld SubscriptionType = iota
)

// Event defines an api event.
type Event struct {
	Type    EventType       `json:"type"`
	Message json.RawMessage `json:"message"`
}

type Subscription struct {
	Type SubscriptionType
	ID   uuid.UUID
}
