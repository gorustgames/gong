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

const (
	PLAYER_SPEED = 6
)

type Bat struct {
	batImage            *ebiten.Image
	batHitImage         *ebiten.Image
	xPos                float64
	yPos                float64
	dx                  float64
	dy                  float64
	speed               float64
	playerLocation      PlayerLocation
	playerType          PlayerType
	telemetry           chan<- ActorTelemetry
	gameNotificationBus chan string
	showHitCounter      int
}

func NewBat(playerLocation PlayerLocation, playerType PlayerType, telemetry chan<- ActorTelemetry, gameNotificationBus chan string) *Bat {

	fileName := "assets/bat00.png"
	if playerLocation == RightPlayer {
		fileName = "assets/bat10.png"
	}
	_batImage, _, err := ebitenutil.NewImageFromFile(fileName, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	fileName = "assets/bat01.png"
	if playerLocation == RightPlayer {
		fileName = "assets/bat11.png"
	}
	_batHitImage, _, err := ebitenutil.NewImageFromFile(fileName, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	_xPos := -40.0
	if playerLocation == RightPlayer {
		_xPos = 680.0
	}

	newBat := &Bat{
		xPos:                _xPos,
		yPos:                0,
		dx:                  0,
		dy:                  0,
		speed:               PLAYER_SPEED,
		batImage:            _batImage,
		batHitImage:         _batHitImage,
		playerLocation:      playerLocation,
		playerType:          playerType,
		telemetry:           telemetry,
		gameNotificationBus: gameNotificationBus,
		showHitCounter:      0,
	}

	go func(gameNotificationBus chan string, b *Bat) {
		for gameNotification := range gameNotificationBus {
			// TODO: use constants instead of literals!
			if b.playerLocation == LeftPlayer && gameNotification == "hitLeftBat" {
				b.showHitCounter = 20
			}

			if b.playerLocation == RightPlayer && gameNotification == "hitRightBat" {
				b.showHitCounter = 20
			}
		}
	}(gameNotificationBus, newBat)

	return newBat
}

func (b *Bat) Update() error {
	movePlayer(b)

	actorType := LeftBatActor

	if b.playerLocation == RightPlayer {
		actorType = RightBatActor
	}

	b.telemetry <- ActorTelemetry{
		ActorType: actorType,
		XPos:      b.xPos,
		YPos:      b.yPos,
	}

	if b.showHitCounter > 0 {
		b.showHitCounter -= 1
	}

	return nil
}

func (b *Bat) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.xPos, b.yPos)
	if b.showHitCounter > 0 {
		screen.DrawImage(b.batHitImage, op)
	} else {
		screen.DrawImage(b.batImage, op)
	}
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

	// make sure we don't cross up/down wall
	if b.yPos < 0 {
		b.yPos = 0
	}

	if b.yPos > 320 {
		b.yPos = 320
	}
}
