package game

import (
	"github.com/gorustgames/gong/game/actor"
	"github.com/gorustgames/gong/pubsub"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"log"
	"os"
	"time"
)

type Game struct {
	actors      []actor.GameActor
	audioPlayer *audio.Player
}

const (
	SCREEN_WIDTH, SCREEN_HEIGHT = 800, 480
	SAMPLE_RATE                 = 44100
	MUSIC_LEN_SEC               = 64
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

	for idx, actor := range g.actors {
		if actor.IsActive() /* update active actor*/ {
			actor.Update()
		} else /* remove inactive actor*/ {
			g.YankActor(idx)
		}
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

// YankActor removes actor at position specified by index parameter
// probably not extremely efficient (might be issue with large number of actors)
// but here it should work just fine
func (g *Game) YankActor(index int) {
	// log.Printf("YankActor: %d", index)
	newActors := append(g.actors[:index], g.actors[index+1:]...)
	g.actors = newActors
}

func transitionToSinglePlayerCallback(_ *pubsub.Message) {
	log.Printf("showing game SP")
	destroyOldActors()
	game.actors = actor.CreateActorsSinglePlayer(notificationBus)
	enableRendering()
}

func transitionToMultiPlayerCallback(_ *pubsub.Message) {
	log.Printf("showing game MP")
	destroyOldActors()
	game.actors = actor.CreateActorsMultiPlayer(notificationBus)
	enableRendering()
}

func transitionToMenuCallback(_ *pubsub.Message) {
	log.Printf("showing menu")
	destroyOldActors()
	// sleep for 0.5 seconds before showing menu, otherwise it
	// is happening that space key hit will be also captured by menu
	// actor and it will proceed directly to single player game without
	// really waiting for user choice.
	time.Sleep(500 * time.Millisecond)
	game.actors = actor.CreateActorsMenu(notificationBus)
	enableRendering()
}

func transitionToGameoverCallback(_ *pubsub.Message) {
	destroyOldActors()
	game.actors = actor.CreateActorsGameOver(notificationBus)
	enableRendering()
}

func createImpactCallback(message *pubsub.Message) {
	switch v := message.GetMessageBody().Data.(type) {
	case pubsub.PositionNotificationPayload:
		impactActor := actor.NewImpact(v.XPos, v.YPos, notificationBus)
		game.actors = append(game.actors, impactActor)
	}
}

func destroyOldActors() {
	disableRendering()
	for _, actor := range game.actors {
		actor.Destroy()
	}

	game.actors = nil // make GC to remove old actors
}

func disableRendering() {
	// disable state update & rendering loop
	changingGameStateInProgress = true
}

func enableRendering() {
	// enable state update & rendering loop
	changingGameStateInProgress = false
}

func createGameBus() {
	notificationBus = pubsub.NewBroker()

	subscriberMN := notificationBus.AddSubscriber("subscriberMN")
	subscriberSP := notificationBus.AddSubscriber("subscriberSP")
	subscriberMP := notificationBus.AddSubscriber("subscriberMP")
	subscriberGO := notificationBus.AddSubscriber("subscriberGO")
	subscriberIM := notificationBus.AddSubscriber("subscriberIM")

	notificationBus.Subscribe(subscriberMN, pubsub.CHANGE_GAME_STATE_MENU_TOPIC)
	notificationBus.Subscribe(subscriberSP, pubsub.CHANGE_GAME_STATE_SINGLE_PLAYER_TOPIC)
	notificationBus.Subscribe(subscriberMP, pubsub.CHANGE_GAME_STATE_MULTI_PLAYER_TOPIC)
	notificationBus.Subscribe(subscriberGO, pubsub.CHANGE_GAME_STATE_GAME_OVER_TOPIC)
	notificationBus.Subscribe(subscriberIM, pubsub.CREATE_IMPACT_TOPIC)

	go subscriberMN.Listen(transitionToMenuCallback)
	go subscriberSP.Listen(transitionToSinglePlayerCallback)
	go subscriberMP.Listen(transitionToMultiPlayerCallback)
	go subscriberGO.Listen(transitionToGameoverCallback)
	go subscriberIM.Listen(createImpactCallback)

}

func prepareAudioPlayerMusic() *audio.Player {
	audioContext := audio.NewContext(SAMPLE_RATE)
	f, err := os.Open("assets/sounds/theme.ogg")
	if err != nil {
		log.Fatal(err)
	}

	data, err := vorbis.DecodeWithSampleRate(SAMPLE_RATE, f)
	if err != nil {
		log.Fatal(err)
	}

	audioPlayer, err := audioContext.NewPlayer(data)
	if err != nil {
		log.Fatal(err)
	}

	return audioPlayer
}

// prepareAudioPlayerMusicInfinite creates infinite music loop
// https://programmer.ink/think/ebiten-learning-infinite-loop-player.html
func prepareAudioPlayerMusicInfinite() *audio.Player {
	audioContext := audio.NewContext(SAMPLE_RATE)
	f, err := os.Open("assets/sounds/theme.ogg")
	if err != nil {
		log.Fatal(err)
	}

	data, err := vorbis.DecodeWithSampleRate(SAMPLE_RATE, f)
	if err != nil {
		log.Fatal(err)
	}

	s := audio.NewInfiniteLoopWithIntro(data, 0, MUSIC_LEN_SEC*4*SAMPLE_RATE)

	audioPlayer, err := audioContext.NewPlayer(s)
	if err != nil {
		log.Fatal(err)
	}
	audioPlayer.SetVolume(0.3)
	return audioPlayer
}

func StartGame() {
	createGameBus()

	actors := actor.CreateActorsMenu(notificationBus)

	audioPlayer := prepareAudioPlayerMusicInfinite()

	game = Game{
		actors:      actors,
		audioPlayer: audioPlayer,
	}

	game.audioPlayer.Rewind()
	game.audioPlayer.Play()

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Go Pong")

	changingGameStateInProgress = false

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
