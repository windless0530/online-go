package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"online-go/game"
)

type Handler struct {
	game        *game.Game
	persistPath string
}

func New(g *game.Game, persistPath string) *Handler {
	return &Handler{game: g, persistPath: persistPath}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/state", h.getState)
	mux.HandleFunc("POST /api/move", h.postMove)
	mux.HandleFunc("POST /api/reset", h.postReset)
}

func (h *Handler) getState(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.game.State())
}

type moveRequest struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"` // "black" or "white"
}

type moveResponse struct {
	State *game.GameState `json:"state"`
	Error string          `json:"error,omitempty"`
}

func (h *Handler) postMove(w http.ResponseWriter, r *http.Request) {
	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, moveResponse{Error: "请求格式错误"})
		return
	}

	player, err := parseColor(req.Color)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, moveResponse{Error: err.Error()})
		return
	}

	state, err := h.game.PlaceStone(req.X, req.Y, player)
	if err != nil {
		writeJSON(w, http.StatusOK, moveResponse{Error: err.Error()})
		return
	}
	h.persist(state)
	writeJSON(w, http.StatusOK, moveResponse{State: &state})
}

func (h *Handler) postReset(w http.ResponseWriter, _ *http.Request) {
	state := h.game.Reset()
	h.persist(state)
	writeJSON(w, http.StatusOK, state)
}

// persist saves the snapshot synchronously. A failure is logged but does not
// fail the request — the user's move already succeeded in memory, and the next
// successful save will overwrite the lost one.
func (h *Handler) persist(s game.GameState) {
	if h.persistPath == "" {
		return
	}
	if err := game.SaveState(h.persistPath, s); err != nil {
		log.Printf("persist state failed: %v", err)
	}
}

func parseColor(s string) (game.Stone, error) {
	switch s {
	case "black":
		return game.Black, nil
	case "white":
		return game.White, nil
	default:
		return game.Empty, fmt.Errorf("无效颜色: %q，需传 \"black\" 或 \"white\"", s)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
	}
}