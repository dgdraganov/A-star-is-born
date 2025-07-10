package main

import (
	"log"

	"github.com/dgdraganov/A-star-is-born/internal/pathfind"
	"github.com/dgdraganov/A-star-is-born/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	astar := pathfind.NewAstar()
	game := ui.NewGame(astar)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("game running error:", err)
	}
}
