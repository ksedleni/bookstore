package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// ContentHashSHA256 returns the canonical content commitment for file bytes.
func ContentHashSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// Metadata describes off-ledger book metadata (content-addressed separately).
type Metadata struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	MIMEType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
	CoverURI    string `json:"cover_uri,omitempty"`
	LicenseURI  string `json:"license_uri,omitempty"`
	ContentHash string `json:"content_hash"`
	StorageURIs []string `json:"storage_uris,omitempty"`
}

// Original is a unique BookStore ownership record for one content hash.
type Original struct {
	ID             string    `json:"id"`
	ContentHash    string    `json:"content_hash"`
	MetadataHash   string    `json:"metadata_hash"`
	MetadataURI    string    `json:"metadata_uri,omitempty"`
	Creator        string    `json:"creator"`
	Owner          string    `json:"owner"`
	RoyaltyPayee   string    `json:"royalty_payee"`
	CreatedAt      time.Time `json:"created_at"`
}

// Listing is an on-ledger offer by the current owner.
type Listing struct {
	ID         string    `json:"id"`
	OriginalID string    `json:"original_id"`
	Seller     string    `json:"seller"`
	PriceAmount string   `json:"price_amount"` // decimal integer string in native denom subunits
	Denom      string    `json:"denom"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

// EventType names for open indexers / third-party storefronts.
const (
	EventMinted   = "Minted"
	EventTransfer = "Transfer"
	EventSale     = "Sale"
	EventListed   = "Listed"
	EventDelisted = "Delisted"
)

// Event is a portable ledger event (source of truth for catalogs).
type Event struct {
	Type       string            `json:"type"`
	Height     int64             `json:"height,omitempty"`
	Timestamp  time.Time         `json:"timestamp"`
	OriginalID string            `json:"original_id,omitempty"`
	Attrs      map[string]string `json:"attrs,omitempty"`
}

// ValidateContentHash checks hex-encoded SHA-256.
func ValidateContentHash(h string) error {
	if len(h) != 64 {
		return fmt.Errorf("content hash must be 64 hex chars (sha-256), got len=%d", len(h))
	}
	if _, err := hex.DecodeString(h); err != nil {
		return fmt.Errorf("content hash must be hex: %w", err)
	}
	return nil
}
