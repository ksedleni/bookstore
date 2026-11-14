package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

func TestOriginalQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	msgs := createNOriginal(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetOriginalRequest
		response *types.QueryGetOriginalResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetOriginalRequest{Id: msgs[0].Id},
			response: &types.QueryGetOriginalResponse{Original: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetOriginalRequest{Id: msgs[1].Id},
			response: &types.QueryGetOriginalResponse{Original: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetOriginalRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Original(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestOriginalQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	msgs := createNOriginal(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllOriginalRequest {
		return &types.QueryAllOriginalRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.OriginalAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Original), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Original),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.OriginalAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Original), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Original),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.OriginalAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Original),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.OriginalAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
