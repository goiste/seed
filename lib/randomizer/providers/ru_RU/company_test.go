package ru

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInnCheckSum(t *testing.T) {
	testCases := []struct {
		name string
		inn  string
		exp  string
	}{
		{"empty", "", ""},
		{"short", "123456", ""},
		{"long", "1234567890", ""},
		{"invalid", "123test00", ""},
		{"7707049388", "770704938", "8"},
		{"7718099790", "771809979", "0"},
		{"7814148471", "781414847", "1"},
		{"7840346335", "784034633", "5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.exp, innCheckSum(tc.inn))
		})
	}
}
