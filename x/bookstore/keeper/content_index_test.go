package keeper_test

import (
	"context"
	"strconv"
	"testing"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	"github.com/ksedleni/bookstore/x/bookstore/keeper"
	"github.com/ksedleni/bookstore/x/bookstore/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNContentIndex(keeper keeper.Keeper, ctx context.Context, n int) []types.ContentIndex {
	items := make([]types.ContentIndex, n)
	for i := range items {
		items[i].ContentHash = strconv.Itoa(i)

		keeper.SetContentIndex(ctx, items[i])
	}
	return items
}

func TestContentIndexGet(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	items := createNContentIndex(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetContentIndex(ctx,
			item.ContentHash,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestContentIndexRemove(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	items := createNContentIndex(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveContentIndex(ctx,
			item.ContentHash,
		)
		_, found := keeper.GetContentIndex(ctx,
			item.ContentHash,
		)
		require.False(t, found)
	}
}

func TestContentIndexGetAll(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	items := createNContentIndex(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllContentIndex(ctx)),
	)
}
