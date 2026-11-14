# BookStore

**Fee-free digital-book ownership protocol** — not a fee-taking store.

BookStore is a **principle and open protocol**: an original file becomes a unique, portable ownership record that anyone can mint, transfer, or sell peer-to-peer — like a paperback. Anyone may build a bookstore on top; **none is required**.

| Rule | Value |
|------|-------|
| Protocol / marketplace fee | **0%** |
| Creator royalty on paid sales | **0.1%** (fixed, not seller-configurable) |
| Ledger | Project-owned **Cosmos SDK** sovereign appchain |
| Incumbent chains (Ethereum, Solana, Base, …) | **Not required** for mint / list / buy / resell |

## Protocol vs storefronts

- **Protocol (this repo):** ownership, content-hash uniqueness, P2P transfer, P2P paid sale, optional on-ledger listings, open events.
- **Storefronts (optional):** any third-party UI or the future reference web app. Filtering policies belong to UIs, not the ledger.
- **Proof of bookstore-less trade:** the CLI below completes mint → P2P sale with only keys + a ledger client.

## Repository layout

```
protocol/                 # Portable domain rules (royalty, ledger semantics, events)
chain/
  app/                    # App wiring (MVP stub → full Cosmos SDK app)
  cmd/bookstored/         # Sovereign node binary entrypoint (scaffold)
  x/bookstore/            # Mint, ownership, P2P transfer/sale (0.1% royalty)
  x/marketplace/          # Optional listings + narrow pause (does not block P2P)
client/cmd/bookstore/     # Reference CLI — mint/list/buy without a storefront
protocol/schemas/         # Event / object schemas for open indexers
docs/                     # ADRs and operator notes
```

## Quick start (local demo ledger)

Requires Go 1.22+.

```bash
git clone https://github.com/ksedleni/bookstore.git
cd bookstore
make test
make build
make demo
```

Or step by step:

```bash
export BOOKSTORE_HOME=/tmp/bookstore-demo
rm -rf "$BOOKSTORE_HOME"

go run ./client/cmd/bookstore faucet --to alice --amount 0
go run ./client/cmd/bookstore faucet --to bob --amount 5000000
printf 'hello paperback\n' > /tmp/book.txt
go run ./client/cmd/bookstore mint --from alice --file /tmp/book.txt --title "Hello" --author "A"
go run ./client/cmd/bookstore sell --buyer bob --seller alice --id book-1 --price 1000000
go run ./client/cmd/bookstore show --id book-1
```

Settlement on that sale: seller receives 99.9%; creator royalty 0.1%; **protocol fee 0**.

Marketplace pause (listings only):

```bash
go run ./client/cmd/bookstore pause-market --paused true
# list/buy blocked; mint + P2P sell/transfer still work
```

Node stub:

```bash
go run ./chain/cmd/bookstored version
```

## What is intentionally scaffold-only

This first commit ships a **buildable MVP scaffold**: protocol invariants, module layout, and a reference CLI over a local ledger. Full Cosmos SDK genesis / CometBFT `start`, public RPC, and light-client sync are the next spike (see architecture plan). Domain rules already match the locked product model so the chain wiring cannot silently introduce a protocol fee or incumbent-chain dependency.

## License

Apache-2.0 — see [LICENSE](LICENSE).
