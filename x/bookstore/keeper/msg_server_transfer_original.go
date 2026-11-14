package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

func (k msgServer) TransferOriginal(goCtx context.Context, msg *types.MsgTransferOriginal) (*types.MsgTransferOriginalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.To == "" {
		return nil, types.ErrInvalidSale
	}
	original, found := k.GetOriginal(ctx, msg.OriginalId)
	if !found {
		return nil, types.ErrOriginalNotFound
	}
	if original.Owner != msg.Creator {
		return nil, types.ErrNotOwner
	}

	from := original.Owner
	original.Owner = msg.To
	k.SetOriginal(ctx, original)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"Transfer",
		sdk.NewAttribute("original_id", strconv.FormatUint(msg.OriginalId, 10)),
		sdk.NewAttribute("from", from),
		sdk.NewAttribute("to", msg.To),
	))

	return &types.MsgTransferOriginalResponse{}, nil
}
