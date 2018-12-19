package api

// EventType defines the type of an event.
type EventType int

// All event types.
const (
	World EventType = iota
)

// Event defines an api event.
type Event struct {
	Type    EventType `json:"type"`
	Message []byte    `json:"message"`
}
