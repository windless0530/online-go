import { useEffect, useState, useCallback } from "react";
import { Board } from "./Board";
import { getState, makeMove, resetGame } from "./api";
import type { GameState } from "./types";
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
  const [moving, setMoving] = useState(false);

  useEffect(() => {
    getState()
      .then(setState)
      .catch(() => setError("无法连接到服务器，请确认后端已启动"))
      .finally(() => setLoading(false));
  }, []);

  const handleMove = useCallback(
    async (x: number, y: number) => {
      if (moving) return;
      setMoving(true);
      setError(null);
      try {
        const res = await makeMove(x, y);
        if (res.error) {
          setError(res.error);
        } else if (res.state) {
          setState(res.state);
        }
      } catch {
        setError("落子失败，请重试");
      } finally {
        setMoving(false);
      }
    },
    [moving]
  );

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
          {state.currentPlayer === "black" ? "黑方" : "白方"}
          {moving ? " ..." : " 落子"}
        </span>
        <span className="capture-info">黑提 {state.blackCaptures} 子</span>
        <span className="capture-info">白提 {state.whiteCaptures} 子</span>
      </div>

      {error && <p className="error-banner">{error}</p>}

      <Board
        board={state.board}
        currentPlayer={state.currentPlayer}
        onMove={handleMove}
        disabled={loading || moving}
      />

      <button className="reset-btn" onClick={handleReset} disabled={loading || moving}>
        重新开始
      </button>
    </div>
  );
}
