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

	// Destroy actor tear down logic
	Destroy()

	// Destroy returns whether actor is active or not (in which case it must be removed from the game)
	IsActive() bool
}

type GameActorBase struct {
	IsActive bool
	Id       string
}

// CreateActorsSinglePlayer
// see https://stackoverflow.com/questions/17077074/array-of-pointers-to-different-struct-implementing-same-interface
func CreateActorsSinglePlayer(notificationBus *pubsub.Broker) map[string]GameActor {
	m := make(map[string]GameActor)

	a1 := NewGameBoard(notificationBus)
	a2 := NewBat(LeftPlayer, Human, notificationBus)
	a3 := NewBat(RightPlayer, Computer, notificationBus)
	a4 := NewBall(1, notificationBus)

	m[a1.Id()] = a1
	m[a2.Id()] = a2
	m[a3.Id()] = a3
	m[a4.Id()] = a4

	return m
}

func CreateActorsMultiPlayer(notificationBus *pubsub.Broker) map[string]GameActor {
	m := make(map[string]GameActor)

	a1 := NewGameBoard(notificationBus)
	a2 := NewBat(LeftPlayer, Human, notificationBus)
	a3 := NewBat(RightPlayer, Human, notificationBus)
	a4 := NewBall(1, notificationBus)

	m[a1.Id()] = a1
	m[a2.Id()] = a2
	m[a3.Id()] = a3
	m[a4.Id()] = a4

	return m
}

func CreateActorsMenu(notificationBus *pubsub.Broker) map[string]GameActor {
	m := make(map[string]GameActor)

	a1 := NewMenu(notificationBus)
	m[a1.Id()] = a1

	return m
}

func CreateActorsGameOver(notificationBus *pubsub.Broker) map[string]GameActor {
	m := make(map[string]GameActor)

	a1 := NewGameOver(notificationBus)
	m[a1.Id()] = a1

	return m
}
