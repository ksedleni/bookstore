package app

import (
	bookstorekeeper "github.com/ksedleni/bookstore/chain/x/bookstore/keeper"
	marketkeeper "github.com/ksedleni/bookstore/chain/x/marketplace/keeper"
	"github.com/ksedleni/bookstore/protocol"
)

// Name is the sovereign BookStore appchain identifier.
const Name = "bookstore"

// App is the MVP application wiring for the BookStore sovereign ledger.
// Full Cosmos SDK / CometBFT node scaffolding (genesis, ABCI, WASM hooks)
// is the next spike; this wires domain keepers for local verification.
type App struct {
	Ledger      *protocol.LocalLedger
	Bookstore   bookstorekeeper.Keeper
	Marketplace marketkeeper.Keeper
}

// NewApp constructs the MVP app with an in-process ledger.
func NewApp() *App {
	ledger := protocol.NewLocalLedger()
	return &App{
		Ledger:      ledger,
		Bookstore:   bookstorekeeper.NewKeeper(ledger, nil),
		Marketplace: marketkeeper.NewKeeper(ledger),
	}
}
