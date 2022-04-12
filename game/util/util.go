package util

import (
	"crypto/rand"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"log"
	"os"
)

const (
	SAMPLE_RATE   = 44100
	MUSIC_LEN_SEC = 64
)

var (
	audioContext *audio.Context
)

func init() {
	audioContext = audio.NewContext(SAMPLE_RATE)
}

// NewAudioPlayerMusicInfinite creates infinite music loop
// https://programmer.ink/think/ebiten-learning-infinite-loop-player.html
func NewAudioPlayerMusicInfinite() *audio.Player {
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

// NewAudioPlayer returns configured audio player for given asset
func NewAudioPlayer(asset string) *audio.Player {
	f, err := os.Open(fmt.Sprintf("assets/sounds/%s.ogg", asset))
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

// GenerateShortId generates short and (hopefully) uniqueue ID
func GenerateShortId() string {
	// alternativelly we could use https://github.com/teris-io/shortid
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	id := fmt.Sprintf("%X-%X", b[0:4], b[4:8])
	return id
}
