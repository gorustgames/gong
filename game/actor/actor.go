package actor

import (
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
)

type GameActor interface {
	// Update Responsible for update of actor state.
	Update() error

	// Draw Responsible for drawing actor onto screen.
	Draw(screen *ebiten.Image)

	// Id returns actor ID
	Id() string
}

// CreateActors
// see https://stackoverflow.com/questions/17077074/array-of-pointers-to-different-struct-implementing-same-interface
func CreateActors(notificationBus *pubsub.Broker) []GameActor {

	return []GameActor{
		NewBat(LeftPlayer, Human, notificationBus),
		NewBat(RightPlayer, Computer, notificationBus), // TODO: this will be Human or Computer based on game mode(single player, multi player, ai2ai demo)!
		NewBall(1, notificationBus),
	}
}
