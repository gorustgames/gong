package game

import (
	"github.com/gorustgames/gong/game/actor"
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"log"
)

type Game struct {
	actors []actor.GameActor
}

const (
	SCREEN_WIDTH, SCREEN_HEIGHT = 800, 480
)

type GameStates int8

const (
	Menu GameStates = iota
	GameStartedSinglePlayer
	GameStartedMultiPlayer
	GameOver
)

var (
	GameState       GameStates
	game            Game
	notificationBus *pubsub.Broker
)

func init() {
	// GameState = Menu
	GameState = GameStartedSinglePlayer
}

// game state updates
func (g *Game) Update(_ *ebiten.Image) error {

	for _, actor := range g.actors {
		actor.Update()
	}

	return nil
}

// game rendering logic
func (g *Game) Draw(screen *ebiten.Image) {
	for _, actor := range g.actors {
		actor.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = SCREEN_WIDTH
	screenHeight = SCREEN_HEIGHT
	return
}

func singlePlayer(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsSinglePlayer(notificationBus)
}

func multiPlayer(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsMultiPlayer(notificationBus)
}

func menu(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsMenu(notificationBus)
}

func gameover(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsGameOver(notificationBus)
}

func destroyOldActors() {
	for _, actor := range game.actors {
		actor.Destroy()
	}
}

func CreateGameBus() {
	notificationBus = pubsub.NewBroker()

	subscriberMN := notificationBus.AddSubscriber()
	subscriberSP := notificationBus.AddSubscriber()
	subscriberMP := notificationBus.AddSubscriber()
	subscriberGO := notificationBus.AddSubscriber()

	notificationBus.Subscribe(subscriberMN, pubsub.CHANGE_GAME_STATE_MENU_TOPIC)
	notificationBus.Subscribe(subscriberSP, pubsub.CHANGE_GAME_STATE_SINGLE_PLAYER_TOPIC)
	notificationBus.Subscribe(subscriberMP, pubsub.CHANGE_GAME_STATE_MULTI_PLAYER_TOPIC)
	notificationBus.Subscribe(subscriberGO, pubsub.CHANGE_GAME_STATE_GAME_OVER_TOPIC)

	go subscriberMN.Listen(menu)
	go subscriberSP.Listen(singlePlayer)
	go subscriberMP.Listen(multiPlayer)
	go subscriberGO.Listen(gameover)

}

func StartGame() {
	actors := actor.CreateActorsMenu(notificationBus)
	// actors := actor.CreateActorsSinglePlayer(notificationBus)
	game = Game{
		actors: actors,
	}

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Go Pong")

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
