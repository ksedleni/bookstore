package keeper

import (
	"context"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

func (k msgServer) CreateListing(goCtx context.Context, msg *types.MsgCreateListing) (*types.MsgCreateListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if pause, found := k.GetPauseState(ctx); found && pause.Paused {
		return nil, types.ErrMarketplacePaused
	}
	if _, ok := new(big.Int).SetString(msg.PriceAmount, 10); !ok {
		return nil, types.ErrInvalidPrice
	}
	if msg.Denom == "" {
		return nil, types.ErrInvalidPrice
	}

	original, found := k.bookstoreKeeper.GetOriginal(ctx, msg.OriginalId)
	if !found {
		return nil, types.ErrListingNotFound
	}
	if original.Owner != msg.Creator {
		return nil, types.ErrNotOwner
	}

	id := k.AppendListing(ctx, types.Listing{
		OriginalId:  msg.OriginalId,
		Seller:      msg.Creator,
		PriceAmount: msg.PriceAmount,
		Denom:       msg.Denom,
		Active:      true,
	})

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"Listed",
		sdk.NewAttribute("listing_id", strconv.FormatUint(id, 10)),
		sdk.NewAttribute("original_id", strconv.FormatUint(msg.OriginalId, 10)),
		sdk.NewAttribute("seller", msg.Creator),
		sdk.NewAttribute("price", msg.PriceAmount),
		sdk.NewAttribute("denom", msg.Denom),
	))

	return &types.MsgCreateListingResponse{Id: id}, nil
}
