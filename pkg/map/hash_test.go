package mapHelpers

import (
	"testing"
)

func TestHashMap(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "PrimitiveValues",
			input: map[string]interface{}{
				"a": 1,
				"b": "hello",
				"c": true,
			},
			expected: "dc315c80ce58e83cbca2c366f31bccffafd3fc557f2b8df0625fb8817b87036d",
		},
		{
			name: "NestedMaps",
			input: map[string]interface{}{
				"a": 1,
				"b": "hello",
				"c": map[string]interface{}{
					"nested": true,
					"value":  "nested value",
				},
			},
			expected: "d7fa61b6bb0552e3cc13fed01e9680b00401e1ad1e981dfb25f47eda1b50fe4e",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash := HashMap(test.input)
			if hash != test.expected {
				t.Errorf("Test case %s failed: expected %s, got %s", test.name, test.expected, hash)
			}
		})
	}
}
