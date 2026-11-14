package content

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HashPlaintext returns the canonical SHA-256 hex digest of file bytes.
func HashPlaintext(plaintext []byte) string {
	sum := sha256.Sum256(plaintext)
	return hex.EncodeToString(sum[:])
}

// VerifyHash checks that plaintext matches the expected hex SHA-256 digest.
func VerifyHash(plaintext []byte, expectedHex string) error {
	got := HashPlaintext(plaintext)
	if got != expectedHex {
		return fmt.Errorf("content hash mismatch: got %s want %s", got, expectedHex)
	}
	return nil
}
