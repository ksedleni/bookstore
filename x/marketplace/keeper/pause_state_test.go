package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	"github.com/ksedleni/bookstore/x/marketplace/keeper"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

func createTestPauseState(keeper keeper.Keeper, ctx context.Context) types.PauseState {
	item := types.PauseState{}
	keeper.SetPauseState(ctx, item)
	return item
}

func TestPauseStateGet(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	item := createTestPauseState(keeper, ctx)
	rst, found := keeper.GetPauseState(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestPauseStateRemove(t *testing.T) {
	keeper, ctx := keepertest.MarketplaceKeeper(t)
	createTestPauseState(keeper, ctx)
	keeper.RemovePauseState(ctx)
	_, found := keeper.GetPauseState(ctx)
	require.False(t, found)
}
