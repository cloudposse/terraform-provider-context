package cases

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
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
		return strcase.ToCamel(value)
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
	return Case{}, fmt.Errorf("%w: %s", ErrUnknownCase, slug)
}

var (
	Unknown        = Case{""}
	None           = Case{"none"}
	CamelCase      = Case{"camel"}
	LowerCase      = Case{"lower"}
	SnakeCase      = Case{"snake"}
	TitleCase      = Case{"title"}
	UpperCase      = Case{"upper"}
	ErrUnknownCase = errors.New("unknown case")
)
