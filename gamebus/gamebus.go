package gamebus

type GameNotificationType int8

const (
	PositionNotification GameNotificationType = iota
	LeftBatHitNotification
	RightBatHitNotification
)

type GameNotificationActorType int8

const (
	LeftBatActor GameNotificationActorType = iota
	RightBatActor
	BallActor
)

type PositionNotificationPayload struct {
	XPos float64
	YPos float64
}

type GameNotification struct {
	ActorType            GameNotificationActorType
	GameNotificationType GameNotificationType
	Data                 any
}

type GameNotificationBus struct {
	Bus chan GameNotification
}

func NewGameNotificationBus() *GameNotificationBus {
	return &GameNotificationBus{
		Bus: make(chan GameNotification),
	}
}
