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

var (
	game                        Game
	notificationBus             *pubsub.Broker
	changingGameStateInProgress bool
)

func init() {
	return
}

// game state updates
func (g *Game) Update(_ *ebiten.Image) error {

	if changingGameStateInProgress {
		return nil
	}

	for _, actor := range g.actors {
		actor.Update()
	}

	return nil
}

// game rendering logic
func (g *Game) Draw(screen *ebiten.Image) {

	if changingGameStateInProgress {
		return
	}

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
	changingGameStateInProgress = false
	enableRendering()
}

func multiPlayer(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsMultiPlayer(notificationBus)
	enableRendering()
}

func menu(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsMenu(notificationBus)
	enableRendering()
}

func gameover(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsGameOver(notificationBus)
	enableRendering()
}

func destroyOldActors() {
	disableRendering()
	for _, actor := range game.actors {
		actor.Destroy()
	}
}

func disableRendering() {
	changingGameStateInProgress = true // disable state update & rendering loop
}

func enableRendering() {
	changingGameStateInProgress = false
}

func createGameBus() {
	notificationBus = pubsub.NewBroker()

	subscriberMN := notificationBus.AddSubscriber("subscriberMN")
	subscriberSP := notificationBus.AddSubscriber("subscriberSP")
	subscriberMP := notificationBus.AddSubscriber("subscriberMP")
	subscriberGO := notificationBus.AddSubscriber("subscriberGO")

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
	createGameBus()

	actors := actor.CreateActorsMenu(notificationBus)
	//actors := actor.CreateActorsSinglePlayer(notificationBus)
	game = Game{
		actors: actors,
	}

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Go Pong")

	changingGameStateInProgress = false

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
