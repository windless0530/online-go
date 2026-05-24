package game

import "sync"

type GameState struct {
	Board         [][]Stone `json:"board"`
	BlackCaptures int       `json:"blackCaptures"`
	WhiteCaptures int       `json:"whiteCaptures"`
}

type Game struct {
	mu            sync.Mutex
	board         *Board
	blackCaptures int
	whiteCaptures int
}

func NewGame() *Game {
	return &Game{board: NewBoard()}
}

// NewGameFromState reconstructs a Game from a persisted snapshot.
// The caller is responsible for validating the state first (see LoadState).
func NewGameFromState(s GameState) *Game {
	b := NewBoard()
	for x := 0; x < BoardSize; x++ {
		for y := 0; y < BoardSize; y++ {
			b.Set(x, y, s.Board[x][y])
		}
	}
	return &Game{
		board:         b,
		blackCaptures: s.BlackCaptures,
		whiteCaptures: s.WhiteCaptures,
	}
}

func (g *Game) State() GameState {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.snapshot()
}

// PlaceStone places a stone of the given color. Turn management is the caller's responsibility.
func (g *Game) PlaceStone(x, y int, player Stone) (GameState, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	captured, err := g.board.PlaceStone(x, y, player)
	if err != nil {
		return GameState{}, err
	}

	if player == Black {
		g.blackCaptures += captured
	} else {
		g.whiteCaptures += captured
	}

	return g.snapshot(), nil
}

func (g *Game) Reset() GameState {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.board = NewBoard()
	g.blackCaptures = 0
	g.whiteCaptures = 0
	return g.snapshot()
}

func (g *Game) snapshot() GameState {
	return GameState{
		Board:         g.board.ToSlice(),
		BlackCaptures: g.blackCaptures,
		WhiteCaptures: g.whiteCaptures,
	}
}
