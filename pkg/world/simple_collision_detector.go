package world

import (
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

// SimpleCollisionDetector takes a simple aproach in detecting the collision of
// creatures with creatures, by brute forcing every combination.
// Implements the evo.CollisionHandler
type SimpleCollisionDetector struct {
	width  int
	height int
}

// NewSimpleCollisionDetector returns a new simpe collisio updater.
func NewSimpleCollisionDetector(width, height int) *SimpleCollisionDetector {
	return &SimpleCollisionDetector{
		width:  width,
		height: height,
	}
}

// DetectCollisions checks the collision for all creatures.
func (s *SimpleCollisionDetector) DetectCollisions(creatures []*entity.Creature) []Collision {
	var collisions []Collision
	for _, c := range creatures {
		// Check if the creature is outside the world boundaries.
		if c.Pos.X < 0.0 {
			collisions = append(collisions, &creatureBorderCollision{c, collision.LEFT, s.width, s.height})
		} else if c.Pos.X > float64(s.width) {
			collisions = append(collisions, &creatureBorderCollision{c, collision.RIGHT, s.width, s.height})
		} else if c.Pos.Y < 0.0 {
			collisions = append(collisions, &creatureBorderCollision{c, collision.TOP, s.width, s.height})
		} else if c.Pos.Y > float64(s.height) {
			collisions = append(collisions, &creatureBorderCollision{c, collision.BOT, s.width, s.height})
		}

		// If the creature is moving, we have to check for a collision with
		// another creature.
		if c.Speed > 0 {
			for _, c2 := range creatures {
				if c == c2 {
					continue
				}
				if collision.CircleCircle(&c.Pos, c.Radius, &c2.Pos, c2.Radius) {
					collisions = append(collisions, &creatureCreatureCollision{c, c2})
				}
			}
		}

		// If the creature has eyes, check if the eyes see anything.
		if len(c.Eyes) > 0 {
			for _, eye := range c.Eyes {
				for _, c2 := range creatures {
					if c == c2 {
						continue
					}

					d := math64.Vec2{X: c2.Pos.X - c.Pos.X, Y: c2.Pos.Y - c.Pos.Y}
					// Check if the other creature is in range of the eye.
					if d.Len()-c2.Radius > eye.Range {
						continue
					}

					// Check if the the other creature is in the fov of the eye.
					if math64.Angle(&d, &c.Dir) > eye.FOV/2 {
						continue
					}

					collisions = append(collisions, &eyeCreatureCollision{eye, c2})
				}
			}
		}
	}
	return collisions
}
