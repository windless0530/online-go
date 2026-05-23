import { memo, useState, useCallback } from "react";
import type { Stone, Player } from "./types";
import "./Board.css";

const LINES = 19;
const CELL = 51;
const PAD = 45;
const STONE_R = 22;
const BOARD_PX = PAD * 2 + CELL * (LINES - 1);

const STAR_POINTS: [number, number][] = [
  [3, 3], [3, 9], [3, 15],
  [9, 3], [9, 9], [9, 15],
  [15, 3], [15, 9], [15, 15],
];

function gx(i: number): number {
  return PAD + i * CELL;
}

interface Props {
  board: Stone[][];
  currentPlayer: Player;
  onMove: (x: number, y: number) => void;
  disabled: boolean;
}

export const Board = memo(function Board({ board, currentPlayer, onMove, disabled }: Props) {
  const [hover, setHover] = useState<[number, number] | null>(null);

  const toGrid = useCallback(
    (clientX: number, clientY: number, rect: DOMRect): [number, number] | null => {
      const scale = BOARD_PX / rect.width;
      const svgX = (clientX - rect.left) * scale;
      const svgY = (clientY - rect.top) * scale;
      const col = Math.round((svgX - PAD) / CELL);
      const row = Math.round((svgY - PAD) / CELL);
      if (col < 0 || col >= LINES || row < 0 || row >= LINES) return null;
      if (Math.abs(svgX - gx(col)) > CELL * 0.46 || Math.abs(svgY - gx(row)) > CELL * 0.46)
        return null;
      return [col, row];
    },
    []
  );

  const handleMouseMove = useCallback(
    (e: React.MouseEvent<SVGSVGElement>) => {
      if (disabled) {
        if (hover !== null) setHover(null);
        return;
      }
      setHover(toGrid(e.clientX, e.clientY, e.currentTarget.getBoundingClientRect()));
    },
    [disabled, hover, toGrid]
  );

  const handleMouseLeave = useCallback(() => setHover(null), []);

  const handleClick = useCallback(
    (e: React.MouseEvent<SVGSVGElement>) => {
      if (disabled) return;
      const pos = toGrid(e.clientX, e.clientY, e.currentTarget.getBoundingClientRect());
      if (pos && board[pos[0]][pos[1]] === 0) {
        onMove(pos[0], pos[1]);
        setHover(null);
      }
    },
    [disabled, toGrid, board, onMove]
  );

  const hoverEmpty = hover !== null && board[hover[0]][hover[1]] === 0;

  return (
    <svg
      className="board"
      viewBox={`0 0 ${BOARD_PX} ${BOARD_PX}`}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      onClick={handleClick}
    >
      <defs>
        <radialGradient id="g-black" cx="38%" cy="32%" r="60%">
          <stop offset="0%" stopColor="#6a6a6a" />
          <stop offset="100%" stopColor="#0a0a0a" />
        </radialGradient>
        <radialGradient id="g-white" cx="40%" cy="32%" r="60%">
          <stop offset="0%" stopColor="#ffffff" />
          <stop offset="100%" stopColor="#c0c0c0" />
        </radialGradient>
      </defs>

      <rect x={0} y={0} width={BOARD_PX} height={BOARD_PX} className="board-bg" rx={4} />

      {Array.from({ length: LINES }, (_, i) => (
        <g key={i}>
          <line x1={gx(i)} y1={gx(0)} x2={gx(i)} y2={gx(LINES - 1)} className="grid-line" />
          <line x1={gx(0)} y1={gx(i)} x2={gx(LINES - 1)} y2={gx(i)} className="grid-line" />
        </g>
      ))}

      {STAR_POINTS.map(([col, row]) => (
        <circle key={`s${col}-${row}`} cx={gx(col)} cy={gx(row)} r={5} className="star-point" />
      ))}

      {board.map((col, x) =>
        col.map((stone, y) => {
          if (stone === 0) return null;
          return (
            <circle
              key={`${x}-${y}`}
              cx={gx(x)}
              cy={gx(y)}
              r={STONE_R}
              fill={stone === 1 ? "url(#g-black)" : "url(#g-white)"}
              stroke={stone === 2 ? "#9a9a9a" : "none"}
              strokeWidth={stone === 2 ? 0.8 : 0}
              className="stone"
            />
          );
        })
      )}

      {hover && hoverEmpty && !disabled && (
        <circle
          cx={gx(hover[0])}
          cy={gx(hover[1])}
          r={STONE_R}
          className={`hover-stone hover-stone--${currentPlayer}`}
          style={{ pointerEvents: "none" }}
        />
      )}
    </svg>
  );
});
