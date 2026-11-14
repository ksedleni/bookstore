package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMintOriginal{}

func NewMsgMintOriginal(creator string, contentHash string, metadataHash string, metadataUri string) *MsgMintOriginal {
	return &MsgMintOriginal{
		Creator:      creator,
		ContentHash:  contentHash,
		MetadataHash: metadataHash,
		MetadataUri:  metadataUri,
	}
}

func (msg *MsgMintOriginal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
