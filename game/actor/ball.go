package actor

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Ball struct {
	ballImage *ebiten.Image
}

func NewBall() *Ball {
	_ballImage, _, err := ebitenutil.NewImageFromFile("assets/ball.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	return &Ball{
		ballImage: _ballImage,
	}
}

func (b *Ball) Update() error {
	// TODO: implement!
	return nil
}

func (b *Ball) Draw(screen *ebiten.Image) {
	// TODO: implement!
}
