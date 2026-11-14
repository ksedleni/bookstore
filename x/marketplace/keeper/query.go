package keeper

import (
	"github.com/ksedleni/bookstore/x/marketplace/types"
)

var _ types.QueryServer = Keeper{}
