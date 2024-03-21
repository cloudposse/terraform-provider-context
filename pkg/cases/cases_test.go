package cases

import (
	"testing"
)

func TestCases(t *testing.T) {
	tests := []struct {
		caseType string
		input    string
		expected string
	}{
		{"none", "mytag", "mytag"},
		{"lower", "MYTAG", "mytag"},
		{"camel", "my_tag_example", "myTagExample"},
		{"camel", "support email", "supportEmail"},
		{"snake", "MyTag", "my_tag"},
		{"none", "", ""},
		{"lower", "AbC", "abc"},
		{"camel", "CamelCase", "camelCase"},
		{"snake", "SNAKE_CASE", "snake_case"},
		{"title", "mytag", "Mytag"},
		{"title", "this is a title", "ThisIsATitle"},
		{"title", "support_email", "SupportEmail"},
		{"title", "support-email", "SupportEmail"},
		{"title", "support-email-address", "SupportEmailAddress"},
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
