package actor

import (
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math"
)

type Ball struct {
	base            GameActorBase
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
	notificationBus *pubsub.Broker
	subscribers     []*pubsub.Subscriber
}

const (
	SCREEN_HALF_WIDTH, SCREEN_HALF_HEIGHT = 400, 240

	BALL_JPG_HEIGHT_TOTAL_PX = 24

	// deduct 12 to compensate for ball.png size (24x24 px) & padding
	BALL_CENTER_X                  = SCREEN_HALF_WIDTH - BALL_JPG_HEIGHT_TOTAL_PX/2
	BALL_CENTER_Y                  = SCREEN_HALF_HEIGHT - BALL_JPG_HEIGHT_TOTAL_PX/2
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
	BALL_SPEED         = 5
)

func NewBall(dx float64, notificationBus *pubsub.Broker) *Ball {
	_ballImage, _, err := ebitenutil.NewImageFromFile("assets/ball.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	newBall := &Ball{
		base: GameActorBase{
			IsActive: true,
		},
		ballImage:       _ballImage,
		xPos:            BALL_CENTER_X,
		yPos:            BALL_CENTER_Y,
		dx:              dx,
		dy:              0,
		speed:           BALL_SPEED,
		notificationBus: notificationBus,
	}

	subscriberPos := notificationBus.AddSubscriber("ball-subscriberPos")

	notificationBus.Subscribe(subscriberPos, pubsub.POSITION_NOTIFICATION_TOPIC)
	go subscriberPos.Listen(newBall.updatePosition)

	newBall.subscribers = make([]*pubsub.Subscriber, 1)
	newBall.subscribers[0] = subscriberPos

	return newBall
}

func (b *Ball) Update() error {
	// moveBallManually(b)
	moveBallAuto(b)

	b.notificationBus.Publish(pubsub.POSITION_NOTIFICATION_TOPIC, pubsub.GameNotification{
		ActorType: pubsub.BallActor,
		Data: pubsub.PositionNotificationPayload{
			XPos: b.xPos,
			YPos: b.yPos,
		},
	})

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

func (b *Ball) Destroy() {
	for _, subscriber := range b.subscribers {
		b.notificationBus.RemoveSubscriber(subscriber)
	}
}

func (b *Ball) IsActive() bool {
	return true
}

func (b *Ball) updatePosition(message *pubsub.Message) {
	switch message.GetMessageBody().ActorType {
	case pubsub.LeftBatActor:
		b.updatePositionOfLeftBat(message)
		break

	case pubsub.RightBatActor:
		b.updatePositionOfRightBat(message)
		break
	default:
		break // ignore other actor's positions
	}
}

func (b *Ball) updatePositionOfLeftBat(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		b.xPosLB = v.XPos
		b.yPosLB = v.YPos
	}
}

func (b *Ball) updatePositionOfRightBat(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		b.xPosRB = v.XPos
		b.yPosRB = v.YPos
	}
}

func moveBallAuto(b *Ball) {
	// Each frame, we move the ball in a series of small steps.
	// The number of steps being based on its speed attribute.
	for i := 1; i <= b.speed; i++ {
		moveBallAutoImpl(b)
	}
}

func moveBallManually(b *Ball) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		moveBallAutoImpl(b)
	}
}

func moveBallAutoImpl(b *Ball) {
	b.xPos += b.dx
	b.yPos += b.dy

	if b.hitBottom() || b.hitTop() {
		b.dy = -b.dy
	}

	if b.hitLeft() {
		b.dx = -b.dx
		b.notificationBus.Publish(pubsub.LEFT_BAT_MISS_NOTIFICATION_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.BallActor,
			Data:      nil,
		})
	}

	if b.hitRight() {
		b.dx = -b.dx
		b.notificationBus.Publish(pubsub.RIGHT_BAT_MISS_NOTIFICATION_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.BallActor,
			Data:      nil,
		})
	}

	if b.hitLeftBat() {
		b.dx = -b.dx                     // reverse x direction
		b.dy += b.deflectionForLeftBat() // deflect y direction based on which half of the bat did the ball hit

		b.notificationBus.Publish(pubsub.LEFT_BAT_HIT_NOTIFICATION_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.BallActor,
			Data:      nil,
		})

		b.notificationBus.Publish(pubsub.CREATE_IMPACT_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.BallActor,
			Data: pubsub.PositionNotificationPayload{
				XPos: b.xPos,
				YPos: b.yPos,
			},
		})

		b.xPos = BALL_MIN_X_BAT
		// Ensure our direction vector is a unit vector, i.e. represents a distance
		// of the equivalent of 1 pixel regardless of its angle.
		b.dx, b.dy = intoUnitVector(b.dx, b.dy)
		b.speed += 1
	}

	if b.hitRightBat() {
		b.dx = -b.dx                      // reverse x direction
		b.dy += b.deflectionForRightBat() // deflect y direction based on which half of the bat did the ball hit

		b.notificationBus.Publish(pubsub.RIGHT_BAT_HIT_NOTIFICATION_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.BallActor,
			Data:      nil,
		})

		b.notificationBus.Publish(pubsub.CREATE_IMPACT_TOPIC, pubsub.GameNotification{
			ActorType: pubsub.BallActor,
			Data: pubsub.PositionNotificationPayload{
				XPos: b.xPos,
				YPos: b.yPos,
			},
		})

		b.xPos = BALL_MAX_X_BAT
		b.dx, b.dy = intoUnitVector(b.dx, b.dy)
		b.speed += 1
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
	return calculateDeflection(b.yPos, b.yPosLB)
}

// calculates change in dy of right bat is hit based on where we did hit it (upper or lowe part)
func (b *Ball) deflectionForRightBat() float64 {
	return calculateDeflection(b.yPos, b.yPosRB)
}

func calculateDeflection(ballY float64, batY float64) float64 {

	batCenterY := batY + 16 + PAD_JPG_HEIGHT_HALF_PADONLY_PX
	ballCenterY := ballY + BALL_JPG_HEIGHT_TOTAL_PX/2

	diffY := ballCenterY - batCenterY

	deflection := diffY / PAD_JPG_HEIGHT_PADONLY_PX

	// Limit the deflection so we don't get into a situation
	// where the ball is bouncing up and down too rapidly.
	deflection = math.Min(math.Max(deflection, -1), 1)

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
