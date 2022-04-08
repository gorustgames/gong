package actor

import (
	"github.com/gorustgames/gong/gamebus"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math"
)

type Ball struct {
	ballImage       *ebiten.Image
	xPos            float64 // position of ball
	yPos            float64
	xPosLB          float64 // position of left bat
	yPosLB          float64
	xPosRB          float64 // position of right bat
	yPosRB          float64
	dx              float64
	dy              float64
	speed           int
	notificationBus *gamebus.GameNotificationBus
}

const (
	SCREEN_HALF_WIDTH, SCREEN_HALF_HEIGHT = 400, 240
	// deduct 12 to compensate for ball.png size (12x12 px) & padding
	BALL_CENTER_X                  = SCREEN_HALF_WIDTH - 12
	BALL_CENTER_Y                  = SCREEN_HALF_HEIGHT - 12
	PAD_JPG_HEIGHT_TOTAL_PX        = 160
	PAD_JPG_HEIGHT_PADONLY_PX      = 128
	PAD_JPG_HEIGHT_HALF_PADONLY_PX = PAD_JPG_HEIGHT_PADONLY_PX / 2
	// this represents upper/lower margin between pad picture and whole jpg top/bottom
	PAD_JPG_PADDING_PX = (PAD_JPG_HEIGHT_TOTAL_PX - PAD_JPG_HEIGHT_PADONLY_PX) / 2
	BALL_MAX_Y         = 443
	BALL_MIN_Y         = 15
	BALL_MAX_X_BAT     = 734 // max X when bat is in front of the ball
	BALL_MIN_X_BAT     = 43  // min X when bat is in front of the ball
	BALL_MAX_X         = BALL_MAX_X_BAT + 27
	BALL_MIN_X         = BALL_MIN_X_BAT - 28
	BALL_SPEED         = 1
)

func NewBall(dx float64, notificationBus *gamebus.GameNotificationBus) *Ball {
	_ballImage, _, err := ebitenutil.NewImageFromFile("assets/ball.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	newBall := &Ball{
		ballImage:       _ballImage,
		xPos:            BALL_CENTER_X,
		yPos:            BALL_CENTER_Y,
		dx:              dx,
		dy:              0,
		speed:           BALL_SPEED,
		notificationBus: notificationBus,
	}

	go func(b *Ball) {
		for notification := range b.notificationBus.Bus {
			switch notification.ActorType {
			case gamebus.LeftBatActor:
				switch v := notification.Data.(type) {
				case gamebus.PositionNotificationPayload:
					b.xPosLB = v.XPos
					b.yPosLB = v.YPos
				}
				break
			case gamebus.RightBatActor:
				switch v := notification.Data.(type) {
				case gamebus.PositionNotificationPayload:
					b.xPosRB = v.XPos
					b.yPosRB = v.YPos
				}
				break
			}
		}
	}(newBall)

	return newBall
}

func (b *Ball) Update() error {
	// moveBallManually(b)
	moveBallAuto(b)

	b.notificationBus.Bus <- gamebus.GameNotification{
		ActorType:            gamebus.BallActor,
		GameNotificationType: gamebus.PositionNotification,
		Data: gamebus.PositionNotificationPayload{
			XPos: b.xPos,
			YPos: b.yPos,
		},
	}

	return nil
}

func (b *Ball) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.xPos, b.yPos)
	screen.DrawImage(b.ballImage, op)
}

func (b *Ball) Id() string {
	return "actor-ball"
}

// will be used only for debugging. In real game ball
// is moving by laws of physics ;)
func moveBallManually(b *Ball) {

	if ebiten.IsKeyPressed(ebiten.KeyY) {
		b.yPos -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		b.yPos += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		b.xPos += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		b.xPos -= 1
	}

	if b.xPos >= BALL_MAX_X {
		b.xPos = BALL_MAX_X
	}

	if b.xPos <= BALL_MIN_X {
		b.xPos = BALL_MIN_X
	}

	if b.yPos >= BALL_MAX_Y {
		b.yPos = BALL_MAX_Y
	}

	if b.yPos <= BALL_MIN_Y {
		b.yPos = BALL_MIN_Y
	}

	if b.hitLeftBat() {
		b.notificationBus.Bus <- gamebus.GameNotification{
			ActorType:            gamebus.BallActor,
			GameNotificationType: gamebus.LeftBatHitNotification,
			Data:                 nil,
		}
		b.xPos = BALL_MIN_X_BAT
	}

	if b.hitRightBat() {
		b.notificationBus.Bus <- gamebus.GameNotification{
			ActorType:            gamebus.BallActor,
			GameNotificationType: gamebus.RightBatHitNotification,
			Data:                 nil,
		}
		b.xPos = BALL_MAX_X_BAT
	}
}

