#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/backend"

echo "→ Running Go tests..."
go test ./...

echo "→ Building backend..."
mkdir -p "$ROOT/bin"
go build -o "$ROOT/bin/server" .

echo "✓ Backend built to bin/server"
