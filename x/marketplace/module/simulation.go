package marketplace

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/ksedleni/bookstore/testutil/sample"
	marketplacesimulation "github.com/ksedleni/bookstore/x/marketplace/simulation"
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

// avoid unused import issue
var (
	_ = marketplacesimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateListing = "op_weight_msg_create_listing"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateListing int = 100

	opWeightMsgCancelListing = "op_weight_msg_cancel_listing"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCancelListing int = 100

	opWeightMsgBuyListing = "op_weight_msg_buy_listing"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBuyListing int = 100

	opWeightMsgSetPaused = "op_weight_msg_set_paused"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetPaused int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	marketplaceGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&marketplaceGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateListing int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateListing, &weightMsgCreateListing, nil,
		func(_ *rand.Rand) {
			weightMsgCreateListing = defaultWeightMsgCreateListing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateListing,
		marketplacesimulation.SimulateMsgCreateListing(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCancelListing int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelListing, &weightMsgCancelListing, nil,
		func(_ *rand.Rand) {
			weightMsgCancelListing = defaultWeightMsgCancelListing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelListing,
		marketplacesimulation.SimulateMsgCancelListing(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgBuyListing int
	simState.AppParams.GetOrGenerate(opWeightMsgBuyListing, &weightMsgBuyListing, nil,
		func(_ *rand.Rand) {
			weightMsgBuyListing = defaultWeightMsgBuyListing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBuyListing,
		marketplacesimulation.SimulateMsgBuyListing(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetPaused int
	simState.AppParams.GetOrGenerate(opWeightMsgSetPaused, &weightMsgSetPaused, nil,
		func(_ *rand.Rand) {
			weightMsgSetPaused = defaultWeightMsgSetPaused
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetPaused,
		marketplacesimulation.SimulateMsgSetPaused(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateListing,
			defaultWeightMsgCreateListing,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				marketplacesimulation.SimulateMsgCreateListing(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCancelListing,
			defaultWeightMsgCancelListing,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				marketplacesimulation.SimulateMsgCancelListing(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgBuyListing,
			defaultWeightMsgBuyListing,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				marketplacesimulation.SimulateMsgBuyListing(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSetPaused,
			defaultWeightMsgSetPaused,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				marketplacesimulation.SimulateMsgSetPaused(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
