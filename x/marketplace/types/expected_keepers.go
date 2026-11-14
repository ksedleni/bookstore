package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bookstoretypes "github.com/ksedleni/bookstore/x/bookstore/types"
)

// BookstoreKeeper is the ownership module surface used by marketplace listings.
type BookstoreKeeper interface {
	GetOriginal(ctx context.Context, id uint64) (val bookstoretypes.Original, found bool)
	SetOriginal(ctx context.Context, original bookstoretypes.Original)
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
