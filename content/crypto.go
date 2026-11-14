package content

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/hkdf"
)

// Encrypt builds an envelope: AES-256-GCM(plaintext) + CEK wrap for recipient.
// plaintextHash must equal HashPlaintext(plaintext); the digest is stored and
// re-verified on decrypt.
func Encrypt(plaintext []byte, recipientAddr string, recipientPub *btcec.PublicKey) (*Envelope, error) {
	contentHash := HashPlaintext(plaintext)

	cek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, cek); err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, []byte(contentHash))

	wrap, err := WrapCEK(cek, recipientAddr, recipientPub)
	if err != nil {
		return nil, err
	}

	return &Envelope{
		Version:     EnvelopeVersion,
		ContentHash: contentHash,
		Cipher:      CipherAESGCM,
		Nonce:       hex.EncodeToString(nonce),
		Ciphertext:  hex.EncodeToString(ciphertext),
		KeyWraps:    []KeyWrap{wrap},
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// Decrypt recovers plaintext for the recipient private key and verifies the hash.
func Decrypt(env *Envelope, recipientAddr string, recipientPriv *btcec.PrivateKey) ([]byte, error) {
	if env.Version != EnvelopeVersion {
		return nil, fmt.Errorf("unsupported envelope version %d", env.Version)
	}
	if env.Cipher != CipherAESGCM {
		return nil, fmt.Errorf("unsupported cipher %s", env.Cipher)
	}
	wrap, ok := env.FindWrap(recipientAddr)
	if !ok {
		return nil, fmt.Errorf("no key wrap for recipient %s", recipientAddr)
	}
	cek, err := UnwrapCEK(wrap, recipientPriv)
	if err != nil {
		return nil, err
	}

	nonce, err := hex.DecodeString(env.Nonce)
	if err != nil {
		return nil, err
	}
	ciphertext, err := hex.DecodeString(env.Ciphertext)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, []byte(env.ContentHash))
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}
	if err := VerifyHash(plaintext, env.ContentHash); err != nil {
		return nil, err
	}
	return plaintext, nil
}

// Rewrap adds or replaces a CEK wrap for newRecipient using an existing wrap
// held by currentPriv (typically the seller after a P2P sale).
func Rewrap(env *Envelope, currentAddr string, currentPriv *btcec.PrivateKey, newAddr string, newPub *btcec.PublicKey) error {
	cek, err := func() ([]byte, error) {
		wrap, ok := env.FindWrap(currentAddr)
		if !ok {
			return nil, fmt.Errorf("no key wrap for current owner %s", currentAddr)
		}
		return UnwrapCEK(wrap, currentPriv)
	}()
	if err != nil {
		return err
	}
	wrap, err := WrapCEK(cek, newAddr, newPub)
	if err != nil {
		return err
	}
	env.UpsertWrap(wrap)
	return nil
}

// WrapCEK encrypts the CEK to recipientPub via ephemeral secp256k1 ECDH + HKDF + AES-GCM.
func WrapCEK(cek []byte, recipientAddr string, recipientPub *btcec.PublicKey) (KeyWrap, error) {
	ephPriv, err := btcec.NewPrivateKey()
	if err != nil {
		return KeyWrap{}, err
	}
	shared := btcec.GenerateSharedSecret(ephPriv, recipientPub)
	wrapKey, err := deriveWrapKey(shared, recipientPub.SerializeCompressed())
	if err != nil {
		return KeyWrap{}, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return KeyWrap{}, err
	}
	block, err := aes.NewCipher(wrapKey)
	if err != nil {
		return KeyWrap{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return KeyWrap{}, err
	}
	wrapped := gcm.Seal(nil, nonce, cek, nil)
	payload := append(nonce, wrapped...)

	return KeyWrap{
		Recipient:       recipientAddr,
		RecipientPubKey: hex.EncodeToString(recipientPub.SerializeCompressed()),
		Scheme:          WrapScheme,
		EphemeralPubKey: hex.EncodeToString(ephPriv.PubKey().SerializeCompressed()),
		WrappedCEK:      hex.EncodeToString(payload),
	}, nil
}

// UnwrapCEK recovers the CEK using the recipient private key.
func UnwrapCEK(wrap KeyWrap, recipientPriv *btcec.PrivateKey) ([]byte, error) {
	if wrap.Scheme != WrapScheme {
		return nil, fmt.Errorf("unsupported wrap scheme %s", wrap.Scheme)
	}
	ephBytes, err := hex.DecodeString(wrap.EphemeralPubKey)
	if err != nil {
		return nil, err
	}
	ephPub, err := btcec.ParsePubKey(ephBytes)
	if err != nil {
		return nil, err
	}
	shared := btcec.GenerateSharedSecret(recipientPriv, ephPub)
	wrapKey, err := deriveWrapKey(shared, recipientPriv.PubKey().SerializeCompressed())
	if err != nil {
		return nil, err
	}
	payload, err := hex.DecodeString(wrap.WrappedCEK)
	if err != nil {
		return nil, err
	}
	if len(payload) < 12+16 {
		return nil, fmt.Errorf("wrapped CEK too short")
	}
	nonce, ct := payload[:12], payload[12:]
	block, err := aes.NewCipher(wrapKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ct, nil)
}

func deriveWrapKey(sharedSecret, recipientCompressedPub []byte) ([]byte, error) {
	r := hkdf.New(sha256.New, sharedSecret, recipientCompressedPub, []byte("bookstore-content-cek-wrap-v1"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(r, key); err != nil {
		return nil, err
	}
	return key, nil
}

// ParsePrivKeyHex parses a 32-byte secp256k1 private key hex string.
func ParsePrivKeyHex(h string) (*btcec.PrivateKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	if len(b) != 32 {
		return nil, fmt.Errorf("privkey must be 32 bytes")
	}
	priv, _ := btcec.PrivKeyFromBytes(b)
	return priv, nil
}

// ParsePubKeyHex parses a compressed secp256k1 public key hex string.
func ParsePubKeyHex(h string) (*btcec.PublicKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return btcec.ParsePubKey(b)
}
