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
