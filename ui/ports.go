package ui

import "github.com/dgdraganov/A-star-is-born/pathfind"

type Pathfinder interface {
	Initialize(cells [][]pathfind.CellState)
}
