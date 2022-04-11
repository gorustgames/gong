package actor

import (
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Menu struct {
	base            GameActorBase
	picture0        *ebiten.Image
	picture1        *ebiten.Image
	singlePlayer    bool
	notificationBus *pubsub.Broker
}

func NewMenu(notificationBus *pubsub.Broker) *Menu {
	picture0, _, err := ebitenutil.NewImageFromFile("assets/menu0.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	picture1, _, err := ebitenutil.NewImageFromFile("assets/menu1.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	newMenu := &Menu{
		base: GameActorBase{
			IsActive: true,
		},
		picture0:        picture0,
		picture1:        picture1,
		singlePlayer:    true,
		notificationBus: notificationBus,
	}

	return newMenu
}

func (m *Menu) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		m.singlePlayer = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		m.singlePlayer = false
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		gameModeTopic := pubsub.CHANGE_GAME_STATE_MULTI_PLAYER_TOPIC
		if m.singlePlayer {
			gameModeTopic = pubsub.CHANGE_GAME_STATE_SINGLE_PLAYER_TOPIC
		}

		m.notificationBus.Publish(gameModeTopic, pubsub.GameNotification{
			ActorType: pubsub.MenuActor,
			Data:      nil,
		})

	}
	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	if m.singlePlayer {
		screen.DrawImage(m.picture0, op)
	} else {
		screen.DrawImage(m.picture1, op)
	}

}

func (m *Menu) Id() string {
	return "actor-menu"
}

func (m *Menu) Destroy() {
	// nothing to do
}

func (m *Menu) IsActive() bool {
	return true
}
