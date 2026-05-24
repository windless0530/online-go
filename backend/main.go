package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"online-go/game"
	"online-go/server"
)

func main() {
	statePath, err := resolveStatePath()
	if err != nil {
		log.Fatalf("resolve state path: %v", err)
	}

	g := loadOrNewGame(statePath)
	h := server.New(g, statePath)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// In production, serve the built frontend from ../frontend/dist.
	// In dev, the Vite dev server proxies /api to this process.
	if _, err := os.Stat("../frontend/dist"); err == nil {
		mux.Handle("/", http.FileServer(http.Dir("../frontend/dist")))
	}

	addr := ":8080"
	log.Printf("backend listening on http://localhost%s", addr)
	log.Printf("state persisted at %s", statePath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// resolveStatePath returns the path to the persistence file, creating the
// parent directory if needed. Honors $ONLINE_GO_STATE for tests/overrides;
// defaults to ./data/state.json relative to the backend binary's CWD.
func resolveStatePath() (string, error) {
	path := os.Getenv("ONLINE_GO_STATE")
	if path == "" {
		path = filepath.Join("data", "state.json")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	return path, nil
}

func loadOrNewGame(path string) *game.Game {
	s, err := game.LoadState(path)
	if err != nil {
		if !errors.Is(err, game.ErrNoState) {
			log.Printf("load state failed, starting fresh: %v", err)
		}
		return game.NewGame()
	}
	log.Printf("restored state from %s", path)
	return game.NewGameFromState(s)
}