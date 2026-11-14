package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

// GetOriginalCount get the total number of original
func (k Keeper) GetOriginalCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.OriginalCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetOriginalCount set the total number of original
func (k Keeper) SetOriginalCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.OriginalCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendOriginal appends an original with a new 1-based id and updates the count.
// The first mint receives id=1 (not 0). Count remains the number of appends.
func (k Keeper) AppendOriginal(
	ctx context.Context,
	original types.Original,
) uint64 {
	count := k.GetOriginalCount(ctx)
	id := count + 1 // 1-based ids
	original.Id = id

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OriginalKey))
	appendedValue := k.cdc.MustMarshal(&original)
	store.Set(GetOriginalIDBytes(original.Id), appendedValue)

	k.SetOriginalCount(ctx, id)
	return id
}

// SetOriginal set a specific original in the store
func (k Keeper) SetOriginal(ctx context.Context, original types.Original) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OriginalKey))
	b := k.cdc.MustMarshal(&original)
	store.Set(GetOriginalIDBytes(original.Id), b)
}

// GetOriginal returns a original from its id
func (k Keeper) GetOriginal(ctx context.Context, id uint64) (val types.Original, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OriginalKey))
	b := store.Get(GetOriginalIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOriginal removes a original from the store
func (k Keeper) RemoveOriginal(ctx context.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OriginalKey))
	store.Delete(GetOriginalIDBytes(id))
}

// GetAllOriginal returns all original
func (k Keeper) GetAllOriginal(ctx context.Context) (list []types.Original) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OriginalKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Original
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetOriginalIDBytes returns the byte representation of the ID
func GetOriginalIDBytes(id uint64) []byte {
	bz := types.KeyPrefix(types.OriginalKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, id)
	return bz
}
