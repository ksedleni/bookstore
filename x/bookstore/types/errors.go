package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrInvalidSigner      = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidContentHash = sdkerrors.Register(ModuleName, 1101, "invalid content hash")
	ErrDuplicateContent   = sdkerrors.Register(ModuleName, 1102, "content hash already minted")
	ErrOriginalNotFound   = sdkerrors.Register(ModuleName, 1103, "original not found")
	ErrNotOwner           = sdkerrors.Register(ModuleName, 1104, "signer is not the owner")
	ErrInvalidSale        = sdkerrors.Register(ModuleName, 1105, "invalid p2p sale")
	ErrInvalidPrice       = sdkerrors.Register(ModuleName, 1106, "invalid price")
)
