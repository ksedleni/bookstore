package keeper

import (
	"context"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ksedleni/bookstore/protocol"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

// P2PSale settles a paid peer-to-peer sale without any bookstore UI.
// Creator (signer) is the buyer. Protocol fee is always 0; royalty is fixed 0.1%.
func (k msgServer) P2PSale(goCtx context.Context, msg *types.MsgP2PSale) (*types.MsgP2PSaleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	buyer := msg.Creator
	if buyer == msg.Seller {
		return nil, types.ErrInvalidSale
	}
	original, found := k.GetOriginal(ctx, msg.OriginalId)
	if !found {
		return nil, types.ErrOriginalNotFound
	}
	if original.Owner != msg.Seller {
		return nil, types.ErrNotOwner
	}

	priceInt, ok := new(big.Int).SetString(msg.PriceAmount, 10)
	if !ok || priceInt.Sign() <= 0 {
		return nil, types.ErrInvalidPrice
	}
	split, err := protocol.SplitSale(priceInt)
	if err != nil {
		return nil, err
	}
	if split.ProtocolFee.Sign() != 0 {
		return nil, types.ErrInvalidSale
	}

	buyerAddr, err := sdk.AccAddressFromBech32(buyer)
	if err != nil {
		return nil, err
	}
	sellerAddr, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return nil, err
	}
	royaltyAddr, err := sdk.AccAddressFromBech32(original.RoyaltyPayee)
	if err != nil {
		return nil, err
	}

	if split.SellerProceeds.Sign() > 0 {
		sellerCoin := sdk.NewCoin(msg.Denom, sdkmath.NewIntFromBigInt(split.SellerProceeds))
		if err := k.bankKeeper.SendCoins(ctx, buyerAddr, sellerAddr, sdk.NewCoins(sellerCoin)); err != nil {
			return nil, err
		}
	}
	if split.Royalty.Sign() > 0 {
		royaltyCoin := sdk.NewCoin(msg.Denom, sdkmath.NewIntFromBigInt(split.Royalty))
		if err := k.bankKeeper.SendCoins(ctx, buyerAddr, royaltyAddr, sdk.NewCoins(royaltyCoin)); err != nil {
			return nil, err
		}
	}

	original.Owner = buyer
	k.SetOriginal(ctx, original)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"Sale",
		sdk.NewAttribute("original_id", strconv.FormatUint(msg.OriginalId, 10)),
		sdk.NewAttribute("buyer", buyer),
		sdk.NewAttribute("seller", msg.Seller),
		sdk.NewAttribute("price", msg.PriceAmount),
		sdk.NewAttribute("royalty", split.Royalty.String()),
		sdk.NewAttribute("seller_proceeds", split.SellerProceeds.String()),
		sdk.NewAttribute("protocol_fee", "0"),
		sdk.NewAttribute("path", "p2p"),
		sdk.NewAttribute("denom", msg.Denom),
	))

	return &types.MsgP2PSaleResponse{
		Royalty:        split.Royalty.String(),
		SellerProceeds: split.SellerProceeds.String(),
		ProtocolFee:    "0",
	}, nil
}
