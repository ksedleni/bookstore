#!/usr/bin/env bash
# Initialize a single-validator BookStore localnet (Cosmos SDK + CometBFT).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="${BIN_DIR:-$ROOT/.bin}/bookstored"
HOME_DIR="${HOME_DIR:-$HOME/.bookstore}"
CHAIN_ID="${CHAIN_ID:-bookstore-1}"
DENOM="${DENOM:-stake}"
MONIKER="${MONIKER:-localvalidator}"
KEY="${KEY:-alice}"

if [[ ! -x "$BIN" ]]; then
  echo "building bookstored..."
  mkdir -p "$(dirname "$BIN")"
  (cd "$ROOT" && go build -o "$BIN" ./cmd/bookstored)
fi

rm -rf "$HOME_DIR"
"$BIN" init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR" >/dev/null
"$BIN" keys add "$KEY" --keyring-backend test --home "$HOME_DIR" >/dev/null 2>&1
"$BIN" genesis add-genesis-account "$KEY" "1000000000000${DENOM}" --keyring-backend test --home "$HOME_DIR"
"$BIN" genesis gentx "$KEY" "100000000000${DENOM}" --chain-id "$CHAIN_ID" --keyring-backend test --home "$HOME_DIR" >/dev/null
"$BIN" genesis collect-gentxs --home "$HOME_DIR" >/dev/null
"$BIN" genesis validate --home "$HOME_DIR" >/dev/null

# Enable LCD for the optional reference storefront (and local tooling).
APP_TOML="$HOME_DIR/config/app.toml"
if [[ -f "$APP_TOML" ]]; then
  python3 - "$APP_TOML" <<'PY'
import pathlib, re, sys
path = pathlib.Path(sys.argv[1])
text = path.read_text()
m = re.search(r"(?ms)^\[api\](.*?)(?=^\[|\Z)", text)
if m:
    section = m.group(0)
    section = re.sub(r"(?m)^enable\s*=\s*false\b", "enable = true", section, count=1)
    section = re.sub(
        r"(?m)^enabled-unsafe-cors\s*=\s*false\b",
        "enabled-unsafe-cors = true",
        section,
        count=1,
    )
    # Bind IPv4 explicitly so Vite/browser defaults to 127.0.0.1 work.
    section = re.sub(
        r'(?m)^address\s*=\s*".*"',
        'address = "tcp://127.0.0.1:1317"',
        section,
        count=1,
    )
    text = text[: m.start()] + section + text[m.end() :]
    path.write_text(text)
PY
fi

echo "Localnet initialized:"
echo "  home:     $HOME_DIR"
echo "  chain-id: $CHAIN_ID"
echo "  key:      $KEY (test keyring)"
echo "  denom:    $DENOM"
echo "  LCD API:  http://127.0.0.1:1317 (enabled for reference storefront)"
echo
echo "Start with:"
echo "  $BIN start --home $HOME_DIR --minimum-gas-prices=0${DENOM}"
echo "  # or: make node-start"
echo "Storefront:"
echo "  make storefront"
