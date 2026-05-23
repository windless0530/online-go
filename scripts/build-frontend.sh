#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/frontend"

echo "→ Installing frontend dependencies..."
npm install

echo "→ Building frontend..."
npm run build

echo "✓ Frontend built to frontend/dist/"
