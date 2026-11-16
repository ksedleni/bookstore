# BookStore

**Fee-free digital-book ownership protocol** — not a fee-taking store.

BookStore is a **principle and open protocol**: an original file becomes a unique, portable ownership record that anyone can mint, transfer, or sell peer-to-peer — like a paperback. Anyone may build a bookstore on top; **none is required**.

| Rule | Value |
|------|-------|
| Protocol / marketplace fee | **0%** |
| Creator royalty on paid sales | **0.1%** (fixed, not seller-configurable) |
| Ledger | Project-owned **Cosmos SDK** sovereign appchain + **CometBFT** |
| Incumbent chains (Ethereum, Solana, Base, …) | **Not required** for mint / list / buy / resell |

## Content uniqueness

Uniqueness is the **SHA-256 digest of the exact file bytes**. The chain stores that digest (and metadata hash) only — not the file contents. Re-encoding the same book produces a different hash and is a different asset. Duplicate remints of the same digest are rejected.

```bash
sha256sum mybook.epub   # → content-hash for mint-original
# Or encrypt + store first (prints the same digest):
.bin/bookstored content put --from alice --file mybook.epub --home ~/.bookstore --keyring-backend test
```

Original ids are **1-based** (the first mint is `id=1`). See [docs/content-delivery.md](docs/content-delivery.md) for owner-gated put/get/rewrap.

## Protocol vs storefronts

- **Protocol (this repo):** ownership, content-hash uniqueness, P2P transfer, P2P paid sale, optional on-ledger listings, open events.
- **Storefronts (optional):** any third-party UI. Filtering policies belong to UIs, not the ledger.
- **Bookstore-less trade:** `bookstored tx bookstore p2p-sale` completes a sale with only keys + a ledger client.

## Requirements

- Go **1.22+** (local builds may use a newer Go toolchain)
- No Docker required for the single-node localnet

## Quick start — local Cosmos node

```bash
git clone https://github.com/ksedleni/bookstore.git
cd bookstore

make build
make smoke-node          # init + start until a block commits, then exit
make e2e                 # live mint → fund buyer → P2P sale assertions

# Persistent localnet:
make node-init           # writes ~/.bookstore (destructive)
make node-start          # bookstored start (CometBFT in-process)
```

RPC: `http://127.0.0.1:26657`. Address prefix: `store`. Bond denom: `stake`. Chain id: `bookstore-1`.

## Walkthrough — mint → fund buyer → P2P sale

Assumes `make node-init` already created `alice` (validator + funded) under `~/.bookstore` with keyring `test`, and `make node-start` is running in another terminal.

```bash
BIN=.bin/bookstored
HOME=~/.bookstore
CHAIN=bookstore-1
COMMON=(--home "$HOME" --keyring-backend test --chain-id "$CHAIN" --yes)
TX=("${COMMON[@]}" --gas auto --gas-adjustment 1.5 --fees 200stake)

# Optional second buyer key + fund from alice
$BIN keys add bob --keyring-backend test --home "$HOME"
BOB=$($BIN keys show bob -a --keyring-backend test --home "$HOME")
ALICE=$($BIN keys show alice -a --keyring-backend test --home "$HOME")
$BIN tx bank send alice "$BOB" 5000000stake "${TX[@]}"
# Wait until the send is indexed before continuing:
#   $BIN q tx <txhash> --home "$HOME"

# Encrypt file for alice, store locally; digest only goes on-chain
printf 'hello paperback\n' > /tmp/book.txt
PUT=$($BIN content put --from alice --file /tmp/book.txt --home "$HOME" --keyring-backend test)
HASH=$(echo "$PUT" | python3 -c 'import json,sys; print(json.load(sys.stdin)["content_hash"])')
META=$(printf '{"title":"Hello"}' | sha256sum | awk '{print $1}')

# Mint (first original is id 1)
$BIN tx bookstore mint-original "$HASH" "$META" "" --from alice "${TX[@]}"
# Wait for q tx <hash> before querying owner
$BIN q bookstore show-original 1 --home "$HOME"
# → owner = alice, id = "1"

# Owner fetches bytes (SHA-256 verified on decrypt)
$BIN content get --from alice --content-hash "$HASH" --out /tmp/book-out.txt \
  --home "$HOME" --keyring-backend test

# P2P sale: signer is the buyer; no storefront UI
$BIN tx bookstore p2p-sale 1 "$ALICE" 1000000 stake --from bob "${TX[@]}"
# Wait for q tx <hash>, then:
$BIN q bookstore show-original 1 --home "$HOME"
# → owner = bob
# Settlement: 0.1% royalty to creator, 0% protocol fee

# Bookstore-less content handoff: seller rewraps CEK for buyer
$BIN content rewrap --from alice --to bob --content-hash "$HASH" \
  --home "$HOME" --keyring-backend test
$BIN content get --from bob --original-id 1 --out /tmp/book-bob.txt \
  --home "$HOME" --keyring-backend test --node tcp://127.0.0.1:26657
```

Always wait for `bookstored q tx <txhash>` to succeed before relying on `show-original` / balances.

## On-chain messages

**`x/bookstore`**

- `mint-original` — permissionless mint from content hash (duplicate digest rejected)
- `transfer-original` — unpaid P2P transfer (no royalty)
- `p2p-sale` — paid P2P sale: buyer pays, **0.1%** royalty, **0%** protocol fee

**`x/marketplace`** (optional)

- `create-listing` / `cancel-listing` / `buy-listing` — same royalty rules
- `set-paused` — **narrow** marketplace pause only; does **not** block mint or P2P

## Offline protocol CLI (no node)

```bash
make demo
go run ./client/cmd/bookstore --help
```

## Optional reference storefront (publisher-friendly)

Thin Vite + React UI under [`web/`](web/): publish editions, list, and buy on localnet via a small local helper. **CLI / P2P remains primary** — this is *a* storefront, not the store (0% protocol fee, 0.1% royalty).

```bash
make node-init
make node-start
```

```bash
make publisher
```

```bash
make storefront
```

Open http://127.0.0.1:5173 — set editions to `2`, mint & list, then buy from the shop. See [web/README.md](web/README.md).

## Repository layout

```
app/                      # Cosmos SDK app wiring (depinject)
cmd/bookstored/           # Node + CLI binary (init/start/tx/query/content)
content/                  # Portable AES-GCM envelopes + local store
web/                      # Optional storefront + local publisher helper
x/bookstore/              # Ownership + P2P sale module
x/marketplace/            # Optional listings + pause
proto/                    # Protobuf definitions
protocol/                 # Portable royalty + local-ledger helpers
client/cmd/bookstore/     # Offline reference CLI
scripts/localnet-init.sh  # Single-validator genesis
scripts/smoke-node.sh     # Start-until-block smoke check
scripts/e2e-localnet.sh   # Live mint → content → P2P sale → rewrap e2e
docs/invariants.md        # Locked product invariants
docs/content-delivery.md  # Envelope format + put/get/rewrap
```

## What’s still stubbed / next

- IPFS/Arweave pinning + HTTP content gateway
- Automatic CEK rewrap inside `MsgP2PSale`
- Full indexer / public multi-node testnet docs
- Keplr (or similar) signing in the reference UI (helper covers localnet today)
- Bridges — explicitly **out of v1**

## License

Apache-2.0 — see [LICENSE](LICENSE).
