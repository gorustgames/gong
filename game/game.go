package game

import (
	"github.com/gorustgames/gong/game/actor"
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
		actor.Draw()
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = SCREEN_WIDTH
	screenHeight = SCREEN_HEIGHT
	return
}

func CreateGame() *Game {
	return &Game{
		actors: actor.CreateActors(background),
	}
}
