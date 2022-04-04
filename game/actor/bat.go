package actor

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Bat struct {
	screen   *ebiten.Image
	batImage *ebiten.Image
	xPos     float64
	yPos     float64
	dx       float64
	dy       float64
}

var (
	batImage *ebiten.Image
)

func (b *Bat) Init(screen *ebiten.Image) {
	b.screen = screen
	b.xPos = 0
	b.yPos = 0
	b.dx = 0
	b.dy = 0

	batImage, _, err := ebitenutil.NewImageFromFile("assets/bat00.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	b.batImage = batImage

}

func (b *Bat) Update() error {
	// TODO: implement!
	return nil
}

func (b *Bat) Draw() {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.xPos, b.yPos)
	b.screen.DrawImage(b.batImage, op)
}
