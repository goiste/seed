package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFillForeignKeys(t *testing.T) {
	testCases := []struct {
		name        string
		table       Table
		foreignKeys map[string][]any
		needErr     bool
		expected    []map[string]any
	}{
		{"empty", Table{}, nil, false, nil},
		{
			"no error",
			Table{
				Name:        "test",
				ForeignKeys: []ForeignKey{{RefTable: "other", Column: "other_id"}},
				Values:      []map[string]any{{"other_id": 0}},
			},
			map[string][]any{"other": {3}},
			false,
			[]map[string]any{{"other_id": 3}},
		},
		{
			"has error",
			Table{
				Name:        "test",
				ForeignKeys: []ForeignKey{{"other", "other_id"}},
				Values:      []map[string]any{{"other_id": "wrong value"}},
			},
			map[string][]any{"other": {3}},
			true,
			nil,
		},
		{
			"two refs",
			Table{
				Name:        "test",
				ForeignKeys: []ForeignKey{{"first", "first_id"}, {"second", "second_id"}},
				Values:      []map[string]any{{"first_id": 0, "second_id": 1}, {"first_id": 1, "second_id": 0}},
			},
			map[string][]any{"first": {3, 5}, "second": {8, 13}},
			false,
			[]map[string]any{{"first_id": 3, "second_id": 13}, {"first_id": 5, "second_id": 8}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tableValues, err := (&Runner{}).fillForeignKeys(tc.table.Values, tc.table.ForeignKeys, tc.foreignKeys)
			if tc.needErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.EqualValues(t, tc.expected, tableValues)
		})
	}
}
