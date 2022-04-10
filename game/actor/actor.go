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

// CreateActorsSinglePlayer
// see https://stackoverflow.com/questions/17077074/array-of-pointers-to-different-struct-implementing-same-interface
func CreateActorsSinglePlayer(notificationBus *pubsub.Broker) []GameActor {

	return []GameActor{
		NewGameBoard(notificationBus),
		NewBat(LeftPlayer, Human, notificationBus),
		NewBat(RightPlayer, Computer, notificationBus),
		NewBall(1, notificationBus),
	}
}

func CreateActorsMultiPlayer(notificationBus *pubsub.Broker) []GameActor {

	return []GameActor{
		NewGameBoard(notificationBus),
		NewBat(LeftPlayer, Human, notificationBus),
		NewBat(RightPlayer, Human, notificationBus),
		NewBall(1, notificationBus),
	}
}

func CreateActorsMenu(notificationBus *pubsub.Broker) []GameActor {
	return []GameActor{
		NewMenu(notificationBus),
	}
}

func CreateActorsGameOver(notificationBus *pubsub.Broker) []GameActor {
	var empty []GameActor
	return empty
}
