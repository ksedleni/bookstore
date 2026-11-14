package marketplace

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/ksedleni/bookstore/api/bookstore/marketplace"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListingAll",
					Use:       "list-listing",
					Short:     "List all listing",
				},
				{
					RpcMethod:      "Listing",
					Use:            "show-listing [id]",
					Short:          "Shows a listing by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "PauseState",
					Use:       "show-pause-state",
					Short:     "show pause-state",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateListing",
					Use:            "create-listing [original-id] [price-amount] [denom]",
					Short:          "Send a create-listing tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "originalId"}, {ProtoField: "priceAmount"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod:      "CancelListing",
					Use:            "cancel-listing [listing-id]",
					Short:          "Send a cancel-listing tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "listingId"}},
				},
				{
					RpcMethod:      "BuyListing",
					Use:            "buy-listing [listing-id]",
					Short:          "Send a buy-listing tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "listingId"}},
				},
				{
					RpcMethod:      "SetPaused",
					Use:            "set-paused [paused]",
					Short:          "Send a set-paused tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "paused"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
