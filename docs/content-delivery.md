# Content delivery (owner-gated, portable)

BookStore keeps **file bytes off-chain**. The ledger stores only the **SHA-256** digest (`contentHash`). Delivery is portable: any host may store/serve envelopes; any compliant client may decrypt if it holds the CEK wrap for the current owner.

## Envelope format (v1)

```json
{
  "version": 1,
  "content_hash": "<sha256 hex of plaintext bytes>",
  "cipher": "AES-256-GCM",
  "nonce": "<12-byte nonce hex>",
  "ciphertext": "<AES-GCM ciphertext+tag hex>",
  "key_wraps": [
    {
      "recipient": "<bech32 store1...>",
      "recipient_pubkey": "<compressed secp256k1 hex>",
      "scheme": "secp256k1-ecdh-hkdf-v1",
      "ephemeral_pubkey": "<compressed secp256k1 hex>",
      "wrapped_cek": "<nonce||aes-gcm(cek) hex>"
    }
  ]
}
```

### Algorithms

1. **Content commitment:** `SHA-256(plaintext)` — must match on-chain `contentHash`.
2. **Content encryption:** random 32-byte CEK; `AES-256-GCM` with AAD = `content_hash` bytes (hex string as used in Go AAD).
3. **CEK wrap (`secp256k1-ecdh-hkdf-v1`):**
   - Ephemeral secp256k1 keypair
   - `shared = ECDH(eph_priv, recipient_pub)` (`btcec.GenerateSharedSecret`)
   - `wrap_key = HKDF-SHA256(shared, salt=recipient_compressed_pubkey, info="bookstore-content-cek-wrap-v1")`
   - Encrypt CEK with AES-256-GCM under `wrap_key`

On download, clients **must** recompute SHA-256 and reject mismatches.

## Local store layout

Default directory: `~/.bookstore/content/<content_hash>.json` (override with `--content-dir`).

This is an MVP stand-in for IPFS/Arweave/HTTP mirrors. Third parties can copy the same JSON objects anywhere.

## CLI

```bash
# After (or before) mint — encrypt for publisher key
.bin/bookstored content put --from alice --file ./book.epub --home ~/.bookstore --keyring-backend test
# → prints content_hash (use that for mint-original)

# Fetch as current owner (verifies hash)
.bin/bookstored content get --from alice --content-hash <hash> --out ./book.epub \
  --home ~/.bookstore --keyring-backend test

# Or resolve hash from chain original id (also checks ledger owner)
.bin/bookstored content get --from bob --original-id 1 --out ./book.epub \
  --home ~/.bookstore --keyring-backend test --node tcp://127.0.0.1:26657

# After P2P sale: seller re-wraps CEK for buyer (bookstore-less)
.bin/bookstored content rewrap --from alice --to bob --content-hash <hash> \
  --home ~/.bookstore --keyring-backend test
```

## Ownership model

- **Encryption** is the portable access control (ciphertext can be public).
- **Gated hosts** should additionally check ledger ownership before serving.
- `content get --original-id` checks the on-chain owner when `--require-owner` is true (default).
- Sale does **not** automatically move CEK wraps — call `rewrap` (offline ceremony between seller and buyer keys). Automatic sale-time rewrap is deferred.

## Deferred

- IPFS/Arweave pinning helpers
- HTTP content gateway
- Automatic rewrap inside `MsgP2PSale`
- Threshold / social recovery wraps
