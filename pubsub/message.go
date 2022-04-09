package pubsub

type GameNotificationActorType int8

const (
	LeftBatActor GameNotificationActorType = iota
	RightBatActor
	BallActor
)

const (
	POSITION_NOTIFICATION_TOPIC      = "PositionNotificationTopic"
	LEFT_BAT_HIT_NOTIFICATION_TOPIC  = "LeftBatHitNotificationTopic"
	RIGHT_BAT_HIT_NOTIFICATION_TOPIC = "RightBatHitNotificationTopic"
)

type PositionNotificationPayload struct {
	XPos float64
	YPos float64
}

type GameNotification struct {
	ActorType GameNotificationActorType
	Data      any
}

type Message struct {
	topic string
	body  GameNotification
}

func NewMessage(msg GameNotification, topic string) *Message {
	// Returns the message object
	return &Message{
		topic: topic,
		body:  msg,
	}
}
func (m *Message) GetTopic() string {
	// returns the topic of the message
	return m.topic
}
func (m *Message) GetMessageBody() GameNotification {
	// returns the message body.
	return m.body
}
