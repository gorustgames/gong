package game

import (
	"fmt"
	"github.com/gorustgames/gong/game/actor"
	"github.com/gorustgames/gong/gamebus"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Game struct {
	actors          []actor.GameActor
	notificationBus *gamebus.GameNotificationBus
}

const (
	SCREEN_WIDTH, SCREEN_HEIGHT = 800, 480
)

var (
	background *ebiten.Image
	xLB        float64 // xPos of left bat
	yLB        float64 // yPos of left bat
	xRB        float64 // xPos of right bat
	yRB        float64 // yPos of right bat
	xB         float64 // xPos of ball
	yB         float64 // xPos of ball
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

func readActorTelemetry(bus *gamebus.GameNotificationBus) {
	for notification := range bus.Bus {
		switch notification.ActorType {
		case gamebus.LeftBatActor:
			switch v := notification.Data.(type) {
			case gamebus.PositionNotificationPayload:
				xLB = v.XPos
				yLB = v.YPos
			}
			break
		case gamebus.RightBatActor:
			switch v := notification.Data.(type) {
			case gamebus.PositionNotificationPayload:
				xRB = v.XPos
				yRB = v.YPos
			}
			break
		case gamebus.BallActor:
			switch v := notification.Data.(type) {
			case gamebus.PositionNotificationPayload:
				xB = v.XPos
				yB = v.YPos
			}
			break
		}
	}
}

func CreateGame() *Game {

	// TODO: _batsTelemetry can be propagated via gameNotificationBus as well!
	notificationBus := gamebus.NewGameNotificationBus()

	actors := actor.CreateActors(notificationBus)
	go readActorTelemetry(notificationBus)

	return &Game{
		actors:          actors,
		notificationBus: notificationBus,
	}
}
