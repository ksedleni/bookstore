package content

import (
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptAndHashVerify(t *testing.T) {
	priv, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	addr := "store1ownerdemo000000000000000000000000000"

	plain := []byte("Once upon a paperback — unique original bytes.\n")
	env, err := Encrypt(plain, addr, priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, HashPlaintext(plain), env.ContentHash)
	require.Equal(t, WrapScheme, env.KeyWraps[0].Scheme)

	out, err := Decrypt(env, addr, priv)
	require.NoError(t, err)
	require.Equal(t, plain, out)
}

func TestDecryptRejectsTamperedCiphertext(t *testing.T) {
	priv, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	addr := "store1owner"

	env, err := Encrypt([]byte("secret book"), addr, priv.PubKey())
	require.NoError(t, err)
	// Flip a ciphertext byte
	b := []byte(env.Ciphertext)
	if b[len(b)-1] == '0' {
		b[len(b)-1] = '1'
	} else {
		b[len(b)-1] = '0'
	}
	env.Ciphertext = string(b)

	_, err = Decrypt(env, addr, priv)
	require.Error(t, err)
}

func TestVerifyHashMismatch(t *testing.T) {
	require.Error(t, VerifyHash([]byte("a"), HashPlaintext([]byte("b"))))
}

func TestRewrapForNewOwner(t *testing.T) {
	seller, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	buyer, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	sellerAddr := "store1seller"
	buyerAddr := "store1buyer"

	plain := []byte("resale chapter")
	env, err := Encrypt(plain, sellerAddr, seller.PubKey())
	require.NoError(t, err)

	require.NoError(t, Rewrap(env, sellerAddr, seller, buyerAddr, buyer.PubKey()))

	// Seller can still decrypt (wrap retained).
	outSeller, err := Decrypt(env, sellerAddr, seller)
	require.NoError(t, err)
	require.Equal(t, plain, outSeller)

	// Buyer can decrypt after rewrap.
	outBuyer, err := Decrypt(env, buyerAddr, buyer)
	require.NoError(t, err)
	require.Equal(t, plain, outBuyer)
}

func TestStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	st, err := NewStore(dir)
	require.NoError(t, err)

	priv, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	env, err := Encrypt([]byte("stored"), "store1x", priv.PubKey())
	require.NoError(t, err)
	require.NoError(t, st.Put(env))

	got, err := st.Get(env.ContentHash)
	require.NoError(t, err)
	require.Equal(t, env.ContentHash, got.ContentHash)
	out, err := Decrypt(got, "store1x", priv)
	require.NoError(t, err)
	require.Equal(t, []byte("stored"), out)
}
