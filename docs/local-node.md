# Local BookStore node

## Init + start

```bash
make build
make node-init    # ~/.bookstore — wipe + genesis + gentx
make node-start
```

Or one shot:

```bash
make node
```

Smoke / e2e:

```bash
make smoke-node   # start until a block commits
make e2e          # mint → P2P sale on a throwaway localnet
```

## Useful queries (once running)

```bash
.bin/bookstored status
.bin/bookstored query bookstore list-original      # ids are 1-based
.bin/bookstored query bookstore show-original 1
.bin/bookstored query marketplace list-listing
.bin/bookstored query marketplace show-pause-state
```

## Content hash

`mint-original` takes the **SHA-256 hex digest of file bytes**. The chain stores that digest only (plus metadata hash/URI). See the README walkthrough for a full mint → `p2p-sale` flow (`chain-id=bookstore-1`, keyring `test`, wait for `q tx` before querying).

## Invariants

See [invariants.md](./invariants.md). Marketplace pause never blocks `P2PSale` or `MintOriginal`. Royalty **0.1%**, protocol fee **0%**.
