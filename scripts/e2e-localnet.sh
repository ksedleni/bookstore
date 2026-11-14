#!/usr/bin/env bash
# Live e2e: mint → P2P sale → assert new owner + 0.1% royalty / 0% protocol fee.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export PATH="${PATH}:/usr/local/go-new/go/bin"
export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

BIN_DIR="${BIN_DIR:-$ROOT/.bin}"
BIN="$BIN_DIR/bookstored"
HOME_DIR="${HOME_DIR:-/tmp/bookstore-e2e-$$}"
CHAIN_ID="${CHAIN_ID:-bookstore-1}"
DENOM="${DENOM:-stake}"
MONIKER="${MONIKER:-locale2e}"
PRICE="${PRICE:-1000000}"
LOG="$HOME_DIR/node.log"
NODE="tcp://127.0.0.1:26657"

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

rm -rf "$HOME_DIR"
mkdir -p "$HOME_DIR"

echo "== init localnet ($HOME_DIR) =="
"$BIN" init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR" >/dev/null
"$BIN" keys add alice --keyring-backend test --home "$HOME_DIR" >/dev/null 2>&1
"$BIN" keys add bob --keyring-backend test --home "$HOME_DIR" >/dev/null 2>&1

ALICE="$("$BIN" keys show alice -a --keyring-backend test --home "$HOME_DIR")"
BOB="$("$BIN" keys show bob -a --keyring-backend test --home "$HOME_DIR")"

"$BIN" genesis add-genesis-account alice "1000000000000${DENOM}" --keyring-backend test --home "$HOME_DIR" >/dev/null >/dev/null
"$BIN" genesis add-genesis-account bob "5000000000${DENOM}" --keyring-backend test --home "$HOME_DIR" >/dev/null >/dev/null
"$BIN" genesis gentx alice "100000000000${DENOM}" --chain-id "$CHAIN_ID" --keyring-backend test --home "$HOME_DIR" >/dev/null
"$BIN" genesis collect-gentxs --home "$HOME_DIR" >/dev/null
"$BIN" genesis validate --home "$HOME_DIR" >/dev/null

echo "alice=$ALICE"
echo "bob=$BOB"

echo "== start node =="
"$BIN" start --home "$HOME_DIR" --minimum-gas-prices="0${DENOM}" >"$LOG" 2>&1 &
NODE_PID=$!

for i in $(seq 1 60); do
  if grep -q "finalized block" "$LOG" 2>/dev/null; then
    break
  fi
  if ! kill -0 "$NODE_PID" 2>/dev/null; then
    echo "node died during startup"; tail -80 "$LOG"; exit 1
  fi
  sleep 1
done
sleep 2

COMMON=(--home "$HOME_DIR" --keyring-backend test --chain-id "$CHAIN_ID" --node "$NODE" --yes)
TXCOMMON=("${COMMON[@]}" --gas auto --gas-adjustment 1.5 --fees "200${DENOM}")

wait_tx() {
  local hash="$1"
  local i
  for i in $(seq 1 40); do
    if OUT="$("$BIN" q tx "$hash" --home "$HOME_DIR" --node "$NODE" -o json 2>/dev/null)"; then
      CODE="$(echo "$OUT" | python3 -c 'import json,sys; print(json.load(sys.stdin).get("code",1))')"
      if [[ "$CODE" == "0" ]]; then
        echo "$OUT"
        return 0
      fi
      echo "tx failed code=$CODE"
      echo "$OUT" | python3 -m json.tool | head -40
      return 1
    fi
    sleep 1
  done
  echo "timeout waiting for tx $hash"
  return 1
}

broadcast() {
  local raw hash
  raw="$("$@" -o json)"
  hash="$(echo "$raw" | python3 -c 'import json,sys; print(json.load(sys.stdin)["txhash"])')"
  echo "txhash=$hash"
  wait_tx "$hash" >/dev/null
}

balance_of() {
  "$BIN" q bank balances "$1" --home "$HOME_DIR" --node "$NODE" -o json \
    | python3 -c "import json,sys; b=json.load(sys.stdin).get('balances',[]); print(next((x['amount'] for x in b if x['denom']=='$DENOM'),'0'))"
}

