package helpers

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{"EXISTING_KEY", "value", "default", "value"},
		{"NON_EXISTING_KEY", "", "default", "default"},
	}

	for _, test := range tests {
		if test.value != "" {
			os.Setenv(test.key, test.value)
		} else {
			os.Unsetenv(test.key)
		}

		result := GetEnvOrDefault(test.key, test.defaultValue)
		if result != test.expected {
			t.Errorf("GetEnvOrDefault(%s, %s) = %s; want %s", test.key, test.defaultValue, result, test.expected)
		}
	}
}
