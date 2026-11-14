package types_test

import (
	"testing"

	"github.com/ksedleni/bookstore/x/bookstore/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				OriginalList: []types.Original{
					{
						Id: 1,
					},
					{
						Id: 2,
					},
				},
				OriginalCount: 2,
				ContentIndexList: []types.ContentIndex{
					{
						ContentHash: "0",
					},
					{
						ContentHash: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated original",
			genState: &types.GenesisState{
				OriginalList: []types.Original{
					{
						Id: 1,
					},
					{
						Id: 1,
					},
				},
				OriginalCount: 1,
			},
			valid: false,
		},
		{
			desc: "invalid original count",
			genState: &types.GenesisState{
				OriginalList: []types.Original{
					{
						Id: 1,
					},
				},
				OriginalCount: 0,
			},
			valid: false,
		},
		{
			desc: "zero original id rejected",
			genState: &types.GenesisState{
				OriginalList: []types.Original{
					{Id: 0},
				},
				OriginalCount: 1,
			},
			valid: false,
		},
		{
			desc: "duplicated contentIndex",
			genState: &types.GenesisState{
				ContentIndexList: []types.ContentIndex{
					{
						ContentHash: "0",
					},
					{
						ContentHash: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
