package pathfind

import "github.com/dgdraganov/A-star-is-born/pkg/queue"

type CellState int

const (
	Empty CellState = iota
	Visited
	Explored
	Obstacle
	Start
	End
	Path
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
	return &Astar{}
}

func (a *Astar) Initialize(cells [][]CellState) {

	a.field = make([][]*Node, len(cells))
	for i := 0; i < len(cells); i++ {
		a.field[i] = make([]*Node, len(cells[0]))
		for j := 0; j < len(cells[0]); j++ {
			a.field[i][j] = &Node{
				StartDist: -1,
				EndDist:   -1,
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

	// a.updateDistancesToEnd(a.endNode)
	// dist := a.getDistance(a.startNode.X, a.startNode.Y, a.endNode)
	// a.startNode.EndDist = dist
	// a.startNode.StartDist = 0
	a.pqExplored.Enqueue(a.startNode, -(a.startNode.StartDist + a.startNode.EndDist))
}
