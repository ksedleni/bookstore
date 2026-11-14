package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		OriginalList:     []Original{},
		ContentIndexList: []ContentIndex{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in original
	originalIdMap := make(map[uint64]bool)
	originalCount := gs.GetOriginalCount()
	for _, elem := range gs.OriginalList {
		if _, ok := originalIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for original")
		}
		if elem.Id >= originalCount {
			return fmt.Errorf("original id should be lower or equal than the last id")
		}
		originalIdMap[elem.Id] = true
	}
	// Check for duplicated index in contentIndex
	contentIndexIndexMap := make(map[string]struct{})

	for _, elem := range gs.ContentIndexList {
		index := string(ContentIndexKey(elem.ContentHash))
		if _, ok := contentIndexIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for contentIndex")
		}
		contentIndexIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
