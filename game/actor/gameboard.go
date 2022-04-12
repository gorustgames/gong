package actor

import (
	"fmt"
	"github.com/gorustgames/gong/game/util"
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"log"
	"math"
)

type GameBoard struct {
	base             GameActorBase
	background       *ebiten.Image
	digitsLeft       [10]*ebiten.Image
	digitsRight      [10]*ebiten.Image
	xLB              float64 // xPos of left bat
	yLB              float64 // yPos of left bat
	xRB              float64 // xPos of right bat
	yRB              float64 // yPos of right bat
	xB               float64 // xPos of ball
	yB               float64 // xPos of ball
	leftScore        int
	rightScore       int
	notificationBus  *pubsub.Broker
	subscribers      []*pubsub.Subscriber
	audioPlayerScore *audio.Player
}

func NewGameBoard(notificationBus *pubsub.Broker) *GameBoard {
	_background, _, err := ebitenutil.NewImageFromFile("assets/table.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	audioPlayerScore := util.NewAudioPlayer("score")

	newGameBoard := &GameBoard{
		base: GameActorBase{
			IsActive: true,
			Id:       fmt.Sprintf("actor-gameboard-%s", util.GenerateShortId()),
		},
		background:       _background,
		xLB:              0,
		yLB:              0,
		xRB:              0,
		yRB:              0,
		xB:               0,
		yB:               0,
		leftScore:        0,
		rightScore:       0,
		notificationBus:  notificationBus,
		audioPlayerScore: audioPlayerScore,
	}

	for i := 0; i < 10; i++ {
		newGameBoard.digitsLeft[i], _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("assets/digit1%d.png", i), ebiten.FilterDefault)
		newGameBoard.digitsRight[i], _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("assets/digit2%d.png", i), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
	}

	subscriberPos := notificationBus.AddSubscriber("gameboard-subscriberPos")
	notificationBus.Subscribe(subscriberPos, pubsub.POSITION_NOTIFICATION_TOPIC)
	go subscriberPos.Listen(newGameBoard.updatePositions)

	subscriberLeftBatMiss := notificationBus.AddSubscriber("gameboard-subscriberLeftBatMiss")
	notificationBus.Subscribe(subscriberLeftBatMiss, pubsub.LEFT_BAT_MISS_NOTIFICATION_TOPIC)
	go subscriberLeftBatMiss.Listen(newGameBoard.leftBatMiss)

	subscriberRightBatMiss := notificationBus.AddSubscriber("gameboard-subscriberRightBatMiss")
	notificationBus.Subscribe(subscriberRightBatMiss, pubsub.RIGHT_BAT_MISS_NOTIFICATION_TOPIC)
	go subscriberRightBatMiss.Listen(newGameBoard.rightBatMiss)

	newGameBoard.subscribers = make([]*pubsub.Subscriber, 3)
	newGameBoard.subscribers[0] = subscriberPos
	newGameBoard.subscribers[1] = subscriberLeftBatMiss
	newGameBoard.subscribers[2] = subscriberRightBatMiss

	return newGameBoard
}

func (g *GameBoard) Update() error {
	return nil
}

func (g *GameBoard) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(g.background, op)

	// left score
	opLS := &ebiten.DrawImageOptions{}
	opLS.GeoM.Translate(400-75, 50) // digit jpg is 75x75

	// right score
	opRS := &ebiten.DrawImageOptions{}
	opRS.GeoM.Translate(400, 50)

	idxL := int(math.Min(float64(g.leftScore), 9))
	idxR := int(math.Min(float64(g.rightScore), 9))
	screen.DrawImage(g.digitsLeft[idxL], opLS)
	screen.DrawImage(g.digitsRight[idxR], opRS)

	// debug print of positions of crucial game actors
	/*ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("LB: x = %f, y = %f RB: x = %f, y = %f B: x = %f, y = %f.   Score: %d:%d",
			g.xLB,
			g.yLB,
			g.xRB,
			g.yRB,
			g.xB,
			g.yB,
			g.leftScore,
			g.rightScore,
		),
	)*/
}

func (g *GameBoard) Id() string {
	return g.base.Id
}

func (g *GameBoard) Destroy() {
	for _, subscriber := range g.subscribers {
		g.notificationBus.RemoveSubscriber(subscriber)
	}
}

func (g *GameBoard) IsActive() bool {
	return true
}

func (g *GameBoard) playScore() {
	g.audioPlayerScore.Rewind()
	g.audioPlayerScore.Play()
}

func (g *GameBoard) checkScore(score int) {
	if score > 9 {
		g.notificationBus.Publish(pubsub.CHANGE_GAME_STATE_GAME_OVER_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.GameOverActor,
			Data: pubsub.GameOverNotificationPayload{
				ScoreLeft:  g.leftScore,
				ScoreRight: g.rightScore,
			},
		})
	}
}

func (g *GameBoard) leftBatMiss(message *pubsub.Message) {
	g.leftScore += 1
	g.playScore()
	g.checkScore(g.leftScore)
}

func (g *GameBoard) rightBatMiss(message *pubsub.Message) {
	g.rightScore += 1
	g.playScore()
	g.checkScore(g.rightScore)
}

func (g *GameBoard) updatePositions(message *pubsub.Message) {
	switch message.GetMessageBody().ActorType {
	case pubsub.LeftBatActor:
		g.updatePositionOfLeftBat(message)
		break
	case pubsub.RightBatActor:
		g.updatePositionOfRightBat(message)
		break
	case pubsub.BallActor:
		g.updatePositionOfBall(message)
		break
	}
}

func (g *GameBoard) updatePositionOfLeftBat(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		g.xLB = v.XPos
		g.yLB = v.YPos
	}
}

func (g *GameBoard) updatePositionOfRightBat(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		g.xRB = v.XPos
		g.yRB = v.YPos
	}
}

func (g *GameBoard) updatePositionOfBall(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		g.xB = v.XPos
		g.yB = v.YPos
	}
}
