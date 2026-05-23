import { useEffect, useState, useCallback, useRef } from "react";
import { Board } from "./Board";
import { getState, makeMove, resetGame } from "./api";
import type { GameState, Stone } from "./types";
import "./App.css";

const BOARD_SIZE = 19;

function emptyState(): GameState {
  return {
    board: Array.from({ length: BOARD_SIZE }, () => Array<0>(BOARD_SIZE).fill(0)),
    currentPlayer: "black",
    blackCaptures: 0,
    whiteCaptures: 0,
  };
}

export function App() {
  const [state, setState] = useState<GameState>(emptyState);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  // Use a ref instead of state so pending-move tracking never triggers a re-render.
  const pendingRef = useRef(false);
  // Keep a ref to the latest state so the stable handleMove callback can read it.
  const stateRef = useRef(state);
  stateRef.current = state;

  useEffect(() => {
    getState()
      .then(setState)
      .catch(() => setError("无法连接到服务器，请确认后端已启动"))
      .finally(() => setLoading(false));
  }, []);

  // Empty dep array → stable reference; reads live state via stateRef.
  const handleMove = useCallback(async (x: number, y: number) => {
    if (pendingRef.current) return;
    pendingRef.current = true;
    setError(null);

    const snapshot = stateRef.current;
    const stone: Stone = snapshot.currentPlayer === "black" ? 1 : 2;
    const nextPlayer = snapshot.currentPlayer === "black" ? "white" : "black";

    // Optimistic update: stone appears immediately without waiting for the API.
    setState(prev => ({
      ...prev,
      board: prev.board.map((col, ci) =>
        col.map((s, ri) => (ci === x && ri === y ? stone : s))
      ),
      currentPlayer: nextPlayer,
    }));

    try {
      const res = await makeMove(x, y);
      if (res.error) {
        setState(snapshot); // revert to pre-move state
        setError(res.error);
      } else if (res.state) {
        setState(res.state); // apply authoritative state (handles captures)
      }
    } catch {
      setState(snapshot);
      setError("落子失败，请重试");
    } finally {
      pendingRef.current = false;
    }
  }, []); // stable — no deps needed because we read state via ref

  const handleReset = useCallback(async () => {
    setError(null);
    setLoading(true);
    try {
      const s = await resetGame();
      setState(s);
    } catch {
      setError("重置失败");
    } finally {
      setLoading(false);
    }
  }, []);

  return (
    <div className="app">
      <h1 className="app-title">围棋</h1>

      <div className="status-bar">
        <span className={`player-chip player-chip--${state.currentPlayer}`}>
          {state.currentPlayer === "black" ? "黑方" : "白方"}落子
        </span>
        <span className="capture-info">黑提 {state.blackCaptures} 子</span>
        <span className="capture-info">白提 {state.whiteCaptures} 子</span>
      </div>

      {error && <p className="error-banner">{error}</p>}

      <Board
        board={state.board}
        currentPlayer={state.currentPlayer}
        onMove={handleMove}
        disabled={loading}
      />

      <button className="reset-btn" onClick={handleReset} disabled={loading}>
        重新开始
      </button>
    </div>
  );
}
