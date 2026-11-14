// Package content implements portable, owner-gated off-chain content delivery
// for BookStore Originals.
//
// The chain stores only the SHA-256 content digest. File bytes live off-ledger.
// Access control is portable: ciphertext may be mirrored by any host; only a
// party that can unwrap the content-encryption key (CEK) for their secp256k1
// key can decrypt. Third-party clients should:
//
//  1. Fetch the envelope (local store, HTTP, IPFS, …).
//  2. Confirm ledger ownership for the caller's address (optional for decrypt,
//     required for honest gated hosts).
//  3. Unwrap the CEK with the owner private key (secp256k1-ecdh-hkdf-v1).
//  4. AES-256-GCM decrypt the ciphertext.
//  5. Recompute SHA-256 and reject if it does not equal content_hash / on-chain digest.
//
// After a P2P sale, run Rewrap so the new owner receives a CEK wrap (bookstore-less;
// no protocol fee).
package content
