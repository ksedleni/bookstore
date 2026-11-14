package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

// SetContentIndex set a specific contentIndex in the store from its index
func (k Keeper) SetContentIndex(ctx context.Context, contentIndex types.ContentIndex) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ContentIndexKeyPrefix))
	b := k.cdc.MustMarshal(&contentIndex)
	store.Set(types.ContentIndexKey(
		contentIndex.ContentHash,
	), b)
}

// GetContentIndex returns a contentIndex from its index
func (k Keeper) GetContentIndex(
	ctx context.Context,
	contentHash string,

) (val types.ContentIndex, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ContentIndexKeyPrefix))

	b := store.Get(types.ContentIndexKey(
		contentHash,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveContentIndex removes a contentIndex from the store
func (k Keeper) RemoveContentIndex(
	ctx context.Context,
	contentHash string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ContentIndexKeyPrefix))
	store.Delete(types.ContentIndexKey(
		contentHash,
	))
}

// GetAllContentIndex returns all contentIndex
func (k Keeper) GetAllContentIndex(ctx context.Context) (list []types.ContentIndex) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ContentIndexKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ContentIndex
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
