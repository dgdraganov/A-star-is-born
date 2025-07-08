package pathfind

import (
	"math"

	"github.com/dgdraganov/A-star-is-born/pkg/queue"
)

type CellState int

const (
	None CellState = 1 << iota
	Empty
	Visited
	Explored
	Obstacle
	Start
	End
	Path
)

func (cs CellState) IsOneOf(val CellState) bool {
	return val&cs == cs
}

const (
	diagonalCost = 14
	straightCost = 10
)

var (
	straights = [][2]int{
		{-1, 0}, // Up
		{1, 0},  // Down
		{0, -1}, // Left
		{0, 1},  // Right
	}

	diagonals = [][2]int{
		{-1, -1}, // Up-Left
		{-1, 1},  // Up-Right
		{1, -1},  // Down-Left
		{1, 1},   // Down-Right
	}
)

type Node struct {
	StartDist int
	EndDist   int
	Status    CellState
	Parent    *Node
	X         int
	Y         int
}

type Astar struct {
	field      [][]*Node
	pqExplored *queue.PriorityQueue[*Node]
	endNode    *Node
	startNode  *Node
}

func NewAstar() *Astar {
	pq := queue.NewPriorityQueue(func(a, b *Node) bool {
		return a.EndDist < b.EndDist
	})

	return &Astar{
		pqExplored: pq,
	}
}

func (a *Astar) Initialize(cells [][]CellState) {
	a.field = make([][]*Node, len(cells))
	for i := 0; i < len(cells); i++ {
		a.field[i] = make([]*Node, len(cells[0]))
		for j := 0; j < len(cells[0]); j++ {
			a.field[i][j] = &Node{
				StartDist: math.MaxInt,
				EndDist:   0,
				Status:    cells[i][j],
				Parent:    nil,
				X:         i,
				Y:         j,
			}

			switch cells[i][j] {
			case End:
				a.endNode = a.field[i][j]
			case Start:
				a.startNode = a.field[i][j]
			}
		}
	}

	a.startNode.StartDist = 0
	a.startNode.EndDist = a.getDistance(a.startNode.X, a.startNode.Y, a.endNode)
	a.pqExplored.Enqueue(a.startNode, (a.startNode.StartDist + a.startNode.EndDist))
}

func (a *Astar) Update() ([][]CellState, bool) {
	current := new(Node)
	var ok bool
	for current.Status.IsOneOf(Visited) {
		current, ok = a.pqExplored.Dequeue()
		if !ok {
			return nil, false
		}
	}

	if current.Status.IsOneOf(End) {
		stateMx := getStateMx(a.field)
		updatePathToEnd(stateMx, current)
		return stateMx, false
	}

	if current.Status.IsOneOf(Explored) {
		current.Status = Visited
	}

	a.exploreNeighbours(current, straights, straightCost)
	a.exploreNeighbours(current, diagonals, diagonalCost)
	stateMx := getStateMx(a.field)

	return stateMx, true
}

func (a *Astar) exploreNeighbours(node *Node, direction [][2]int, cost int) {
	for i := 0; i < len(direction); i++ {
		nX := node.X + direction[i][0]
		nY := node.Y + direction[i][1]
		if isOutsideField(a.field, nX, nY) || a.field[nX][nY].Status.IsOneOf(Obstacle|Start|Visited) {
			continue
		}

		neighbour := a.field[nX][nY]
		newStartDist := node.StartDist + cost

		if newStartDist >= neighbour.StartDist {
			continue
		}

		neighbour.StartDist = newStartDist
		neighbour.Parent = node

		if neighbour.EndDist == 0 {
			neighbour.EndDist = a.getDistance(nX, nY, a.endNode)
		}

		fScore := neighbour.StartDist + neighbour.EndDist

		if neighbour.Status.IsOneOf(Empty) {
			neighbour.Status = Explored
		}

		a.pqExplored.Enqueue(neighbour, fScore)
	}
}

func (a *Astar) getDistance(i int, j int, endNode *Node) int {
	diffX := math.Abs(float64(endNode.X - i))
	diffY := math.Abs(float64(endNode.Y - j))

	var res float64
	min := math.Min(diffX, diffY)
	res += min * diagonalCost
	res += math.Abs(diffX-diffY) * straightCost
	return int(res)
}

func getStateMx(node [][]*Node) [][]CellState {
	res := make([][]CellState, len(node))
	for i := 0; i < len(node); i++ {
		res[i] = make([]CellState, len(node[i]))
		for j := 0; j < len(node[i]); j++ {
			res[i][j] = node[i][j].Status
		}
	}
	return res
}

func isOutsideField(field [][]*Node, i, j int) bool {
	if i < 0 || i >= len(field) || j < 0 || j >= len(field[0]) {
		return true
	}
	return false
}

func updatePathToEnd(stateMx [][]CellState, current *Node) {
	child := current.Parent
	for child.Status != Start {
		stateMx[child.X][child.Y] = Path
		child = child.Parent
	}
}
