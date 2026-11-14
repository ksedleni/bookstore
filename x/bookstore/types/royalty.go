package types

import "github.com/ksedleni/bookstore/protocol"

// Locked fee model — do not change without an ADR.
const (
	RoyaltyBPS     = protocol.RoyaltyBPS     // 10 = 0.1%
	ProtocolFeeBPS = protocol.ProtocolFeeBPS // 0
)
