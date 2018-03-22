package api

import "github.com/relnod/evo/num"

type EventType int

const (
	TGetWorld EventType = iota
	TGetEntityAt
)

type Event struct {
	Type    EventType `json:"type"`
	Message []byte    `json:"message"`
}

type GetEntityAt struct {
	Pos num.Vec2 `json:"pos"`
}
