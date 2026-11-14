# BookStore reference storefront

Optional UI for publishers and buyers. **The protocol works without it** (CLI / P2P). Protocol fee **0%**. Creator royalty **0.1%** on paid sales.

## Run (publisher-friendly)

Three terminals from the repo root:

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

Open http://127.0.0.1:5173

You should see **Connected** when the ledger API answers. The publisher helper (http://127.0.0.1:8787) signs localnet actions with the test keyring (`alice` / `bob`) — non-custodial local keys, not a hosted wallet. Keplr can replace the helper later.

### Publish two editions and list them

1. Open **Publish a book**
2. Choose a file
3. Set **Numbered editions** to `2`
4. Keep **List for sale** on and set a price (default `1000000` stake)
5. Click **Mint & list**

Each edition is a distinct file fingerprint and original id (edition bytes append a small footer so hashes differ). Both appear under **Books for sale**.

### Buy as the localnet buyer

On **Books for sale**, click **Buy**. The helper settles `buy-listing` as `bob`, rewraps content from `alice` → `bob`, and writes the decrypted file under `/tmp`.

## Advanced

Ledger fingerprint / LCD details stay under **Advanced / ledger details**.

Direct LCD (CORS enabled by `make node-init`):

```bash
curl -s 'http://127.0.0.1:1317/ksedleni/bookstore/bookstore/original?pagination.limit=100'
```

## Build

```bash
cd web
npm install
npm run build
```

## Smoke

With a node already running (`make node-start`):

```bash
make publisher-smoke
```
