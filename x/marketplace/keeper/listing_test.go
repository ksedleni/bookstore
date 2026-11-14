package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	"github.com/ksedleni/bookstore/x/marketplace/keeper"
	"github.com/ksedleni/bookstore/x/marketplace/types"
	"github.com/stretchr/testify/require"
)

func createNListing(keeper keeper.Keeper, ctx context.Context, n int) []types.Listing {
	items := make([]types.Listing, n)
	for i := range items {
		items[i].Id = keeper.AppendListing(ctx, items[i])
	}
	return items
}

func TestListingGet(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNListing(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetListing(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestListingRemove(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNListing(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveListing(ctx, item.Id)
		_, found := keeper.GetListing(ctx, item.Id)
		require.False(t, found)
	}
}

func TestListingGetAll(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNListing(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllListing(ctx)),
	)
}

func TestListingCount(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	items := createNListing(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetListingCount(ctx))
}
