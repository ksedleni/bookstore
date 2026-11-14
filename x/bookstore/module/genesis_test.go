package bookstore_test

import (
	"testing"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	bookstore "github.com/ksedleni/bookstore/x/bookstore/module"
	"github.com/ksedleni/bookstore/x/bookstore/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		OriginalList: []types.Original{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		OriginalCount: 2,
		ContentIndexList: []types.ContentIndex{
			{
				ContentHash: "0",
			},
			{
				ContentHash: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BookstoreKeeper(t)
	bookstore.InitGenesis(ctx, k, genesisState)
	got := bookstore.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.OriginalList, got.OriginalList)
	require.Equal(t, genesisState.OriginalCount, got.OriginalCount)
	require.ElementsMatch(t, genesisState.ContentIndexList, got.ContentIndexList)
	// this line is used by starport scaffolding # genesis/test/assert
}
