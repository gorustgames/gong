package pubsub

type GameNotificationActorType int8

const (
	LeftBatActor GameNotificationActorType = iota
	RightBatActor
	BallActor
	MenuActor
	GameBoardActor
	GameOverActor
)

const (
	POSITION_NOTIFICATION_TOPIC = "POSITION_NOTIFICATION_TOPIC"

	LEFT_BAT_HIT_NOTIFICATION_TOPIC  = "LEFT_BAT_HIT_NOTIFICATION_TOPIC"
	RIGHT_BAT_HIT_NOTIFICATION_TOPIC = "RIGHT_BAT_HIT_NOTIFICATION_TOPIC"

	LEFT_BAT_MISS_NOTIFICATION_TOPIC  = "LEFT_BAT_MISS_NOTIFICATION_TOPIC"
	RIGHT_BAT_MISS_NOTIFICATION_TOPIC = "RIGHT_BAT_MISS_NOTIFICATION_TOPIC"

	CHANGE_GAME_STATE_SINGLE_PLAYER_TOPIC = "CHANGE_GAME_STATE_SINGLE_PLAYER_TOPIC"
	CHANGE_GAME_STATE_MULTI_PLAYER_TOPIC  = "CHANGE_GAME_STATE_MULTI_PLAYER_TOPIC"
	CHANGE_GAME_STATE_GAME_OVER_TOPIC     = "CHANGE_GAME_STATE_GAME_OVER_TOPIC"
	CHANGE_GAME_STATE_MENU_TOPIC          = "CHANGE_GAME_STATE_MENU_TOPIC"

	CREATE_IMPACT_TOPIC = "CREATE_IMPACT_TOPIC"
)

type PositionNotificationPayload struct {
	XPos float64
	YPos float64
}

type GameOverNotificationPayload struct {
	ScoreLeft  int
	ScoreRight int
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
