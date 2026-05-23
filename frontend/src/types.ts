export type Stone = 0 | 1 | 2; // 0=empty 1=black 2=white
export type Player = "black" | "white";
export type PlaceMode = "black" | "white" | "alternate";

export interface GameState {
  board: Stone[][];
  blackCaptures: number;
  whiteCaptures: number;
}

export interface MoveResponse {
  state: GameState | null;
  error?: string;
}
