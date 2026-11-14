package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ContentIndexKeyPrefix is the prefix to retrieve all ContentIndex
	ContentIndexKeyPrefix = "ContentIndex/value/"
)

// ContentIndexKey returns the store key to retrieve a ContentIndex from the index fields
func ContentIndexKey(
	contentHash string,
) []byte {
	var key []byte

	contentHashBytes := []byte(contentHash)
	key = append(key, contentHashBytes...)
	key = append(key, []byte("/")...)

	return key
}
