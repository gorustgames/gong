package actor

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Bat struct {
	batImage *ebiten.Image
	xPos     float64
	yPos     float64
	dx       float64
	dy       float64
	speed    float64
}

func NewBat() *Bat {

	_batImage, _, err := ebitenutil.NewImageFromFile("assets/bat00.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	return &Bat{
		xPos:     -40,
		yPos:     0,
		dx:       0,
		dy:       0,
		speed:    5,
		batImage: _batImage,
	}
}

func (b *Bat) Update() error {
	movePlayer(b)
	return nil
}

func (b *Bat) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.xPos, b.yPos)
	screen.DrawImage(b.batImage, op)
}

func movePlayer(b *Bat) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		b.yPos -= b.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		b.yPos += b.speed
	}
}
