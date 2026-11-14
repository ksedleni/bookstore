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

Smoke (CI-friendly):

```bash
make smoke-node
```

## Useful queries (once running)

```bash
.bin/bookstored status
.bin/bookstored query bookstore list-original
.bin/bookstored query marketplace list-listing
.bin/bookstored query marketplace show-pause-state
```

## Invariants

See [invariants.md](./invariants.md). Marketplace pause never blocks `P2PSale` or `MintOriginal`.
