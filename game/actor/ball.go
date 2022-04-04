package actor

import "github.com/hajimehoshi/ebiten"

type Ball struct {
	screen *ebiten.Image
}

func (b *Ball) Init(screen *ebiten.Image) {
	b.screen = screen
}

func (b *Ball) Update() error {
	// TODO: implement!
	return nil
}

func (b *Ball) Draw() {
	// TODO: implement!
}
