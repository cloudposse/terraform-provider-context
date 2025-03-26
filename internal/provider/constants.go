package provider

const (
	// CaseNone represents no case transformation.
	CaseNone = "none"
	// CaseCamel represents camelCase transformation.
	CaseCamel = "camel"
	// CaseLower represents lowercase transformation.
	CaseLower = "lower"
	// CaseSnake represents snake_case transformation.
	CaseSnake = "snake"
	// CaseTitle represents TitleCase transformation.
	CaseTitle = "title"
	// CaseUpper represents UPPERCASE transformation.
	CaseUpper = "upper"
)

// ValidCases contains all valid case values.
var ValidCases = []string{CaseNone, CaseCamel, CaseLower, CaseSnake, CaseTitle, CaseUpper}
