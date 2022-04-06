package game

import (
	"fmt"
	"github.com/gorustgames/gong/game/actor"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Game struct {
	actors []actor.GameActor
}

const (
	SCREEN_WIDTH, SCREEN_HEIGHT = 800, 480
)

// TODO: this should be moved into game struct probably (with background exception)!
var (
	background    *ebiten.Image
	xLB           float64 // xPos of left bat
	yLB           float64 // yPos of left bat
	xRB           float64 // xPos of right bat
	yRB           float64 // yPos of right bat
	xB            float64 // xPos of ball
	yB            float64 // xPos of ball
	batsTelemetry chan<- actor.ActorTelemetry
)

func init() {
	_background, _, err := ebitenutil.NewImageFromFile("assets/table.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	background = _background
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

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	for _, actor := range g.actors {
		actor.Draw(screen)
	}

	// debug print of positions of crucial game actors
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("LB: x = %f, y = %f RB: x = %f, y = %f B: x = %f, y = %f",
			xLB,
			yLB,
			xRB,
			yRB,
			xB,
			yB,
		),
	)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = SCREEN_WIDTH
	screenHeight = SCREEN_HEIGHT
	return
}

func readActorTelemetry(telemetry <-chan actor.ActorTelemetry) {
	for telemetryItem := range telemetry {
		switch telemetryItem.ActorType {
		case actor.LeftBatActor:
			xLB = telemetryItem.XPos
			yLB = telemetryItem.YPos
			// forward bat position to ball
			batsTelemetry <- telemetryItem
			break
		case actor.RightBatActor:
			xRB = telemetryItem.XPos
			yRB = telemetryItem.YPos
			// forward bat position to ball
			batsTelemetry <- telemetryItem
			break
		case actor.BallActor:
			xB = telemetryItem.XPos
			yB = telemetryItem.YPos
			break
		}
	}
}

func CreateGame() *Game {

	actors, telemetry, _batsTelemetry := actor.CreateActors()
	go readActorTelemetry(telemetry)

	batsTelemetry = _batsTelemetry

	return &Game{
		actors: actors,
	}
}
