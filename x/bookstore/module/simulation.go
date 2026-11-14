package bookstore

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/ksedleni/bookstore/testutil/sample"
	bookstoresimulation "github.com/ksedleni/bookstore/x/bookstore/simulation"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

// avoid unused import issue
var (
	_ = bookstoresimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgMintOriginal = "op_weight_msg_mint_original"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMintOriginal int = 100

	opWeightMsgTransferOriginal = "op_weight_msg_transfer_original"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTransferOriginal int = 100

	opWeightMsgP2PSale = "op_weight_msg_p_2_p_sale"
	// TODO: Determine the simulation weight value
	defaultWeightMsgP2PSale int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	bookstoreGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&bookstoreGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgMintOriginal int
	simState.AppParams.GetOrGenerate(opWeightMsgMintOriginal, &weightMsgMintOriginal, nil,
		func(_ *rand.Rand) {
			weightMsgMintOriginal = defaultWeightMsgMintOriginal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintOriginal,
		bookstoresimulation.SimulateMsgMintOriginal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgTransferOriginal int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferOriginal, &weightMsgTransferOriginal, nil,
		func(_ *rand.Rand) {
			weightMsgTransferOriginal = defaultWeightMsgTransferOriginal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferOriginal,
		bookstoresimulation.SimulateMsgTransferOriginal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgP2PSale int
	simState.AppParams.GetOrGenerate(opWeightMsgP2PSale, &weightMsgP2PSale, nil,
		func(_ *rand.Rand) {
			weightMsgP2PSale = defaultWeightMsgP2PSale
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgP2PSale,
		bookstoresimulation.SimulateMsgP2PSale(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgMintOriginal,
			defaultWeightMsgMintOriginal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				bookstoresimulation.SimulateMsgMintOriginal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgTransferOriginal,
			defaultWeightMsgTransferOriginal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				bookstoresimulation.SimulateMsgTransferOriginal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgP2PSale,
			defaultWeightMsgP2PSale,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				bookstoresimulation.SimulateMsgP2PSale(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
