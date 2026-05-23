package game

import "sync"

type GameState struct {
	Board         [][]Stone `json:"board"`
	CurrentPlayer string    `json:"currentPlayer"`
	BlackCaptures int       `json:"blackCaptures"`
	WhiteCaptures int       `json:"whiteCaptures"`
}

type Game struct {
	mu            sync.Mutex
	board         *Board
	currentPlayer Stone
	blackCaptures int
	whiteCaptures int
}

func NewGame() *Game {
	return &Game{
		board:         NewBoard(),
		currentPlayer: Black,
	}
}

func (g *Game) State() GameState {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.snapshot()
}

func (g *Game) PlaceStone(x, y int) (GameState, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	captured, err := g.board.PlaceStone(x, y, g.currentPlayer)
	if err != nil {
		return GameState{}, err
	}

	if g.currentPlayer == Black {
		g.blackCaptures += captured
		g.currentPlayer = White
	} else {
		g.whiteCaptures += captured
		g.currentPlayer = Black
	}

	return g.snapshot(), nil
}

func (g *Game) Reset() GameState {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.board = NewBoard()
	g.currentPlayer = Black
	g.blackCaptures = 0
	g.whiteCaptures = 0
	return g.snapshot()
}

func (g *Game) snapshot() GameState {
	return GameState{
		Board:         g.board.ToSlice(),
		CurrentPlayer: stoneToPlayer(g.currentPlayer),
		BlackCaptures: g.blackCaptures,
		WhiteCaptures: g.whiteCaptures,
	}
}

func stoneToPlayer(s Stone) string {
	if s == Black {
		return "black"
	}
	return "white"
}
