package actor

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type PlayerLocation int8

const (
	LeftPlayer PlayerLocation = iota
	RightPlayer
)

type PlayerType int8

const (
	Human PlayerType = iota
	Computer
)

type Bat struct {
	batImage       *ebiten.Image
	xPos           float64
	yPos           float64
	dx             float64
	dy             float64
	speed          float64
	playerLocation PlayerLocation
	playerType     PlayerType
}

func NewBat(playerLocation PlayerLocation, playerType PlayerType) *Bat {

	_batImage, _, err := ebitenutil.NewImageFromFile("assets/bat00.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	_xPos := -40.0
	if playerLocation == RightPlayer {
		_xPos = 690.0
	}

	return &Bat{
		xPos:           _xPos,
		yPos:           0,
		dx:             0,
		dy:             0,
		speed:          5,
		batImage:       _batImage,
		playerLocation: playerLocation,
		playerType:     playerType,
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
	if b.playerType == Human {
		if b.playerLocation == LeftPlayer {
			if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyA) {
				b.yPos -= b.speed
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyZ) {
				b.yPos += b.speed
			}
		} else {
			if ebiten.IsKeyPressed(ebiten.KeyK) {
				b.yPos -= b.speed
			}
			if ebiten.IsKeyPressed(ebiten.KeyM) {
				b.yPos += b.speed
			}
		}
	} else {
		// TODO: implement AI!
	}
}
