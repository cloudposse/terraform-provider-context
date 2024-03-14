package stringHelpers

import (
	"testing"
)

func TestHashString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"world", "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"This is a test string.", "3eec256a587cccf72f71d2342b6dfab0bbca01697c7e7014540bdd62b72120da"},
	}

	for _, tc := range testCases {
		actual := HashString(tc.input)
		if actual != tc.expected {
			t.Errorf("Expected hash for input '%s' to be '%s', but got '%s'", tc.input, tc.expected, actual)
		}
	}
}
