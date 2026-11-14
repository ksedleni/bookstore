package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrInvalidSigner     = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrMarketplacePaused = sdkerrors.Register(ModuleName, 1101, "marketplace module paused (P2P sale/transfer still allowed)")
	ErrNotOwner          = sdkerrors.Register(ModuleName, 1102, "signer is not the owner")
	ErrListingNotFound   = sdkerrors.Register(ModuleName, 1103, "listing not found")
	ErrListingInactive   = sdkerrors.Register(ModuleName, 1104, "listing not active")
	ErrInvalidPrice      = sdkerrors.Register(ModuleName, 1105, "invalid price")
	ErrInvalidSale       = sdkerrors.Register(ModuleName, 1106, "invalid listed sale")
	ErrUnauthorized      = sdkerrors.Register(ModuleName, 1107, "unauthorized")
)
