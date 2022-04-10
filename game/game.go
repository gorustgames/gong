package game

import (
	"github.com/gorustgames/gong/game/actor"
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	actors []actor.GameActor
}

const (
	SCREEN_WIDTH, SCREEN_HEIGHT = 800, 480
)

func init() {
	// not used for now
}

// game state updates
func (g *Game) Update(_ *ebiten.Image) error {

	for _, actor := range g.actors {
		actor.Update()
	}

	return nil
}

// game rendering logic
func (g *Game) Draw(screen *ebiten.Image) {
	for _, actor := range g.actors {
		actor.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = SCREEN_WIDTH
	screenHeight = SCREEN_HEIGHT
	return
}

func CreateGame() *Game {

	notificationBus := pubsub.NewBroker()
	actors := actor.CreateActors(notificationBus)

	return &Game{
		actors: actors,
	}
}
