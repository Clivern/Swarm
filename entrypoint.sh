#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${OPENROUTER_API_KEY:-}" ]]; then
  echo "OPENROUTER_API_KEY is required" >&2
  exit 1
fi
if [[ -z "${PROMPT:-}" ]]; then
  echo "PROMPT is required" >&2
  exit 1
fi
if [[ -z "${PI_MODEL:-}" ]]; then
  echo "PI_MODEL is required" >&2
  exit 1
fi
if [[ ! -d /repo/.git ]]; then
  echo "/repo must be a mount of a git repository" >&2
  exit 1
fi

cd /repo
BASE="$(git rev-parse HEAD)"

echo "Running Pi (model=${PI_MODEL}) in /repo ..." >&2
pi -p --model "$PI_MODEL" --no-session "$PROMPT" | tee /tmp/pi-summary.txt

DIFF="$(git diff "$BASE")"
if [[ -d /out ]]; then
  printf '%s' "$DIFF" > /out/patch.diff
  cp /tmp/pi-summary.txt /out/summary.txt
  echo "Wrote /out/patch.diff and /out/summary.txt" >&2
else
  echo "--- patch.diff ---" >&2
  printf '%s' "$DIFF"
fi

if [[ -z "$DIFF" ]]; then
  echo "warning: empty diff (no tracked file changes?)" >&2
fi
