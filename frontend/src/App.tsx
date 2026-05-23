import { useEffect, useState, useCallback, useRef } from "react";
import { Board } from "./Board";
import { getState, makeMove, resetGame } from "./api";
import type { GameState, Stone, Player, PlaceMode } from "./types";
import "./App.css";

const BOARD_SIZE = 19;

function emptyState(): GameState {
  return {
    board: Array.from({ length: BOARD_SIZE }, () => Array<0>(BOARD_SIZE).fill(0)),
    blackCaptures: 0,
    whiteCaptures: 0,
  };
}

const MODE_LABELS: Record<PlaceMode, string> = {
  black: "摆黑子",
  white: "摆白子",
  alternate: "交替落子",
};

export function App() {
  const [state, setState] = useState<GameState>(emptyState);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [placeMode, setPlaceMode] = useState<PlaceMode>("alternate");
  const [currentTurn, setCurrentTurn] = useState<Player>("black");

  const pendingRef = useRef(false);
  const stateRef = useRef(state);
  const turnRef = useRef(currentTurn);
  const modeRef = useRef(placeMode);
  stateRef.current = state;
  turnRef.current = currentTurn;
  modeRef.current = placeMode;

  useEffect(() => {
    getState()
      .then(setState)
      .catch(() => setError("无法连接到服务器，请确认后端已启动"))
      .finally(() => setLoading(false));
  }, []);

  const activeColor = (mode: PlaceMode, turn: Player): Player =>
    mode === "alternate" ? turn : mode;

  const handleMove = useCallback(async (x: number, y: number) => {
    if (pendingRef.current) return;
    pendingRef.current = true;
    setError(null);

    const snapshot = stateRef.current;
    const mode = modeRef.current;
    const turn = turnRef.current;
    const color = activeColor(mode, turn);
    const stone: Stone = color === "black" ? 1 : 2;

    // Optimistic update
    setState(prev => ({
      ...prev,
      board: prev.board.map((col, ci) =>
        col.map((s, ri) => (ci === x && ri === y ? stone : s))
      ),
    }));

    // Flip turn for alternate mode immediately so hover preview updates
    if (mode === "alternate") {
      setCurrentTurn(t => (t === "black" ? "white" : "black"));
    }

    try {
      const res = await makeMove(x, y, color);
      if (res.error) {
        setState(snapshot);
        if (mode === "alternate") setCurrentTurn(color); // revert turn flip
        setError(res.error);
      } else if (res.state) {
        setState(res.state);
      }
    } catch {
      setState(snapshot);
      if (mode === "alternate") setCurrentTurn(color);
      setError("落子失败，请重试");
    } finally {
      pendingRef.current = false;
    }
  }, []);

  const handleReset = useCallback(async () => {
    setError(null);
    setLoading(true);
    setCurrentTurn("black");
    try {
      const s = await resetGame();
      setState(s);
    } catch {
      setError("重置失败");
    } finally {
      setLoading(false);
    }
  }, []);

  const hoverColor = activeColor(placeMode, currentTurn);

  return (
    <div className="app">
      <h1 className="app-title">围棋</h1>

      <div className="capture-bar">
        <span className="capture-info">黑提 {state.blackCaptures} 子</span>
        <span className="capture-info">白提 {state.whiteCaptures} 子</span>
      </div>

      {error && <p className="error-banner">{error}</p>}

      <Board
        board={state.board}
        currentPlayer={hoverColor}
        onMove={handleMove}
        disabled={loading}
      />

      <div className="bottom-bar">
        <div className="mode-selector" role="group" aria-label="落子模式">
          {(["black", "white", "alternate"] as PlaceMode[]).map(mode => (
            <label
              key={mode}
              className={`mode-option mode-option--${mode}${placeMode === mode ? " mode-option--active" : ""}`}
            >
              <input
                type="radio"
                name="placeMode"
                value={mode}
                checked={placeMode === mode}
                onChange={() => setPlaceMode(mode)}
              />
              {MODE_LABELS[mode]}
            </label>
          ))}
        </div>

        <button type="button" className="reset-btn" onClick={handleReset} disabled={loading}>
          重新开始
        </button>
      </div>
    </div>
  );
}
