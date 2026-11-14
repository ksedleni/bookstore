#!/usr/bin/env bash
# Offline local-store demo: put → get → rewrap → get (no CometBFT node).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export PATH="${PATH}:/usr/local/go-new/go/bin"
export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

BIN_DIR="${BIN_DIR:-$ROOT/.bin}"
BIN="$BIN_DIR/bookstored"
HOME_DIR="${HOME_DIR:-/tmp/bookstore-content-demo-$$}"

mkdir -p "$BIN_DIR"
if [[ ! -x "$BIN" ]]; then
  (cd "$ROOT" && go build -o "$BIN" ./cmd/bookstored)
fi

rm -rf "$HOME_DIR"
mkdir -p "$HOME_DIR"

echo "== keys =="
"$BIN" keys add alice --keyring-backend test --home "$HOME_DIR" >/dev/null 2>&1
"$BIN" keys add bob --keyring-backend test --home "$HOME_DIR" >/dev/null 2>&1

BOOK="$HOME_DIR/demo-book.txt"
printf 'Once upon a portable envelope\n' >"$BOOK"
EXPECTED="$(sha256sum "$BOOK" | awk '{print $1}')"

echo "== content put =="
PUT="$("$BIN" content put --from alice --file "$BOOK" --home "$HOME_DIR" --keyring-backend test)"
HASH="$(echo "$PUT" | python3 -c 'import json,sys; print(json.load(sys.stdin)["content_hash"])')"
echo "$PUT"
[[ "$HASH" == "$EXPECTED" ]] || { echo "hash mismatch"; exit 1; }

echo "== content get (alice) =="
"$BIN" content get --from alice --content-hash "$HASH" --out "$HOME_DIR/alice.bin" \
  --home "$HOME_DIR" --keyring-backend test
cmp -s "$BOOK" "$HOME_DIR/alice.bin"

echo "== content rewrap → bob =="
"$BIN" content rewrap --from alice --to bob --content-hash "$HASH" \
  --home "$HOME_DIR" --keyring-backend test

echo "== content get (bob) =="
"$BIN" content get --from bob --content-hash "$HASH" --out "$HOME_DIR/bob.bin" \
  --home "$HOME_DIR" --keyring-backend test
cmp -s "$BOOK" "$HOME_DIR/bob.bin"

"$BIN" content show "$HASH" --home "$HOME_DIR"

echo
echo "CONTENT DEMO PASSED (local store under $HOME_DIR/content)"
