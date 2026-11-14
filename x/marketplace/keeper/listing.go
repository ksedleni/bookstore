package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

// GetListingCount get the total number of listing
func (k Keeper) GetListingCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.ListingCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetListingCount set the total number of listing
func (k Keeper) SetListingCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.ListingCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendListing appends a listing in the store with a new id and update the count
func (k Keeper) AppendListing(
	ctx context.Context,
	listing types.Listing,
) uint64 {
	// Create the listing
	count := k.GetListingCount(ctx)

	// Set the ID of the appended value
	listing.Id = count

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ListingKey))
	appendedValue := k.cdc.MustMarshal(&listing)
	store.Set(GetListingIDBytes(listing.Id), appendedValue)

	// Update listing count
	k.SetListingCount(ctx, count+1)

	return count
}

// SetListing set a specific listing in the store
func (k Keeper) SetListing(ctx context.Context, listing types.Listing) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ListingKey))
	b := k.cdc.MustMarshal(&listing)
	store.Set(GetListingIDBytes(listing.Id), b)
}

// GetListing returns a listing from its id
func (k Keeper) GetListing(ctx context.Context, id uint64) (val types.Listing, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ListingKey))
	b := store.Get(GetListingIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveListing removes a listing from the store
func (k Keeper) RemoveListing(ctx context.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ListingKey))
	store.Delete(GetListingIDBytes(id))
}

// GetAllListing returns all listing
func (k Keeper) GetAllListing(ctx context.Context) (list []types.Listing) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ListingKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Listing
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetListingIDBytes returns the byte representation of the ID
func GetListingIDBytes(id uint64) []byte {
	bz := types.KeyPrefix(types.ListingKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, id)
	return bz
}
