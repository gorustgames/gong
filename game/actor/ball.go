package actor

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math"
)

type Ball struct {
	ballImage *ebiten.Image
	xPos      float64 // position of ball
	yPos      float64
	xPosLB    float64 // position of left bat
	yPosLB    float64
	xPosRB    float64 // position of right bat
	yPosRB    float64
	dx        float64
	dy        float64
	speed     int
	telemetry chan<- ActorTelemetry
}

const (
	SCREEN_HALF_WIDTH, SCREEN_HALF_HEIGHT = 400, 240
	// deduct 12 to compensate for ball.png size (12x12 px) & padding
	BALL_CENTER_X             = SCREEN_HALF_WIDTH - 12
	BALL_CENTER_Y             = SCREEN_HALF_HEIGHT - 12
	PAD_JPG_HEIGHT_TOTAL_PX   = 160
	PAD_JPG_HEIGHT_PADONLY_PX = 128
	// this represents upper/lower margin between pad picture and whole jpg top/bottom
	PAD_JPG_PADDING_PX = (PAD_JPG_HEIGHT_TOTAL_PX - PAD_JPG_HEIGHT_PADONLY_PX) / 2
	BALL_MAX_Y         = 443
	BALL_MIN_Y         = 15
	BALL_MAX_X_BAT     = 734 // max X when bat is in front of the ball
	BALL_MIN_X_BAT     = 43  // min X when bat is in front of the ball
	BALL_MAX_X         = BALL_MAX_X_BAT + 27
	BALL_MIN_X         = BALL_MIN_X_BAT - 28
)

func NewBall(dx float64, telemetry chan<- ActorTelemetry, batsTelemetry <-chan ActorTelemetry) *Ball {
	_ballImage, _, err := ebitenutil.NewImageFromFile("assets/ball.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	newBal := &Ball{
		ballImage: _ballImage,
		xPos:      BALL_CENTER_X,
		yPos:      BALL_CENTER_Y,
		dx:        dx,
		dy:        0,
		speed:     5,
		telemetry: telemetry,
	}

	go func(telemetry <-chan ActorTelemetry, b *Ball) {
		for telemetryItem := range telemetry {
			switch telemetryItem.ActorType {
			case LeftBatActor:
				b.xPosLB = telemetryItem.XPos
				b.yPosLB = telemetryItem.YPos
				break
			case RightBatActor:
				b.xPosRB = telemetryItem.XPos
				b.yPosRB = telemetryItem.YPos
				break
			default:
				// should never happen
			}
		}
	}(batsTelemetry, newBal)

	return newBal
}

func (b *Ball) Update() error {
	moveBallManually(b)
	b.telemetry <- ActorTelemetry{
		ActorType: BallActor,
		XPos:      b.xPos,
		YPos:      b.yPos,
	}
	if b.hitLeftBat() {
		fmt.Println("hitLeftBat")
	}

	if b.hitRightBat() {
		fmt.Println("hitRightBat")
	}
	return nil
}

func (b *Ball) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.xPos, b.yPos)
	screen.DrawImage(b.ballImage, op)
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
}

func moveBallAuto(b *Ball) {
	// Each frame, we move the ball in a series of small steps.
	// The number of steps being based on its speed attribute.
	for i := 1; i <= b.speed; i++ {
		//TODO: implement
		b.xPos += b.dx
		b.yPos += b.dy
	}
}

// returns normalized vector (unit vector) with same direction as input vector (expressed as dx, dy)
// unit vector as same direction as original vector but always has same length (1 unit length)
func normalizedDxDy(dx float64, dy float64) (float64, float64) {
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
