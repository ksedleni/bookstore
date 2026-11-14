package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/ksedleni/bookstore/testutil/keeper"
	"github.com/ksedleni/bookstore/testutil/nullify"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestContentIndexQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	msgs := createNContentIndex(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetContentIndexRequest
		response *types.QueryGetContentIndexResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetContentIndexRequest{
				ContentHash: msgs[0].ContentHash,
			},
			response: &types.QueryGetContentIndexResponse{ContentIndex: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetContentIndexRequest{
				ContentHash: msgs[1].ContentHash,
			},
			response: &types.QueryGetContentIndexResponse{ContentIndex: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetContentIndexRequest{
				ContentHash: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ContentIndex(ctx, tc.request)
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

func TestContentIndexQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.BookstoreKeeper(t)
	msgs := createNContentIndex(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllContentIndexRequest {
		return &types.QueryAllContentIndexRequest{
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
			resp, err := keeper.ContentIndexAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ContentIndex), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ContentIndex),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ContentIndexAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ContentIndex), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ContentIndex),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ContentIndexAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ContentIndex),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ContentIndexAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
