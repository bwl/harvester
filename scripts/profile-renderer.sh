#!/usr/bin/env bash
set -euo pipefail
export VR_PROFILE=1
runs=${1:-3}
secs=${2:-5}
for i in $(seq 1 "$runs"); do
  echo "Run $i for ${secs}s..." >&2
  timeout ${secs}s go run ./cmd/game 2>profile_run_${i}.log || true
  echo "Captured: profile_run_${i}.log" >&2
  awk '/\[vr\] render=/{print $3}' profile_run_${i}.log | sed 's/[^0-9\.msµns]//g' > tmp_render_times_${i}.txt || true
  awk '/\[vr\] stringify=/{print $3}' profile_run_${i}.log | sed 's/[^0-9\.msµns]//g' > tmp_stringify_times_${i}.txt || true
  echo "Render stats (approx):" >&2
  cat tmp_render_times_${i}.txt | awk 'END{print NR " frames"}' >&2
  echo >&2
done
