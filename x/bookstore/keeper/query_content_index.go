package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ksedleni/bookstore/x/bookstore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ContentIndexAll(ctx context.Context, req *types.QueryAllContentIndexRequest) (*types.QueryAllContentIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var contentIndexs []types.ContentIndex

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	contentIndexStore := prefix.NewStore(store, types.KeyPrefix(types.ContentIndexKeyPrefix))

	pageRes, err := query.Paginate(contentIndexStore, req.Pagination, func(key []byte, value []byte) error {
		var contentIndex types.ContentIndex
		if err := k.cdc.Unmarshal(value, &contentIndex); err != nil {
			return err
		}

		contentIndexs = append(contentIndexs, contentIndex)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllContentIndexResponse{ContentIndex: contentIndexs, Pagination: pageRes}, nil
}

func (k Keeper) ContentIndex(ctx context.Context, req *types.QueryGetContentIndexRequest) (*types.QueryGetContentIndexResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetContentIndex(
		ctx,
		req.ContentHash,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetContentIndexResponse{ContentIndex: val}, nil
}
