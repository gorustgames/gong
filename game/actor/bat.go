package actor

import (
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math"
	"math/rand"
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
	MAX_AI_SPEED = 6
	HALF_WIDTH   = 400   //SCREEN_WIDTH / 2
	TARGET_BAT_Y = 168.0 //y pos when bat is ~ in the middle of the playing area
)

type Bat struct {
	batImage        *ebiten.Image
	batHitImage     *ebiten.Image
	batEffectImage  *ebiten.Image
	xPos            float64
	yPos            float64
	dx              float64
	dy              float64
	xPosBall        float64
	yPosBall        float64
	aiOffset        float64
	speed           float64
	playerLocation  PlayerLocation
	playerType      PlayerType
	notificationBus *pubsub.Broker
	showHitCounter  int
	showMissCounter int
	subscribers     []*pubsub.Subscriber
}

func NewBat(playerLocation PlayerLocation, playerType PlayerType, notificationBus *pubsub.Broker) *Bat {

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

	fileName = "assets/effect0.png"
	if playerLocation == RightPlayer {
		fileName = "assets/effect1.png"
	}
	_batEffectImage, _, err := ebitenutil.NewImageFromFile(fileName, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	_xPos := -40.0
	if playerLocation == RightPlayer {
		_xPos = 680.0
	}

	newBat := &Bat{
		xPos:            _xPos,
		yPos:            0,
		dx:              0,
		dy:              0,
		xPosBall:        0,
		yPosBall:        0,
		aiOffset:        0,
		speed:           PLAYER_SPEED,
		batImage:        _batImage,
		batHitImage:     _batHitImage,
		batEffectImage:  _batEffectImage,
		playerLocation:  playerLocation,
		playerType:      playerType,
		notificationBus: notificationBus,
		showHitCounter:  0,
		showMissCounter: 0,
	}

	subscriberBatHit := notificationBus.AddSubscriber("bat-subscriberBatHit")
	if playerLocation == LeftPlayer {
		notificationBus.Subscribe(subscriberBatHit, pubsub.LEFT_BAT_HIT_NOTIFICATION_TOPIC)
	} else {
		notificationBus.Subscribe(subscriberBatHit, pubsub.RIGHT_BAT_HIT_NOTIFICATION_TOPIC)
	}
	go subscriberBatHit.Listen(newBat.initHitCounter)

	subscriberBatMiss := notificationBus.AddSubscriber("bat-subscriberBatMiss")
	if playerLocation == LeftPlayer {
		notificationBus.Subscribe(subscriberBatMiss, pubsub.LEFT_BAT_MISS_NOTIFICATION_TOPIC)
	} else {
		notificationBus.Subscribe(subscriberBatMiss, pubsub.RIGHT_BAT_MISS_NOTIFICATION_TOPIC)
	}
	go subscriberBatMiss.Listen(newBat.initMissCounter)

	subscriberBallPos := notificationBus.AddSubscriber("bat-subscriberBallPos")
	notificationBus.Subscribe(subscriberBallPos, pubsub.POSITION_NOTIFICATION_TOPIC)
	go subscriberBallPos.Listen(newBat.updateBallPosition)

	newBat.subscribers = make([]*pubsub.Subscriber, 3)
	newBat.subscribers[0] = subscriberBatHit
	newBat.subscribers[1] = subscriberBatMiss
	newBat.subscribers[2] = subscriberBallPos

	return newBat
}

func (b *Bat) Update() error {
	movePlayer(b)

	gameNotificationActorType := pubsub.LeftBatActor

	if b.playerLocation == RightPlayer {
		gameNotificationActorType = pubsub.RightBatActor
	}

	b.notificationBus.Publish(pubsub.POSITION_NOTIFICATION_TOPIC, pubsub.GameNotification{
		ActorType: gameNotificationActorType,
		Data: pubsub.PositionNotificationPayload{
			XPos: b.xPos,
			YPos: b.yPos,
		},
	})

	if b.showHitCounter > 0 {
		b.showHitCounter -= 1
	}

	if b.showMissCounter > 0 {
		b.showMissCounter -= 1
	}

	return nil
}

func (b *Bat) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.xPos, b.yPos)

	if b.showHitCounter > 0 {
		screen.DrawImage(b.batHitImage, op)
	} else if b.showMissCounter > 0 {
		opEffect := &ebiten.DrawImageOptions{}
		opEffect.GeoM.Translate(0, 0)
		screen.DrawImage(b.batEffectImage, opEffect)
	} else {
		screen.DrawImage(b.batImage, op)
	}
}

func (b *Bat) Id() string {
	if b.playerLocation == LeftPlayer {
		return "actor-left-bat"
	} else {
		return "actor-right-bat"
	}
}

func (b *Bat) Destroy() {
	for _, subscriber := range b.subscribers {
		b.notificationBus.RemoveSubscriber(subscriber)
	}
}

func (b *Bat) initHitCounter(_ *pubsub.Message) {
	b.showHitCounter = 20
	b.aiOffset = randInRange(-10, 10)
}

func (b *Bat) initMissCounter(_ *pubsub.Message) {
	b.showMissCounter = 20
}

func (b *Bat) updateBallPosition(message *pubsub.Message) {
	switch message.GetMessageBody().ActorType {
	case pubsub.BallActor:
		switch v := message.GetMessageBody().Data.(type) {
		case pubsub.PositionNotificationPayload:
			b.xPosBall = v.XPos
			b.yPosBall = v.YPos
		}
	}
}

func randInRange(min int, max int) float64 {
	return float64(rand.Intn(max-min) + min)
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
	} else /* AI */ {
		b.yPos += aiDy(b.xPosBall, b.yPosBall, b.xPos, b.yPos, b.aiOffset)
	}

	// make sure we don't cross up/down wall
	if b.yPos < 0 {
		b.yPos = 0
	}

	if b.yPos > 320 {
		b.yPos = 320
	}
}

func aiDy(xPosBall float64, yPosBall float64, xPos float64, yPos float64, aiOffset float64) float64 {
	xDistance := xPos - xPosBall

	// If the ball is far away, we move towards the centre of the screen (HALF_HEIGHT), on the basis that we don't
	// yet know whether the ball will be in the top or bottom half of the screen when it reaches our position on
	// the X axis. By waiting at a central position, we're as ready as it's possible to be for all eventualities.
	targetY1 := TARGET_BAT_Y
	// deduct 16(upper margin) + 64 (half of the bat width), since we use top base Y, not middle based Y
	// also deduct number offset to make computer more human like aka error-prone
	targetY2 := yPosBall + float64(aiOffset) - 16 - 64

	/*
		The final step is to work out the actual Y position we want to move towards. We use what's called a weighted
		average - taking the average of the two target Y positions we've previously calculated, but shifting the
		balance towards one or the other depending on how far away the ball is. If the ball is more than 400 pixels
		(half the screen width) away on the X axis, our target will be half the screen height (target_y_1). If the
		ball is at the same position as us on the X axis, our target will be target_y_2. If it's 200 pixels away,
		we'll aim for halfway between target_y_1 and target_y_2. This reflects the idea that as the ball gets closer,
		we have a better idea of where it's going to end up.
	*/
	weight1 := math.Min(1, xDistance/HALF_WIDTH)
	weight2 := 1 - weight1
	targetY := (weight1 * targetY1) + (weight2 * targetY2)

	// make sure we do not move faster than allowed (MAX_AI_SPEED) regardless of the direction we are moving to.
	targetDy := math.Min(MAX_AI_SPEED, math.Max(-MAX_AI_SPEED, targetY-yPos))
	return targetDy
}
