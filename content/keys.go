package content

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// PrivFromCosmos converts a cosmos secp256k1 privkey to btcec.
func PrivFromCosmos(priv cryptotypes.PrivKey) (*btcec.PrivateKey, error) {
	sk, ok := priv.(*secp256k1.PrivKey)
	if !ok {
		return nil, fmt.Errorf("expected secp256k1 privkey, got %T", priv)
	}
	return ParsePrivKeyHex(hex.EncodeToString(sk.Bytes()))
}

// PubFromCosmos converts a cosmos secp256k1 pubkey to btcec.
func PubFromCosmos(pub cryptotypes.PubKey) (*btcec.PublicKey, error) {
	pk, ok := pub.(*secp256k1.PubKey)
	if !ok {
		return nil, fmt.Errorf("expected secp256k1 pubkey, got %T", pub)
	}
	return btcec.ParsePubKey(pk.Bytes())
}
