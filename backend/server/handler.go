package server

import (
	"encoding/json"
	"net/http"

	"online-go/game"
)

type Handler struct {
	game *game.Game
}

func New(g *game.Game) *Handler {
	return &Handler{game: g}
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
	X int `json:"x"`
	Y int `json:"y"`
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

	state, err := h.game.PlaceStone(req.X, req.Y)
	if err != nil {
		writeJSON(w, http.StatusOK, moveResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, moveResponse{State: &state})
}

func (h *Handler) postReset(w http.ResponseWriter, _ *http.Request) {
	state := h.game.Reset()
	writeJSON(w, http.StatusOK, state)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
	}
}
