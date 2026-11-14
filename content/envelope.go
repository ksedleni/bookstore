package content

import "time"

const (
	// EnvelopeVersion is the portable schema version.
	EnvelopeVersion = 1
	// WrapScheme identifies the CEK wrap algorithm for third-party clients.
	WrapScheme = "secp256k1-ecdh-hkdf-v1"
	// CipherAESGCM identifies the content cipher.
	CipherAESGCM = "AES-256-GCM"
)

// KeyWrap binds a CEK to one recipient secp256k1 pubkey / address.
type KeyWrap struct {
	Recipient        string `json:"recipient"`          // bech32 address
	RecipientPubKey  string `json:"recipient_pubkey"`   // compressed secp256k1 hex
	Scheme           string `json:"scheme"`             // WrapScheme
	EphemeralPubKey  string `json:"ephemeral_pubkey"`   // compressed secp256k1 hex
	WrappedCEK       string `json:"wrapped_cek"`        // hex(nonce||ciphertext||tag) of CEK
}

// Envelope is the portable off-chain object for one contentHash.
type Envelope struct {
	Version     int       `json:"version"`
	ContentHash string    `json:"content_hash"` // SHA-256 of plaintext
	Cipher      string    `json:"cipher"`
	Nonce       string    `json:"nonce"`      // hex, 12 bytes
	Ciphertext  string    `json:"ciphertext"` // hex
	KeyWraps    []KeyWrap `json:"key_wraps"`
	CreatedAt   time.Time `json:"created_at"`
	// OriginalID is optional bookkeeping for local demos (not authoritative).
	OriginalID string `json:"original_id,omitempty"`
}

// FindWrap returns the wrap for recipient address, if any.
func (e *Envelope) FindWrap(recipient string) (KeyWrap, bool) {
	for _, w := range e.KeyWraps {
		if w.Recipient == recipient {
			return w, true
		}
	}
	return KeyWrap{}, false
}

// UpsertWrap replaces or appends a wrap for the recipient.
func (e *Envelope) UpsertWrap(w KeyWrap) {
	for i := range e.KeyWraps {
		if e.KeyWraps[i].Recipient == w.Recipient {
			e.KeyWraps[i] = w
			return
		}
	}
	e.KeyWraps = append(e.KeyWraps, w)
}
