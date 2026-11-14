package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateListing{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCancelListing{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBuyListing{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetPaused{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
