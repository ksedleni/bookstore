package types

const (
	ModuleName = "marketplace"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	ListingPrefix = []byte{0x01}
	PausedKey     = []byte{0x02}
)
