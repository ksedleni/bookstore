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

Dev server proxies `/lcd/*` → `VITE_LCD_PROXY_TARGET` (default `http://127.0.0.1:1317`) so the browser avoids CORS. Override in `.env` (see `.env.example`).

In the UI, leave the LCD field as `/lcd` while using the Vite proxy, or set a full LCD URL (and enable CORS on that API if needed).

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
