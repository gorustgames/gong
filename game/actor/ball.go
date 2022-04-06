package actor

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math"
)

type Ball struct {
	ballImage *ebiten.Image
	xPos      float64
	yPos      float64
	dx        float64
	dy        float64
	speed     int
	telemetry chan<- ActorTelemetry
}

const (
	SCREEN_HALF_WIDTH, SCREEN_HALF_HEIGHT = 400, 240
	BALL_MAX_Y                            = 443
	BALL_MIN_Y                            = 15
	BALL_MAX_X_BAT                        = 734 // max X when bat is in front of the ball
	BALL_MIN_X_BAT                        = 43  // min X when bat is in front of the ball
	BALL_MAX_X                            = BALL_MAX_X_BAT + 27
	BALL_MIN_X                            = BALL_MIN_X_BAT - 28
)

func NewBall(dx float64, telemetry chan<- ActorTelemetry) *Ball {
	_ballImage, _, err := ebitenutil.NewImageFromFile("assets/ball.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	return &Ball{
		ballImage: _ballImage,
		xPos:      SCREEN_HALF_WIDTH - 12, // deduct 12 to compensate for ball.png size (12x12 px) & padding
		yPos:      SCREEN_HALF_HEIGHT - 12,
		dx:        dx,
		dy:        0,
		speed:     5,
		telemetry: telemetry,
	}
}

func (b *Ball) Update() error {
	moveBallManually(b)
	b.telemetry <- ActorTelemetry{
		ActorType: BallActor,
		XPos:      b.xPos,
		YPos:      b.yPos,
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
