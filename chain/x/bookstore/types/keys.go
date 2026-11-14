package types

const (
	// ModuleName is the Cosmos SDK module name for ownership + P2P sale/transfer.
	ModuleName = "bookstore"

	StoreKey = ModuleName

	// RouterKey is the message route.
	RouterKey = ModuleName
)

// Prefixes for KV store keys.
var (
	OriginalPrefix     = []byte{0x01}
	ContentHashPrefix  = []byte{0x02}
	OriginalCountKey   = []byte{0x03}
)
