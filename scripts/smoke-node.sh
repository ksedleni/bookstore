#!/usr/bin/env bash
# Smoke-test: init a throwaway localnet, start bookstored, wait for a committed block.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export HOME_DIR="${HOME_DIR:-/tmp/bookstore-smoke-$$}"
export BIN_DIR="${BIN_DIR:-$ROOT/.bin}"
export CHAIN_ID="${CHAIN_ID:-bookstore-1}"
export DENOM="${DENOM:-stake}"
BIN="$BIN_DIR/bookstored"
LOG="$HOME_DIR/smoke.log"

cleanup() {
  if [[ -n "${NODE_PID:-}" ]] && kill -0 "$NODE_PID" 2>/dev/null; then
    kill "$NODE_PID" 2>/dev/null || true
    wait "$NODE_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

mkdir -p "$BIN_DIR"
if [[ ! -x "$BIN" ]]; then
  (cd "$ROOT" && go build -o "$BIN" ./cmd/bookstored)
fi

"$ROOT/scripts/localnet-init.sh"

mkdir -p "$HOME_DIR"
"$BIN" start --home "$HOME_DIR" --minimum-gas-prices="0${DENOM}" >"$LOG" 2>&1 &
NODE_PID=$!

echo "waiting for committed block (pid=$NODE_PID)..."
for i in $(seq 1 60); do
  if grep -q "finalized block" "$LOG" 2>/dev/null; then
    HEIGHT=$(grep -o 'height=[0-9]*' "$LOG" | tail -1 | cut -d= -f2 || true)
    echo "OK: node produced blocks (last seen height=${HEIGHT:-?})"
    echo "genesis modules include bookstore + marketplace:"
    python3 - <<PY
import json
g=json.load(open("$HOME_DIR/config/genesis.json"))
assert "bookstore" in g["app_state"], "bookstore missing from genesis"
assert "marketplace" in g["app_state"], "marketplace missing from genesis"
print("  bookstore:", list(g["app_state"]["bookstore"].keys()))
print("  marketplace:", list(g["app_state"]["marketplace"].keys()))
PY
    exit 0
  fi
  if ! kill -0 "$NODE_PID" 2>/dev/null; then
    echo "node exited early; log:"
    tail -50 "$LOG" || true
    exit 1
  fi
  sleep 1
done

echo "timeout waiting for block; log:"
tail -80 "$LOG" || true
exit 1
