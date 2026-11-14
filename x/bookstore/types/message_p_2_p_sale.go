package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgP2PSale{}

func NewMsgP2PSale(creator string, originalId uint64, seller string, priceAmount string, denom string) *MsgP2PSale {
	return &MsgP2PSale{
		Creator:     creator,
		OriginalId:  originalId,
		Seller:      seller,
		PriceAmount: priceAmount,
		Denom:       denom,
	}
}

func (msg *MsgP2PSale) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
