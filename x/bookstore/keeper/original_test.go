package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	"github.com/ksedleni/bookstore/x/bookstore/keeper"
	"github.com/ksedleni/bookstore/x/bookstore/types"
	"github.com/stretchr/testify/require"
)

func createNOriginal(keeper keeper.Keeper, ctx context.Context, n int) []types.Original {
	items := make([]types.Original, n)
	for i := range items {
		items[i].Id = keeper.AppendOriginal(ctx, items[i])
	}
	return items
}

func TestOriginalGet(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	items := createNOriginal(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetOriginal(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestOriginalRemove(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	items := createNOriginal(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveOriginal(ctx, item.Id)
		_, found := keeper.GetOriginal(ctx, item.Id)
		require.False(t, found)
	}
}

func TestOriginalGetAll(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	items := createNOriginal(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOriginal(ctx)),
	)
}

func TestOriginalCount(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	items := createNOriginal(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetOriginalCount(ctx))
}
