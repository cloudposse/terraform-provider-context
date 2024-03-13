package stringHelpers

import (
	"testing"
)

func TestTruncateAndHash(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"myidunique1", 12, "myidunique1"},
		{"myidunique2", 8, "myid9641"},
		{"abcdefghijk", 8, "abc42571"},
		{"12345", 5, "12345"},
		{"abcdefghijklmnopqrstuvwxyz", 10, "abcde21474"},
		{"foo-bar-baz", 10, "foo-ba6094"},
	}

	for _, test := range tests {
		result := TruncateWithHash(test.input, test.maxLength)
		if result != test.expected {
			t.Errorf("For input %s with maxLength %d, expected %s, but got %s", test.input, test.maxLength, test.expected, result)
		}
	}
}
