package game

import "errors"

const BoardSize = 19

type Stone int8

const (
	Empty Stone = 0
	Black Stone = 1
	White Stone = 2
)

type Point struct {
	X, Y int
}

var adj = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func inBounds(x, y int) bool {
	return x >= 0 && x < BoardSize && y >= 0 && y < BoardSize
}

type Board struct {
	cells [BoardSize][BoardSize]Stone
}

func NewBoard() *Board {
	return &Board{}
}

func (b *Board) Get(x, y int) Stone {
	return b.cells[x][y]
}

func (b *Board) Set(x, y int, s Stone) {
	b.cells[x][y] = s
}

func (b *Board) copy() *Board {
	nb := &Board{}
	nb.cells = b.cells
	return nb
}

// ToSlice returns the board as a column-major 2D slice for JSON.
func (b *Board) ToSlice() [][]Stone {
	out := make([][]Stone, BoardSize)
	for x := 0; x < BoardSize; x++ {
		out[x] = make([]Stone, BoardSize)
		for y := 0; y < BoardSize; y++ {
			out[x][y] = b.cells[x][y]
		}
	}
	return out
}

// group returns all stones connected to (x,y) of the same color.
func (b *Board) group(x, y int) []Point {
	color := b.cells[x][y]
	if color == Empty {
		return nil
	}
	visited := map[Point]bool{{x, y}: true}
	queue := []Point{{x, y}}
	for i := 0; i < len(queue); i++ {
		p := queue[i]
		for _, d := range adj {
			nx, ny := p.X+d[0], p.Y+d[1]
			np := Point{nx, ny}
			if inBounds(nx, ny) && !visited[np] && b.cells[nx][ny] == color {
				visited[np] = true
				queue = append(queue, np)
			}
		}
	}
	return queue
}

// liberties returns the number of distinct empty intersections adjacent to the group.
func (b *Board) liberties(grp []Point) int {
	seen := map[Point]bool{}
	for _, p := range grp {
		for _, d := range adj {
			nx, ny := p.X+d[0], p.Y+d[1]
			if inBounds(nx, ny) && b.cells[nx][ny] == Empty {
				seen[Point{nx, ny}] = true
			}
		}
	}
	return len(seen)
}

func opponent(s Stone) Stone {
	if s == Black {
		return White
	}
	return Black
}

// PlaceStone places a stone for player at (x,y), removes captured opponent groups,
// and rejects suicide. Returns the number of captured stones.
func (b *Board) PlaceStone(x, y int, player Stone) (captured int, err error) {
	if !inBounds(x, y) {
		return 0, errors.New("落子超出棋盘范围")
	}
	if b.cells[x][y] != Empty {
		return 0, errors.New("该位置已有棋子")
	}

	// Work on a copy so we can roll back on suicide without partial mutation.
	next := b.copy()
	next.cells[x][y] = player
	opp := opponent(player)

	for _, d := range adj {
		nx, ny := x+d[0], y+d[1]
		if !inBounds(nx, ny) || next.cells[nx][ny] != opp {
			continue
		}
		grp := next.group(nx, ny)
		if next.liberties(grp) == 0 {
			for _, p := range grp {
				next.cells[p.X][p.Y] = Empty
				captured++
			}
		}
	}

	if next.liberties(next.group(x, y)) == 0 {
		return 0, errors.New("禁止自杀落子")
	}

	b.cells = next.cells
	return captured, nil
}
