package actor

import (
	"github.com/hajimehoshi/ebiten"
)

type GameActor interface {
	// Init Actor initialization. When being initiated we pass pointer to screen so that actor can draw its state on it when asked.
	Init(screen *ebiten.Image)

	// Update Responsible for update of actor state.
	Update() error

	// Draw Responsible for drawing actor onto screen.
	Draw()
}

// CreateActors
// see https://stackoverflow.com/questions/17077074/array-of-pointers-to-different-struct-implementing-same-interface
func CreateActors(screen *ebiten.Image) []GameActor {
	var batPlayer *Bat
	batPlayer = new(Bat)
	batPlayer.Init(screen)

	var ball *Ball
	ball = new(Ball)
	ball.Init(screen)

	return []GameActor{batPlayer, ball}
}
