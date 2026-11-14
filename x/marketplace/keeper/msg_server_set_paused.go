package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

// SetPaused toggles the narrow marketplace pause. Does not block P2P/mint.
func (k msgServer) SetPaused(goCtx context.Context, msg *types.MsgSetPaused) (*types.MsgSetPausedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Creator != k.GetAuthority() {
		return nil, types.ErrUnauthorized
	}

	k.SetPauseState(ctx, types.PauseState{
		Paused:  msg.Paused,
	})

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"MarketplacePause",
		sdk.NewAttribute("paused", strconv.FormatBool(msg.Paused)),
	))

	return &types.MsgSetPausedResponse{}, nil
}
