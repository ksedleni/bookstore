package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgTransferOriginal{}

func NewMsgTransferOriginal(creator string, originalId uint64, to string) *MsgTransferOriginal {
	return &MsgTransferOriginal{
		Creator:    creator,
		OriginalId: originalId,
		To:         to,
	}
}

func (msg *MsgTransferOriginal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
