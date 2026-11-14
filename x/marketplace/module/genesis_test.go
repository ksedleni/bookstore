package marketplace_test

import (
	"testing"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	marketplace "github.com/ksedleni/bookstore/x/marketplace/module"
	"github.com/ksedleni/bookstore/x/marketplace/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		ListingList: []types.Listing{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		ListingCount: 2,
		PauseState: &types.PauseState{
			Paused: false,
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MarketplaceKeeper(t)
	marketplace.InitGenesis(ctx, k, genesisState)
	got := marketplace.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.ListingList, got.ListingList)
	require.Equal(t, genesisState.ListingCount, got.ListingCount)
	require.Equal(t, genesisState.PauseState, got.PauseState)
	// this line is used by starport scaffolding # genesis/test/assert
}
