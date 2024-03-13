package cases

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	tc "golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Case struct {
	slug string
}

func (r Case) String() string {
	return r.slug
}

func (r Case) Apply(value string) string {

	switch r.slug {
	case "none":
		return value
	case "camel":
		return strcase.ToLowerCamel(value)
	case "lower":
		return strings.ToLower(value)
	case "snake":
		return strcase.ToSnake(value)
	case "title":
		return tc.Title(language.English).String(value)
	case "upper":
		return strings.ToUpper(value)
	}
	return value
}

func FromString(slug string) (Case, error) {
	switch slug {
	case "":
		return Unknown, nil
	case "none":
		return None, nil
	case "camel":
		return CamelCase, nil
	case "lower":
		return LowerCase, nil
	case "snake":
		return SnakeCase, nil
	case "title":
		return TitleCase, nil
	case "upper":
		return UpperCase, nil
	}
	return Case{}, fmt.Errorf("unknown case: %s", slug)
}

var (
	Unknown   = Case{""}
	None      = Case{"none"}
	CamelCase = Case{"camel"}
	LowerCase = Case{"lower"}
	SnakeCase = Case{"snake"}
	TitleCase = Case{"title"}
	UpperCase = Case{"upper"}
)
