package keeper

import (
	"fmt"

	"github.com/ksedleni/bookstore/protocol"
)

// BankBridge is the narrow bank interface the ownership module needs.
// Wired to cosmos-sdk bank keeper in a full app; local demo uses protocol.LocalLedger.
type BankBridge interface {
	SendCoins(from, to, amount, denom string) error
}

// Keeper holds BookStore ownership state. Full Cosmos SDK store wiring lands
// when the appchain binary is generated; MVP logic lives in protocol.LocalLedger.
type Keeper struct {
	Ledger *protocol.LocalLedger
	Bank   BankBridge
}

func NewKeeper(ledger *protocol.LocalLedger, bank BankBridge) Keeper {
	return Keeper{Ledger: ledger, Bank: bank}
}

func (k Keeper) Mint(creator, contentHash, metadataHash, metadataURI string) (*protocol.Original, error) {
	return k.Ledger.Mint(creator, contentHash, metadataHash, metadataURI)
}

func (k Keeper) Transfer(from, to, originalID string) error {
	return k.Ledger.Transfer(from, to, originalID)
}

func (k Keeper) P2PSale(buyer, seller, originalID, priceAmount string) (*protocol.SaleSplit, error) {
	split, err := k.Ledger.P2PSale(buyer, seller, originalID, priceAmount)
	if err != nil {
		return nil, err
	}
	if split.ProtocolFee.Sign() != 0 {
		return nil, fmt.Errorf("invariant violated: protocol fee must be 0")
	}
	return split, nil
}
