# BookStore

**Fee-free digital-book ownership protocol** — not a fee-taking store.

BookStore is a **principle and open protocol**: an original file becomes a unique, portable ownership record that anyone can mint, transfer, or sell peer-to-peer — like a paperback. Anyone may build a bookstore on top; **none is required**.

| Rule | Value |
|------|-------|
| Protocol / marketplace fee | **0%** |
| Creator royalty on paid sales | **0.1%** (fixed, not seller-configurable) |
| Ledger | Project-owned **Cosmos SDK** sovereign appchain + **CometBFT** |
| Incumbent chains (Ethereum, Solana, Base, …) | **Not required** for mint / list / buy / resell |

## Protocol vs storefronts

- **Protocol (this repo):** ownership, content-hash uniqueness, P2P transfer, P2P paid sale, optional on-ledger listings, open events.
- **Storefronts (optional):** any third-party UI. Filtering policies belong to UIs, not the ledger.
- **Bookstore-less trade:** `MsgP2PSale` / CLI demo complete a sale with only keys + a ledger client.

## Requirements

- Go **1.22+** (CI/local builds use Go 1.26 toolchain when needed)
- No Docker required for the single-node localnet

## Quick start — local Cosmos node

```bash
git clone https://github.com/ksedleni/bookstore.git
cd bookstore

make build
make smoke-node          # init + start until a block commits, then exit

# Or run a persistent localnet:
make node-init           # writes ~/.bookstore (destructive)
make node-start          # bookstored start (CometBFT in-process)
```

Manual equivalent:

```bash
go build -o .bin/bookstored ./cmd/bookstored
./scripts/localnet-init.sh
.bin/bookstored start --home ~/.bookstore --minimum-gas-prices=0stake
```

RPC defaults to `http://127.0.0.1:26657`. Address prefix: `store`. Bond denom: `stake`.

### On-chain messages (wired)

**`x/bookstore`**

- `MintOriginal` — permissionless mint from content hash (duplicate hash rejected)
- `TransferOriginal` — unpaid P2P transfer (no royalty)
- `P2PSale` — paid P2P sale: buyer pays, **0.1%** royalty to royalty payee, **0%** protocol fee, ownership moves

**`x/marketplace`** (optional convenience)

- `CreateListing` / `CancelListing` / `BuyListing` — same royalty rules
- `SetPaused` — **narrow** marketplace pause only (authority); does **not** block mint or P2P

## Offline protocol CLI (no node)

Still available for royalty/ledger demos without CometBFT:

```bash
make demo
# or
go run ./client/cmd/bookstore --help
```

## Repository layout

```
app/                      # Cosmos SDK app wiring (depinject)
cmd/bookstored/           # Node + CLI binary (init/start/tx/query)
x/bookstore/              # Ownership + P2P sale module
x/marketplace/            # Optional listings + pause
proto/                    # Protobuf definitions
protocol/                 # Portable royalty + local-ledger helpers
client/cmd/bookstore/     # Offline reference CLI (bookstore-less demo)
scripts/localnet-init.sh  # Single-validator genesis
scripts/smoke-node.sh     # Start-until-block smoke check
docs/invariants.md        # Locked product invariants
```

## What’s still stubbed / next

- Full indexer / public multi-node testnet docs
- End-to-end tx integration tests against a live localnet (mint → P2P sale via `bookstored tx`)
- Reference web storefront (optional by design)
- Bridges — explicitly **out of v1**

## License

Apache-2.0 — see [LICENSE](LICENSE).
