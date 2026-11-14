package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

func (k msgServer) CancelListing(goCtx context.Context, msg *types.MsgCancelListing) (*types.MsgCancelListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	listing, found := k.GetListing(ctx, msg.ListingId)
	if !found {
		return nil, types.ErrListingNotFound
	}
	if listing.Seller != msg.Creator {
		return nil, types.ErrNotOwner
	}
	if !listing.Active {
		return nil, types.ErrListingInactive
	}

	listing.Active = false
	k.SetListing(ctx, listing)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"Delisted",
		sdk.NewAttribute("listing_id", strconv.FormatUint(msg.ListingId, 10)),
		sdk.NewAttribute("original_id", strconv.FormatUint(listing.OriginalId, 10)),
	))

	return &types.MsgCancelListingResponse{}, nil
}
