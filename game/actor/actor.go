package actor

import (
	"github.com/hajimehoshi/ebiten"
)

type ActorType int8

const (
	LeftBatActor ActorType = iota
	RightBatActor
	BallActor
)

type ActorTelemetry struct {
	ActorType ActorType
	XPos      float64
	YPos      float64
}

type GameActor interface {
	// Update Responsible for update of actor state.
	Update() error

	// Draw Responsible for drawing actor onto screen.
	Draw(screen *ebiten.Image)
}

// CreateActors
// see https://stackoverflow.com/questions/17077074/array-of-pointers-to-different-struct-implementing-same-interface
func CreateActors() ([]GameActor, <-chan ActorTelemetry) {

	telemetry := make(chan ActorTelemetry)

	return []GameActor{
		NewBat(LeftPlayer, Human, telemetry),
		NewBat(RightPlayer, Human, telemetry), // TODO: this will be Human or Computer based on game mode(single player, multi player, ai2ai demo)!
		NewBall(5, telemetry),
	}, telemetry
}
