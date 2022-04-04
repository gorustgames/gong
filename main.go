package main

import (
	"github.com/gorustgames/gong/game"
	"github.com/hajimehoshi/ebiten"
	"log"
)

func main() {
	ebiten.SetWindowSize(game.SCREEN_WIDTH, game.SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Go Pong")
	if err := ebiten.RunGame(game.CreateGame()); err != nil {
		log.Fatal(err)
	}
}
