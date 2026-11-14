package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ksedleni/bookstore/x/marketplace/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListingAll(ctx context.Context, req *types.QueryAllListingRequest) (*types.QueryAllListingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var listings []types.Listing

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	listingStore := prefix.NewStore(store, types.KeyPrefix(types.ListingKey))

	pageRes, err := query.Paginate(listingStore, req.Pagination, func(key []byte, value []byte) error {
		var listing types.Listing
		if err := k.cdc.Unmarshal(value, &listing); err != nil {
			return err
		}

		listings = append(listings, listing)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllListingResponse{Listing: listings, Pagination: pageRes}, nil
}

func (k Keeper) Listing(ctx context.Context, req *types.QueryGetListingRequest) (*types.QueryGetListingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	listing, found := k.GetListing(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetListingResponse{Listing: listing}, nil
}
