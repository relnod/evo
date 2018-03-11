package platform

import "github.com/relnod/evo/num"

type KeyListener func(e *KeyEvent)

type KeyEvent struct {
	KeyCode int
}

type MouseListener func(e *MouseEvent)

type MouseEvent struct {
	Pos num.Vec2
}
