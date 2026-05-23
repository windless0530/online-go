package main

import (
	"log"
	"net/http"
	"os"

	"online-go/game"
	"online-go/server"
)

func main() {
	g := game.NewGame()
	h := server.New(g)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// In production, serve the built frontend from ../frontend/dist.
	// In dev, the Vite dev server proxies /api to this process.
	if _, err := os.Stat("../frontend/dist"); err == nil {
		mux.Handle("/", http.FileServer(http.Dir("../frontend/dist")))
	}

	addr := ":8080"
	log.Printf("backend listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
