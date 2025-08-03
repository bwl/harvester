#!/usr/bin/env bash
set -euo pipefail

dir=$(cd "$(dirname "$0")/.." && pwd)
cd "$dir"

json='{
  "seed": 1,
  "width": 120,
  "height": 60,
  "dt": 1,
  "steps": [
    {"key":"right","ticks":10},
    {"key":">","ticks":1},
    {"key":"up","ticks":3},
    {"key":"left","ticks":2}
  ]
}'

go run ./cmd/sim <<<"$json"
