package keeper

import (
	"context"
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

func (k msgServer) MintOriginal(goCtx context.Context, msg *types.MsgMintOriginal) (*types.MsgMintOriginalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.ContentHash) != 64 {
		return nil, types.ErrInvalidContentHash
	}
	if _, err := hex.DecodeString(msg.ContentHash); err != nil {
		return nil, types.ErrInvalidContentHash
	}

	if _, found := k.GetContentIndex(ctx, msg.ContentHash); found {
		return nil, types.ErrDuplicateContent
	}

	original := types.Original{
		ContentHash:  msg.ContentHash,
		MetadataHash: msg.MetadataHash,
		MetadataUri:  msg.MetadataUri,
		Owner:        msg.Creator,
		RoyaltyPayee: msg.Creator,
		Creator:      msg.Creator,
	}
	id := k.AppendOriginal(ctx, original)
	k.SetContentIndex(ctx, types.ContentIndex{
		ContentHash: msg.ContentHash,
		OriginalId:  id,
	})

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"Minted",
		sdk.NewAttribute("original_id", strconv.FormatUint(id, 10)),
		sdk.NewAttribute("content_hash", msg.ContentHash),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("owner", msg.Creator),
	))

	return &types.MsgMintOriginalResponse{Id: id}, nil
}
