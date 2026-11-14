package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ksedleni/bookstore/x/bookstore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OriginalAll(ctx context.Context, req *types.QueryAllOriginalRequest) (*types.QueryAllOriginalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var originals []types.Original

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	originalStore := prefix.NewStore(store, types.KeyPrefix(types.OriginalKey))

	pageRes, err := query.Paginate(originalStore, req.Pagination, func(key []byte, value []byte) error {
		var original types.Original
		if err := k.cdc.Unmarshal(value, &original); err != nil {
			return err
		}

		originals = append(originals, original)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOriginalResponse{Original: originals, Pagination: pageRes}, nil
}

func (k Keeper) Original(ctx context.Context, req *types.QueryGetOriginalRequest) (*types.QueryGetOriginalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	original, found := k.GetOriginal(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetOriginalResponse{Original: original}, nil
}
