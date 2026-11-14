# BookStore protocol invariants (MVP)

Locked decisions from the architecture plan:

1. **Protocol fee = 0%.** No marketplace take rate may be introduced without an explicit governance/ADR change.
2. **Creator royalty = 0.1% (10 bps)** on paid ledger sales only (P2P or listed). Unpaid transfers: no royalty. Not seller-configurable.
3. **Permissionless mint/list.** No seller allowlist in protocol.
4. **P2P sale/transfer is first-class** and must work without any bookstore UI.
5. **Marketplace pause is narrow** — listings/buy-listed only; never blocks mint or P2P.
6. **No incumbent-chain dependency** for core trade (no required Ethereum / Solana / Base RPCs).
7. **One mint per content hash** (SHA-256 of exact bytes) in registry scope. The ledger stores the digest only — not file contents.
8. **Original ids are 1-based** (first mint is `id=1`).
9. **Content delivery is off-chain and portable.** Clients verify SHA-256 on download. Owner access is CEK-wrap gated (`secp256k1-ecdh-hkdf-v1`); rewrap after sale is bookstore-less (no protocol fee).
10. **Storefronts are optional.** The reference `web/` UI is a thin client; P2P CLI trade remains valid without any bookstore UI.
