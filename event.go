package evo

import "github.com/relnod/evo/world"

// type EventType int

// const (
// 	EventRecieveWorld EventType = iota
// )

// type Event interface {
// 	Type() EventType
// 	SetData(interface{})
// }

// Stream defines a world stream.
type Stream func(*world.World)
