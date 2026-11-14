package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

// SetPauseState set pauseState in the store
func (k Keeper) SetPauseState(ctx context.Context, pauseState types.PauseState) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PauseStateKey))
	b := k.cdc.MustMarshal(&pauseState)
	store.Set([]byte{0}, b)
}

// GetPauseState returns pauseState
func (k Keeper) GetPauseState(ctx context.Context) (val types.PauseState, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PauseStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePauseState removes pauseState from the store
func (k Keeper) RemovePauseState(ctx context.Context) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PauseStateKey))
	store.Delete([]byte{0})
}
