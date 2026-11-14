package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ListingList: []Listing{},
		PauseState:  nil,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in listing
	listingIdMap := make(map[uint64]bool)
	listingCount := gs.GetListingCount()
	for _, elem := range gs.ListingList {
		if _, ok := listingIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for listing")
		}
		if elem.Id >= listingCount {
			return fmt.Errorf("listing id should be lower or equal than the last id")
		}
		listingIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