func moveBallAuto(b *Ball) {
	// Each frame, we move the ball in a series of small steps.
	// The number of steps being based on its speed attribute.
	for i := 1; i <= b.speed; i++ {
		b.xPos += b.dx
		b.yPos += b.dy

		// TODO: for now bounce off but this means new set in real game
		if b.hitBottom() || b.hitTop() {
			b.dy = -b.dy
		}

		// TODO: for now bounce off but this means new set in real game
		if b.hitLeft() || b.hitRight() {
			b.dx = -b.dx
		}

		if b.hitLeftBat() {
			b.dx = -b.dx
			b.dy += b.deflectionForLeftBat()
			b.notificationBus.Bus <- gamebus.GameNotification{
				ActorType:            gamebus.BallActor,
				GameNotificationType: gamebus.LeftBatHitNotification,
				Data:                 nil,
			}
			b.xPos = BALL_MIN_X_BAT
		}

		if b.hitRightBat() {
			b.dx = -b.dx
			b.dy += b.deflectionForRightBat()
			b.notificationBus.Bus <- gamebus.GameNotification{
				ActorType:            gamebus.BallActor,
				GameNotificationType: gamebus.RightBatHitNotification,
				Data:                 nil,
			}
			b.xPos = BALL_MAX_X_BAT
		}

		// Ensure our direction vector is a unit vector, i.e. represents a distance
		// of the equivalent of 1 pixel regardless of its angle.
		b.dx, b.dy = intoUnitVector(b.dx, b.dy)
	}
}

// returns normalized vector (unit vector) with same direction as input vector (expressed as dx, dy)
// unit vector as same direction as original vector but always has same length (1 unit length)
// see https://en.wikipedia.org/wiki/Unit_vector
func intoUnitVector(dx float64, dy float64) (float64, float64) {
	vecLen := math.Hypot(dx, dy)
	return dx / vecLen, dy / vecLen
}

func (b *Ball) hitLeftBat() bool {
	if b.xPos >= BALL_CENTER_X {
		return false
	}

	ballUpperBound := b.yPosLB + PAD_JPG_PADDING_PX
	ballLowerBound := b.yPosLB + PAD_JPG_PADDING_PX + PAD_JPG_HEIGHT_PADONLY_PX

	return b.yPos >= ballUpperBound && b.yPos <= ballLowerBound && b.xPos <= BALL_MIN_X_BAT
}

func (b *Ball) hitRightBat() bool {
	if b.xPos <= BALL_CENTER_X {
		return false
	}

	ballUpperBound := b.yPosRB + PAD_JPG_PADDING_PX
	ballLowerBound := b.yPosRB + PAD_JPG_PADDING_PX + PAD_JPG_HEIGHT_PADONLY_PX

	return b.yPos >= ballUpperBound && b.yPos <= ballLowerBound && b.xPos >= BALL_MAX_X_BAT
}

// calculates change in dy of left bat is hit based on where we did hit it (upper or lowe part)
func (b *Ball) deflectionForLeftBat() float64 {
	diffY := b.yPosLB + PAD_JPG_PADDING_PX - b.yPos // ball padding is small & ignored here
	if diffY < PAD_JPG_HEIGHT_HALF_PADONLY_PX {
		diffY = -diffY
	}
	deflection := diffY / PAD_JPG_HEIGHT_PADONLY_PX
	return deflection
}

// calculates change in dy of right bat is hit based on where we did hit it (upper or lowe part)
func (b *Ball) deflectionForRightBat() float64 {
	diffY := b.yPosRB + PAD_JPG_PADDING_PX - b.yPos // ball padding is small & ignored here
	if diffY < PAD_JPG_HEIGHT_HALF_PADONLY_PX {
		diffY = -diffY
	}
	deflection := diffY / PAD_JPG_HEIGHT_PADONLY_PX
	return deflection
}

func (b *Ball) hitLeft() bool {
	return b.xPos <= BALL_MIN_X
}

func (b *Ball) hitRight() bool {
	return b.xPos >= BALL_MAX_X
}

func (b *Ball) hitTop() bool {
	return b.yPos <= BALL_MIN_Y
}

func (b *Ball) hitBottom() bool {
	return b.yPos >= BALL_MAX_Y
}
