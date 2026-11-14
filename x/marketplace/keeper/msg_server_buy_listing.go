package keeper

import (
	"context"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ksedleni/bookstore/protocol"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

func (k msgServer) BuyListing(goCtx context.Context, msg *types.MsgBuyListing) (*types.MsgBuyListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if pause, found := k.GetPauseState(ctx); found && pause.Paused {
		return nil, types.ErrMarketplacePaused
	}

	listing, found := k.GetListing(ctx, msg.ListingId)
	if !found || !listing.Active {
		return nil, types.ErrListingInactive
	}

	buyer := msg.Creator
	if buyer == listing.Seller {
		return nil, types.ErrInvalidSale
	}

	original, found := k.bookstoreKeeper.GetOriginal(ctx, listing.OriginalId)
	if !found {
		return nil, types.ErrListingNotFound
	}
	if original.Owner != listing.Seller {
		return nil, types.ErrNotOwner
	}

	priceInt, ok := new(big.Int).SetString(listing.PriceAmount, 10)
	if !ok || priceInt.Sign() <= 0 {
		return nil, types.ErrInvalidPrice
	}
	split, err := protocol.SplitSale(priceInt)
	if err != nil {
		return nil, err
	}

	buyerAddr, err := sdk.AccAddressFromBech32(buyer)
	if err != nil {
		return nil, err
	}
	sellerAddr, err := sdk.AccAddressFromBech32(listing.Seller)
	if err != nil {
		return nil, err
	}
	royaltyAddr, err := sdk.AccAddressFromBech32(original.RoyaltyPayee)
	if err != nil {
		return nil, err
	}

	if split.SellerProceeds.Sign() > 0 {
		sellerCoin := sdk.NewCoin(listing.Denom, sdkmath.NewIntFromBigInt(split.SellerProceeds))
		if err := k.bankKeeper.SendCoins(ctx, buyerAddr, sellerAddr, sdk.NewCoins(sellerCoin)); err != nil {
			return nil, err
		}
	}
	if split.Royalty.Sign() > 0 {
		royaltyCoin := sdk.NewCoin(listing.Denom, sdkmath.NewIntFromBigInt(split.Royalty))
		if err := k.bankKeeper.SendCoins(ctx, buyerAddr, royaltyAddr, sdk.NewCoins(royaltyCoin)); err != nil {
			return nil, err
		}
	}

	original.Owner = buyer
	k.bookstoreKeeper.SetOriginal(ctx, original)
	listing.Active = false
	k.SetListing(ctx, listing)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"Sale",
		sdk.NewAttribute("original_id", strconv.FormatUint(listing.OriginalId, 10)),
		sdk.NewAttribute("listing_id", strconv.FormatUint(msg.ListingId, 10)),
		sdk.NewAttribute("buyer", buyer),
		sdk.NewAttribute("seller", listing.Seller),
		sdk.NewAttribute("price", listing.PriceAmount),
		sdk.NewAttribute("royalty", split.Royalty.String()),
		sdk.NewAttribute("seller_proceeds", split.SellerProceeds.String()),
		sdk.NewAttribute("protocol_fee", "0"),
		sdk.NewAttribute("path", "listed"),
		sdk.NewAttribute("denom", listing.Denom),
	))

	return &types.MsgBuyListingResponse{
		Royalty:        split.Royalty.String(),
		SellerProceeds: split.SellerProceeds.String(),
		ProtocolFee:    "0",
	}, nil
}
