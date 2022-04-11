package actor

import (
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type GameOver struct {
	base            GameActorBase
	picture         *ebiten.Image
	notificationBus *pubsub.Broker
}

func NewGameOver(notificationBus *pubsub.Broker) *GameOver {
	picture, _, err := ebitenutil.NewImageFromFile("assets/over.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	newGameOver := &GameOver{
		base: GameActorBase{
			IsActive: true,
		},
		picture:         picture,
		notificationBus: notificationBus,
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
}

func (a *GameOver) Id() string {
	return "actor-gameover"
}

func (a *GameOver) Destroy() {
	// nothing to do
}

func (a *GameOver) IsActive() bool {
	return true
}
