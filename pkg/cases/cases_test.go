package cases

import (
	"testing"
)

func TestTCases(t *testing.T) {
	tests := []struct {
		caseType string
		input    string
		expected string
	}{
		{"none", "mytag", "mytag"},
		{"lower", "MYTAG", "mytag"},
		{"camel", "my_tag_example", "myTagExample"},
		{"snake", "MyTag", "my_tag"},
		{"title", "mytag", "Mytag"},
		{"none", "", ""},
		{"lower", "AbC", "abc"},
		{"camel", "CamelCase", "camelCase"},
		{"snake", "SNAKE_CASE", "snake_case"},
		{"title", "this is a title", "This Is A Title"},
		{"upper", "upper case", "UPPER CASE"},
	}

	for _, test := range tests {
		testCase, err := FromString(test.caseType)
		if err != nil {
			t.Errorf("For caseType %s with input %s, expected %s, but got error %s", test.caseType, test.input, test.expected, err.Error())
		}

		result := testCase.Apply(test.input)
		if result != test.expected {
			t.Errorf("For caseType %s with input %s, expected %s, but got %s", test.caseType, test.input, test.expected, result)
		}
	}
}
