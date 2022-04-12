package actor

import (
	"fmt"
	"github.com/gorustgames/gong/game/util"
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math"
)

type GameOver struct {
	base            GameActorBase
	digitsLeft      [10]*ebiten.Image
	digitsRight     [10]*ebiten.Image
	picture         *ebiten.Image
	notificationBus *pubsub.Broker
	scoreLeft       int
	scoreRight      int
}

func NewGameOver(notificationBus *pubsub.Broker, scoreLeft int, scoreRight int) *GameOver {
	picture, _, err := ebitenutil.NewImageFromFile("assets/over.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	newGameOver := &GameOver{
		base: GameActorBase{
			IsActive: true,
			Id:       fmt.Sprintf("actor-gameover-%s", util.GenerateShortId()),
		},
		picture:         picture,
		notificationBus: notificationBus,
		scoreLeft:       scoreLeft,
		scoreRight:      scoreRight,
	}

	for i := 0; i < 10; i++ {
		newGameOver.digitsLeft[i], _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("assets/digit1%d.png", i), ebiten.FilterDefault)
		newGameOver.digitsRight[i], _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("assets/digit2%d.png", i), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
	}

	return newGameOver
}

func (a *GameOver) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		a.notificationBus.Publish(pubsub.CHANGE_GAME_STATE_MENU_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.GameOverActor,
			Data:      nil,
		})
	}
	return nil
}

func (a *GameOver) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(a.picture, op)

	// left score
	opLS := &ebiten.DrawImageOptions{}
	opLS.GeoM.Translate(400-75, 50) // digit jpg is 75x75

	// right score
	opRS := &ebiten.DrawImageOptions{}
	opRS.GeoM.Translate(400, 50)

	idxL := int(math.Min(float64(a.scoreLeft), 9))
	idxR := int(math.Min(float64(a.scoreRight), 9))
	screen.DrawImage(a.digitsLeft[idxL], opLS)
	screen.DrawImage(a.digitsRight[idxR], opRS)

}

func (a *GameOver) Id() string {
	return a.base.Id
}

func (a *GameOver) Destroy() {
	// nothing to do
}

func (a *GameOver) IsActive() bool {
	return true
}
