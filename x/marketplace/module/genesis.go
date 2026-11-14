package marketplace

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ksedleni/bookstore/x/marketplace/keeper"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the listing
	for _, elem := range genState.ListingList {
		k.SetListing(ctx, elem)
	}

	// Set listing count
	k.SetListingCount(ctx, genState.ListingCount)
	// Set if defined
	if genState.PauseState != nil {
		k.SetPauseState(ctx, *genState.PauseState)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ListingList = k.GetAllListing(ctx)
	genesis.ListingCount = k.GetListingCount(ctx)
	// Get all pauseState
	pauseState, found := k.GetPauseState(ctx)
	if found {
		genesis.PauseState = &pauseState
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
