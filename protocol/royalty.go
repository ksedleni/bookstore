// Package protocol defines BookStore ownership and sale rules independent of
// any incumbent public chain (Ethereum, Solana, Base, etc.).
package protocol

import (
	"fmt"
	"math/big"
)

// RoyaltyBPS is the fixed creator royalty on paid ledger sales: 0.1% = 10 bps.
// Not seller-configurable.
const RoyaltyBPS uint64 = 10

// ProtocolFeeBPS is always zero. BookStore is not a fee-taking store.
const ProtocolFeeBPS uint64 = 0

// BPSDenominator is 10_000 basis points = 100%.
const BPSDenominator uint64 = 10_000

// SaleSplit is the settlement outcome of a paid sale.
type SaleSplit struct {
	Price         *big.Int
	Royalty       *big.Int // to creator / royalty payee
	SellerProceeds *big.Int
	ProtocolFee   *big.Int // always zero
}

// SplitSale applies the locked fee model: 0.1% creator royalty, 0% protocol fee.
func SplitSale(price *big.Int) (SaleSplit, error) {
	if price == nil || price.Sign() <= 0 {
		return SaleSplit{}, fmt.Errorf("price must be positive")
	}

	royalty := new(big.Int).Mul(price, big.NewInt(int64(RoyaltyBPS)))
	royalty.Div(royalty, big.NewInt(int64(BPSDenominator)))

	// Round-down royalty; seller receives the remainder. Protocol fee is 0.
	seller := new(big.Int).Sub(price, royalty)
	protocol := big.NewInt(0)

	if ProtocolFeeBPS != 0 {
		return SaleSplit{}, fmt.Errorf("protocol fee must remain 0%% (got %d bps)", ProtocolFeeBPS)
	}

	return SaleSplit{
		Price:          new(big.Int).Set(price),
		Royalty:        royalty,
		SellerProceeds: seller,
		ProtocolFee:    protocol,
	}, nil
}