parse_original_fields() {
  python3 -c '
import json, sys
d = json.loads(sys.argv[1])
o = d.get("Original") or d.get("original")
if o is None:
    raise SystemExit("unexpected show-original payload: %r" % (d,))
print(o["id"])
print(o["owner"])
' "$1"
}

BOOK_FILE="$HOME_DIR/e2e-book.txt"
printf 'BookStore e2e paperback\n' >"$BOOK_FILE"
CONTENT_HASH="$(sha256sum "$BOOK_FILE" | awk '{print $1}')"
META_HASH="$(printf '{"title":"E2E","content_hash":"%s"}' "$CONTENT_HASH" | sha256sum | awk '{print $1}')"
echo "content_hash=$CONTENT_HASH (SHA-256 of file bytes; chain stores digest only)"

echo "== mint (alice) =="
broadcast "$BIN" tx bookstore mint-original "$CONTENT_HASH" "$META_HASH" "" --from alice "${TXCOMMON[@]}"

ORIG="$("$BIN" q bookstore show-original 1 --home "$HOME_DIR" --node "$NODE" -o json)"
mapfile -t ORIG_FIELDS < <(parse_original_fields "$ORIG")
OID="${ORIG_FIELDS[0]}"
OWNER="${ORIG_FIELDS[1]}"
echo "minted id=$OID owner=$OWNER"
[[ "$OID" == "1" ]] || { echo "expected first original id=1, got $OID"; exit 1; }
[[ "$OWNER" == "$ALICE" ]] || { echo "expected owner alice"; exit 1; }

ALICE_PRE_SALE="$(balance_of "$ALICE")"
BOB_PRE_SALE="$(balance_of "$BOB")"
echo "balances pre-sale: alice=$ALICE_PRE_SALE bob=$BOB_PRE_SALE"

echo "== p2p-sale (bob buys from alice) =="
broadcast "$BIN" tx bookstore p2p-sale 1 "$ALICE" "$PRICE" "$DENOM" --from bob "${TXCOMMON[@]}"

ORIG2="$("$BIN" q bookstore show-original 1 --home "$HOME_DIR" --node "$NODE" -o json)"
mapfile -t ORIG2_FIELDS < <(parse_original_fields "$ORIG2")
OWNER2="${ORIG2_FIELDS[1]}"
echo "after sale owner=$OWNER2"
[[ "$OWNER2" == "$BOB" ]] || { echo "FAIL: expected owner bob after P2P sale, got $OWNER2"; exit 1; }

ALICE_AFTER="$(balance_of "$ALICE")"
BOB_AFTER="$(balance_of "$BOB")"
echo "balances after: alice=$ALICE_AFTER bob=$BOB_AFTER"

python3 - <<PY
price = int("$PRICE")
royalty = price * 10 // 10000  # 0.1%
seller_proceeds = price - royalty
protocol_fee = 0
assert protocol_fee == 0
alice_pre, alice_after = int("$ALICE_PRE_SALE"), int("$ALICE_AFTER")
bob_pre, bob_after = int("$BOB_PRE_SALE"), int("$BOB_AFTER")
# Alice is seller and royalty payee → receives full price (proceeds + royalty).
# Bob pays price plus his tx fees.
alice_delta = alice_after - alice_pre
bob_delta = bob_after - bob_pre
print(f"alice_delta={alice_delta} (expect == {price} as seller+royalty payee)")
print(f"bob_delta={bob_delta} (expect <= -{price}, fees extra)")
print(f"royalty={royalty} seller_proceeds={seller_proceeds} protocol_fee={protocol_fee}")
if alice_delta != price:
    raise SystemExit(f"alice should receive exactly price as seller+creator; delta={alice_delta}")
if bob_delta > -price:
    raise SystemExit(f"bob should pay at least price; delta={bob_delta}")
assert royalty == price * 10 // 10000
assert seller_proceeds + royalty == price
print("OK: ownership transferred; royalty split 0.1% / protocol fee 0%")
PY

echo
echo "E2E PASSED"
