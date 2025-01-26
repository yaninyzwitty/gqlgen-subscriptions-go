package tests

import (
	"testing"

	"github.com/yaninyzwitty/gqlgen-subscriptions-go/helpers"
)

func TestUintToString(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{1234567890, "1234567890"},
		{18446744073709551615, "18446744073709551615"}, // Max uint64 value
	}

	for _, test := range tests {
		result := helpers.UintToString(test.input)
		if result != test.expected {
			t.Errorf("UintToString(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}
