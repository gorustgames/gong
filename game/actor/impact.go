package actor

import (
	"fmt"
	"github.com/gorustgames/gong/game/util"
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Impact struct {
	base            GameActorBase
	pictures        []*ebiten.Image
	xPos            float64
	yPos            float64
	currentPicture  int
	notificationBus *pubsub.Broker
}

func NewImpact(xPos float64, yPos float64, notificationBus *pubsub.Broker) *Impact {
	picture0, _, err := ebitenutil.NewImageFromFile("assets/impact0.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	picture1, _, err := ebitenutil.NewImageFromFile("assets/impact1.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	picture2, _, err := ebitenutil.NewImageFromFile("assets/impact2.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	picture3, _, err := ebitenutil.NewImageFromFile("assets/impact3.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	picture4, _, err := ebitenutil.NewImageFromFile("assets/impact4.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	newImpact := &Impact{
		base: GameActorBase{
			IsActive: true,
			Id:       fmt.Sprintf("actor-impact-%s", util.GenerateShortId()),
		},
		pictures:        []*ebiten.Image{picture0, picture1, picture2, picture3, picture4},
		xPos:            xPos,
		yPos:            yPos,
		currentPicture:  0,
		notificationBus: notificationBus,
	}

	return newImpact
}

func (a *Impact) Update() error {

	a.currentPicture += 1
	if a.currentPicture > 4 {
		a.base.IsActive = false
	}
	return nil
}

func (a *Impact) Draw(screen *ebiten.Image) {
	if a.base.IsActive {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(a.xPos, a.yPos)
		screen.DrawImage(a.pictures[a.currentPicture], op)
	}
}

func (a *Impact) Id() string {
	return a.base.Id
}

func (a *Impact) Destroy() {
	// nothing to do
}

func (a *Impact) IsActive() bool {
	return a.base.IsActive
}
