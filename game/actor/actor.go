package actor

import (
	"github.com/hajimehoshi/ebiten"
)

type GameActor interface {
	// Update Responsible for update of actor state.
	Update() error

	// Draw Responsible for drawing actor onto screen.
	Draw(screen *ebiten.Image)
}

// CreateActors
// see https://stackoverflow.com/questions/17077074/array-of-pointers-to-different-struct-implementing-same-interface
func CreateActors() []GameActor {
	return []GameActor{
		NewBat(LeftPlayer, Human),
		NewBat(RightPlayer, Human), // TODO: switch to Computer here!
		NewBall(0),                 //for now ball will be not moving initially
	}
}
