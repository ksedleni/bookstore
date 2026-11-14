package bookstore

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/ksedleni/bookstore/api/bookstore/bookstore"
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
					RpcMethod: "OriginalAll",
					Use:       "list-original",
					Short:     "List all original",
				},
				{
					RpcMethod:      "Original",
					Use:            "show-original [id]",
					Short:          "Shows a original by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ContentIndexAll",
					Use:       "list-content-index",
					Short:     "List all content-index",
				},
				{
					RpcMethod:      "ContentIndex",
					Use:            "show-content-index [id]",
					Short:          "Shows a content-index",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "contentHash"}},
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
					RpcMethod:      "MintOriginal",
					Use:            "mint-original [content-hash] [metadata-hash] [metadata-uri]",
					Short:          "Send a mint-original tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "contentHash"}, {ProtoField: "metadataHash"}, {ProtoField: "metadataUri"}},
				},
				{
					RpcMethod:      "TransferOriginal",
					Use:            "transfer-original [original-id] [to]",
					Short:          "Send a transfer-original tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "originalId"}, {ProtoField: "to"}},
				},
				{
					RpcMethod:      "P2PSale",
					Use:            "p-2-p-sale [original-id] [seller] [price-amount] [denom]",
					Short:          "Send a p2p-sale tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "originalId"}, {ProtoField: "seller"}, {ProtoField: "priceAmount"}, {ProtoField: "denom"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
