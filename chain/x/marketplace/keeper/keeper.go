package keeper

import (
	"github.com/ksedleni/bookstore/protocol"
)

// Keeper for optional on-ledger listings. Pause affects only this module.
type Keeper struct {
	Ledger *protocol.LocalLedger
}

func NewKeeper(ledger *protocol.LocalLedger) Keeper {
	return Keeper{Ledger: ledger}
}

func (k Keeper) List(seller, originalID, price, denom string) (*protocol.Listing, error) {
	return k.Ledger.List(seller, originalID, price, denom)
}

func (k Keeper) Delist(seller, listingID string) error {
	return k.Ledger.Delist(seller, listingID)
}

func (k Keeper) Buy(buyer, listingID string) (*protocol.SaleSplit, error) {
	return k.Ledger.BuyListed(buyer, listingID)
}

func (k Keeper) SetPaused(paused bool) {
	k.Ledger.SetMarketplacePaused(paused)
}
