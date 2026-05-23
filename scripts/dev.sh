#!/usr/bin/env bash
# Run backend + frontend dev server concurrently.
# Backend:  http://localhost:8080  (API)
# Frontend: http://localhost:5173  (open this in the browser)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

cleanup() {
  echo ""
  echo "→ Shutting down..."
  kill "$BACKEND_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

echo "→ Starting backend on :8080..."
cd "$ROOT/backend"
go run . &
BACKEND_PID=$!

echo "→ Starting frontend dev server..."
cd "$ROOT/frontend"
npm install --silent

echo ""
echo "Open http://localhost:5173 in your browser"
echo ""
npm run dev
