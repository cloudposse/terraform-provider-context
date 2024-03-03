package client

import (
	"fmt"
	"regexp"
	"strings"
)

type Property struct {
	IncludeInTags   bool
	MaxLength       int
	MinLength       int
	Name            string
	Required        bool
	ValidationRegex string
}

func (p *Property) Validate(value string) []error {
	errors := []error{}

	if err := validateRequired(p.Required, value, p.Name); err != nil {
		errors = append(errors, err)
	}

	if err := validateMinLength(p.MinLength, value, p.Name); err != nil {
		errors = append(errors, err)
	}

	if err := validateMaxLength(p.MaxLength, value, p.Name); err != nil {
		errors = append(errors, err)
	}

	if err := validateRegex(p.ValidationRegex, value, p.Name); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func validateRequired(required bool, value string, propertyName string) error {
	if required && strings.TrimSpace(value) == "" {
		return fmt.Errorf("value for property %s is required", propertyName)
	}

	return nil
}

func validateMinLength(minLength int, value string, propertyName string) error {
	if minLength == 0 {
		return nil
	}

	if len(value) < minLength {
		return fmt.Errorf("value %s for property %s is less than the minimum length of %d", value, propertyName, minLength)
	}

	return nil
}

func validateMaxLength(maxLength int, value string, propertyName string) error {
	if maxLength == 0 {
		return nil
	}

	if len(value) > maxLength {
		return fmt.Errorf("value %s for property %s is greater than the maximum length of %d", value, propertyName, maxLength)
	}

	return nil
}

func validateRegex(regex string, value string, propertyName string) error {
	if regex == "" {
		return nil
	}

	r, err := regexp.Compile(regex)

	if err != nil {
		return fmt.Errorf("regex %s for property %s is invalid", regex, propertyName)
	}

	if !r.MatchString(value) {
		return fmt.Errorf("value %s for property %s does not match the regex %s", value, propertyName, regex)
	}

	return nil
}

func NewProperty(name string, options ...func(*Property)) *Property {
	defaults := &Property{
		IncludeInTags:   true,
		MaxLength:       0,
		MinLength:       0,
		Name:            name,
		Required:        false,
		ValidationRegex: "",
	}

	for _, option := range options {
		option(defaults)
	}

	return defaults
}

func WithRequired() func(*Property) {
	return func(obj *Property) {
		obj.Required = true
	}
}

func WithExcludeFromTags() func(*Property) {
	return func(obj *Property) {
		obj.IncludeInTags = false
	}
}

func WithMinLength(minLength int) func(*Property) {
	return func(obj *Property) {
		obj.MinLength = minLength
	}
}

func WithMaxLength(maxLength int) func(*Property) {
	return func(obj *Property) {
		obj.MaxLength = maxLength
	}
}

func WithValidationRegex(regex string) func(*Property) {
	return func(obj *Property) {
		obj.ValidationRegex = regex
	}
}
