package ui

import (
	"image/color"

	"github.com/dgdraganov/A-star-is-born/internal/pathfind"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	startX       = 28
	startY       = 1
	endX         = 1
	endY         = 28
	gridSize     = 10
	layOutSizeX  = 301
	layOutSizeY  = 301
	buttonHeight = 30
	windowSizeX  = 600
	windowSizeY  = 600
)

var (
	btnX      = float64(layOutSizeX/2 - btnWidth/2)
	btnY      = float64(layOutSizeY) + 5
	btnWidth  = gridSize * 5
	btnHeight = gridSize * 2
)

var (
	colorMap = map[pathfind.CellState]color.Color{
		pathfind.Empty:    emptyColor,
		pathfind.Visited:  visitedColor,
		pathfind.Explored: exploredColor,
		pathfind.Obstacle: obstacleColor,
		pathfind.Start:    startColor,
		pathfind.End:      endColor,
		pathfind.Path:     pathColor,
	}

	buttonColor   = color.RGBA{77, 77, 77, 133}
	startColor    = color.RGBA{0, 255, 255, 0}
	endColor      = color.RGBA{255, 128, 0, 0}
	emptyColor    = color.RGBA{0, 0, 0, 0}
	obstacleColor = color.RGBA{255, 255, 199, 0}
	visitedColor  = color.RGBA{255, 0, 0, 0}
	exploredColor = color.RGBA{0, 255, 0, 0}
	pathColor     = color.RGBA{0, 0, 255, 0}
)

type Cell struct {
	State pathfind.CellState
}

type UI struct {
	field       [][]Cell
	astar       Pathfinder
	square      *ebiten.Image
	gameEnded   bool
	gameStarted bool
	isDrawing   bool
}

func NewGame(pf *pathfind.Astar) *UI {
	ebiten.SetWindowSize(windowSizeX, windowSizeY)
	ebiten.SetWindowTitle("A star is born")

	rect := ebiten.NewImage(gridSize-1, gridSize-1)

	ui := &UI{
		gameStarted: false,
		field:       initializeField(),
		astar:       pf,
		square:      rect,
	}

	return ui
}

func (g *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return layOutSizeX, layOutSizeY + buttonHeight
}

func (g *UI) Update() error {
	if !g.gameStarted {
		g.drawingPhase()
		return nil
	}

	if !g.gameEnded {
		updatedMx, ok := g.astar.Update()
		if !ok {
			g.gameEnded = true
		}

		g.updateGameField(updatedMx)
	}

	return nil
}

func (g *UI) Draw(screen *ebiten.Image) {
	g.drawGridLines(screen)
	g.drawField(screen)
	if !g.gameStarted {
		g.drawStartButton(screen)
	}
}

// =================================== PRIVATE METHODS =========================================

func (g *UI) drawField(screen *ebiten.Image) {
	for i := 0; i < len(g.field); i++ {
		for j := 0; j < len(g.field[i]); j++ {
			g.drawSquare(screen, j, i, colorMap[g.field[i][j].State])
		}
	}
}

func (g *UI) drawStartButton(screen *ebiten.Image) {
	button := ebiten.NewImage(btnWidth, btnHeight)
	button.Fill(buttonColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(btnX, btnY)
	screen.DrawImage(button, op)

	ebitenutil.DebugPrintAt(screen, " Start", int(btnX), int(btnY))
}

func (g *UI) drawSquare(screen *ebiten.Image, x, y int, col color.Color) {
	g.square.Fill(col)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(1+x*gridSize), float64(1+y*gridSize))
	screen.DrawImage(g.square, op)
}

func (g *UI) drawGridLines(screen *ebiten.Image) {
	currentColumn := 1
	for currentColumn <= layOutSizeY {
		vector.StrokeLine(
			screen,
			float32(0), float32(currentColumn),
			float32(layOutSizeX), float32(currentColumn),
			1,
			color.Gray{Y: 125},
			false,
		)
		currentColumn += gridSize
	}

	currentRow := 1
	for currentRow <= layOutSizeX {
		vector.StrokeLine(
			screen,
			float32(currentRow), float32(0),
			float32(currentRow), float32(layOutSizeY),
			1,
			color.Gray{Y: 125},
			false,
		)
		currentRow += gridSize
	}
}

func (g *UI) drawingPhase() {
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.isDrawing = false
		return
	}

	x, y := ebiten.CursorPosition()
	gridX := (x - 1) / gridSize
	gridY := (y - 1) / gridSize
	if g.isValidCell(gridX, gridY) {
		g.isDrawing = true
		g.field[gridY][gridX].State = pathfind.Obstacle
	} else if !g.isDrawing && isStartBtnHover(x, y) {
		g.gameStarted = true
		cellStateMx := getCellStateMx(g.field)
		g.astar.Initialize(cellStateMx)
	}
}

func (g *UI) updateGameField(updatedMx [][]pathfind.CellState) {
	for i := 0; i < len(updatedMx); i++ {
		for j := 0; j < len(updatedMx[i]); j++ {
			g.field[i][j].State = updatedMx[i][j]
		}
	}
}

func (g *UI) isValidCell(gridX, gridY int) bool {
	if isOutsideField(gridX, gridY) {
		return false
	}
	if g.field[gridY][gridX].State != pathfind.Empty {
		return false
	}
	return true
}

func isOutsideField(x, y int) bool {
	if x < 0 || x >= layOutSizeX/gridSize || y < 0 || y >= layOutSizeY/gridSize {
		return true
	}
	return false
}

func isStartBtnHover(x, y int) bool {
	if x >= int(btnX) && x <= int(btnX)+btnWidth && y >= int(btnY) && y <= int(btnY)+btnHeight {
		return true
	}
	return false
}

func getCellStateMx(cell [][]Cell) [][]pathfind.CellState {
	res := make([][]pathfind.CellState, len(cell))
	for i := 0; i < len(cell); i++ {
		res[i] = make([]pathfind.CellState, len(cell[i]))
		for j := 0; j < len(cell[i]); j++ {
			res[i][j] = cell[i][j].State
		}
	}
	return res
}

func initializeField() [][]Cell {
	fieldArr := make([][]Cell, layOutSizeY/gridSize)

	for i := range fieldArr {
		fieldArr[i] = make([]Cell, layOutSizeX/gridSize)
		for j := range fieldArr[i] {
			fieldArr[i][j] = Cell{
				State: pathfind.Empty,
			}
		}
	}

	fieldArr[startY][startX].State = pathfind.Start
	fieldArr[endY][endX].State = pathfind.End
	return fieldArr
}
