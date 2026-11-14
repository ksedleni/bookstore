package protocol

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitSale_PointOnePercent(t *testing.T) {
	price := big.NewInt(1_000_000)
	split, err := SplitSale(price)
	require.NoError(t, err)
	require.Equal(t, "1000", split.Royalty.String()) // 0.1% of 1_000_000
	require.Equal(t, "999000", split.SellerProceeds.String())
	require.Equal(t, "0", split.ProtocolFee.String())
	require.Equal(t, uint64(0), ProtocolFeeBPS)
	require.Equal(t, uint64(10), RoyaltyBPS)
}

func TestSplitSale_RejectsNonPositive(t *testing.T) {
	_, err := SplitSale(big.NewInt(0))
	require.Error(t, err)
}

func TestLocalLedger_P2PWithoutBookstore(t *testing.T) {
	led := NewLocalLedger()
	require.NoError(t, led.Faucet("alice", "1000000"))
	require.NoError(t, led.Faucet("bob", "2000000"))

	hash := ContentHashSHA256([]byte("hello paperback"))
	orig, err := led.Mint("alice", hash, "meta", "")
	require.NoError(t, err)

	// Marketplace pause must not block P2P sale.
	led.SetMarketplacePaused(true)
	_, err = led.List("alice", orig.ID, "1000000", "ustore")
	require.Error(t, err)

	split, err := led.P2PSale("bob", "alice", orig.ID, "1000000")
	require.NoError(t, err)
	require.Equal(t, "1000", split.Royalty.String())
	require.Equal(t, "0", split.ProtocolFee.String())

	got, err := led.GetOriginal(orig.ID)
	require.NoError(t, err)
	require.Equal(t, "bob", got.Owner)

	// Resell without re-approval.
	require.NoError(t, led.Faucet("carol", "2000000"))
	_, err = led.P2PSale("carol", "bob", orig.ID, "500000")
	require.NoError(t, err)
	got, err = led.GetOriginal(orig.ID)
	require.NoError(t, err)
	require.Equal(t, "carol", got.Owner)
}

func TestLocalLedger_DuplicateHashRejected(t *testing.T) {
	led := NewLocalLedger()
	hash := ContentHashSHA256([]byte("same bytes"))
	_, err := led.Mint("alice", hash, "m", "")
	require.NoError(t, err)
	_, err = led.Mint("bob", hash, "m2", "")
	require.Error(t, err)
}

func TestLocalLedger_ListedBuy(t *testing.T) {
	led := NewLocalLedger()
	require.NoError(t, led.Faucet("buyer", "5000000"))
	hash := ContentHashSHA256([]byte("listed book"))
	orig, err := led.Mint("seller", hash, "m", "")
	require.NoError(t, err)
	listing, err := led.List("seller", orig.ID, "2000000", "ustore")
	require.NoError(t, err)
	split, err := led.BuyListed("buyer", listing.ID)
	require.NoError(t, err)
	require.Equal(t, "2000", split.Royalty.String())
	got, _ := led.GetOriginal(orig.ID)
	require.Equal(t, "buyer", got.Owner)
}
