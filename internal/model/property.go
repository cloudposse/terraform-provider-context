package model

import (
	"errors"
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

var (
	ErrPropertyRequired = errors.New("property is required")
	ErrValueTooShort    = errors.New("value is less than minimum length")
	ErrValueTooLong     = errors.New("value is greater than maximum length")
	ErrInvalidRegex     = errors.New("regex is invalid")
	ErrRegexMismatch    = errors.New("value does not match regex")
)

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
		return fmt.Errorf("%w: value for property %s", ErrPropertyRequired, propertyName)
	}
	return nil
}

func validateMinLength(minLength int, value string, propertyName string) error {
	if minLength == 0 {
		return nil
	}

	if len(value) < minLength {
		return fmt.Errorf("%w: value %s for property %s is less than %d", ErrValueTooShort, value, propertyName, minLength)
	}
	return nil
}

func validateMaxLength(maxLength int, value string, propertyName string) error {
	if maxLength == 0 {
		return nil
	}

	if len(value) > maxLength {
		return fmt.Errorf("%w: value %s for property %s is greater than %d", ErrValueTooLong, value, propertyName, maxLength)
	}
	return nil
}

func validateRegex(regex string, value string, propertyName string) error {
	if regex == "" || value == "" {
		return nil
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("%w: %s for property %s", ErrInvalidRegex, regex, propertyName)
	}

	if !r.MatchString(value) {
		return fmt.Errorf("%w: value %s for property %s does not match %s", ErrRegexMismatch, value, propertyName, regex)
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
