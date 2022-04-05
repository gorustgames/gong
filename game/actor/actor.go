package actor

import (
	"github.com/hajimehoshi/ebiten"
)

type GameActor interface {
	// Init Actor initialization.
	Init()

	// Update Responsible for update of actor state.
	Update() error

	// Draw Responsible for drawing actor onto screen.
	Draw(screen *ebiten.Image)
}

// CreateActors
// see https://stackoverflow.com/questions/17077074/array-of-pointers-to-different-struct-implementing-same-interface
func CreateActors() []GameActor {
	var batPlayer *Bat
	batPlayer = new(Bat)
	batPlayer.Init()

	var ball *Ball
	ball = new(Ball)
	ball.Init()

	return []GameActor{batPlayer, ball}
}
