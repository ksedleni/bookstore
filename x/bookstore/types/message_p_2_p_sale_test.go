package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ksedleni/bookstore/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgP2PSale_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgP2PSale
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgP2PSale{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgP2PSale{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
