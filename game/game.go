package game

import (
	"fmt"
	"github.com/gorustgames/gong/game/actor"
	"github.com/gorustgames/gong/pubsub"
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

func updatePosition(message *pubsub.Message) {
	switch message.GetMessageBody().ActorType {
	case pubsub.LeftBatActor:
		updatePositionOfLeftBat(message)
		break
	case pubsub.RightBatActor:
		updatePositionOfRightBat(message)
		break
	case pubsub.BallActor:
		updatePositionOfBall(message)
		break
	}
}

func updatePositionOfLeftBat(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		xLB = v.XPos
		yLB = v.YPos
	}
}

func updatePositionOfRightBat(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		xRB = v.XPos
		yRB = v.YPos
	}
}

func updatePositionOfBall(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		xB = v.XPos
		yB = v.YPos
	}
}

func CreateGame() *Game {

	notificationBus := pubsub.NewBroker()
	actors := actor.CreateActors(notificationBus)

	subscriberPos := notificationBus.AddSubscriber()
	notificationBus.Subscribe(subscriberPos, pubsub.POSITION_NOTIFICATION_TOPIC)
	go subscriberPos.Listen(updatePosition)

	return &Game{
		actors: actors,
	}
}
