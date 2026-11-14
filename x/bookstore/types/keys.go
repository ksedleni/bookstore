package types

const (
	// ModuleName defines the module name
	ModuleName = "bookstore"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_bookstore"
)

var (
	ParamsKey = []byte("p_bookstore")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	OriginalKey      = "Original/value/"
	OriginalCountKey = "Original/count/"
)
