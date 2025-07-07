package ui

import "github.com/dgdraganov/A-star-is-born/internal/pathfind"

type Pathfinder interface {
	Initialize(cells [][]pathfind.CellState)
	Update() ([][]pathfind.CellState, bool)
}
