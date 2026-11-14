package keeper

import (
	"github.com/ksedleni/bookstore/x/bookstore/types"
)

var _ types.QueryServer = Keeper{}
