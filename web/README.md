# BookStore reference storefront

Optional, thin, forkable web UI over the BookStore ledger. **The protocol remains usable without it** — mint, P2P sale, and content put/get stay on the `bookstored` CLI.

| Principle | Value |
|-----------|-------|
| Protocol fee | **0%** |
| Creator royalty | **0.1%** on paid sales |
| Custody | **None** — this UI does not hold keys |

## Run

Prerequisites: Node 20+, a local `bookstored` node with the API enabled (LCD default `http://127.0.0.1:1317`).

```bash
# terminal 1 — from repo root
make node-init && make node-start   # enables LCD :1317

# terminal 2
cd web
npm install
# if the host has a low file-watch limit: ulimit -n 65536
npm run dev
# → http://127.0.0.1:5173
```

Default LCD in the UI is `http://127.0.0.1:1317` (direct). `make node-init` turns on the API + CORS and binds `127.0.0.1:1317`.

Optional Vite proxy: set `VITE_LCD_URL=/lcd` and use `/lcd/*` → `VITE_LCD_PROXY_TARGET` (see `.env.example`).

### REST paths (this chain)

Ignite/grpc-gateway annotations use the `/ksedleni/bookstore/...` prefix — **not** `/bookstore/.../v1/...` (those return 501 here).

```bash
# originals (works)
curl -s 'http://127.0.0.1:1317/ksedleni/bookstore/bookstore/original?pagination.limit=100'

# listings (works)
curl -s 'http://127.0.0.1:1317/ksedleni/bookstore/marketplace/listing?pagination.limit=100'
```

## What it shows

- Originals: `id`, `contentHash`, `metadataUri`, `owner`, `creator`
- Active marketplace listings joined by `originalId`
- Copy that this is *a* storefront; CLI/P2P is primary

## Keys / wallet

Localnet: use `bookstored keys … --keyring-backend test` and the CLI for mint/sale/content. Keplr (or similar) can be added later against chain id `bookstore-1` — still non-custodial; not required to sell.

## Build

```bash
npm run build
npm run preview
```
