import type { GameState, MoveResponse } from "./types";

const BASE = "/api";

export async function getState(): Promise<GameState> {
  const res = await fetch(`${BASE}/state`);
  if (!res.ok) throw new Error(`GET /api/state failed: ${res.status}`);
  return res.json() as Promise<GameState>;
}

export async function makeMove(x: number, y: number): Promise<MoveResponse> {
  const res = await fetch(`${BASE}/move`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ x, y }),
  });
  if (!res.ok) throw new Error(`POST /api/move failed: ${res.status}`);
  return res.json() as Promise<MoveResponse>;
}

export async function resetGame(): Promise<GameState> {
  const res = await fetch(`${BASE}/reset`, { method: "POST" });
  if (!res.ok) throw new Error(`POST /api/reset failed: ${res.status}`);
  return res.json() as Promise<GameState>;
}
