package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/ksedleni/bookstore/x/bookstore/keeper"
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

func SimulateMsgP2PSale(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgP2PSale{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the P2PSale simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "P2PSale simulation not implemented"), nil, nil
	}
}
