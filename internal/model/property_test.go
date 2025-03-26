package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPropertyWithDefaults(t *testing.T) {
	p := NewProperty("test")

	assert.Equal(t, "test", p.Name)
	assert.Equal(t, false, p.Required)
	assert.Equal(t, 0, p.MinLength)
	assert.Equal(t, 0, p.MaxLength)
	assert.Equal(t, "", p.ValidationRegex)
}

func TestNewPropertyWithRequired(t *testing.T) {
	p := NewProperty("test", WithRequired())

	assert.Equal(t, "test", p.Name)
	assert.Equal(t, true, p.Required)
	assert.Equal(t, 0, p.MinLength)
	assert.Equal(t, 0, p.MaxLength)
	assert.Equal(t, "", p.ValidationRegex)
}

func TestNewPropertyWithMinLength(t *testing.T) {
	p := NewProperty("test", WithMinLength(5))

	assert.Equal(t, "test", p.Name)
	assert.Equal(t, false, p.Required)
	assert.Equal(t, 5, p.MinLength)
	assert.Equal(t, 0, p.MaxLength)
	assert.Equal(t, "", p.ValidationRegex)
}

func TestNewPropertyWithMaxLength(t *testing.T) {
	p := NewProperty("test", WithMaxLength(5))

	assert.Equal(t, "test", p.Name)
	assert.Equal(t, false, p.Required)
	assert.Equal(t, 0, p.MinLength)
	assert.Equal(t, 5, p.MaxLength)
	assert.Equal(t, "", p.ValidationRegex)
}

func TestNewPropertyWithValidationRegex(t *testing.T) {
	p := NewProperty("test", WithValidationRegex("^[a-z]+$"))

	assert.Equal(t, "test", p.Name)
	assert.Equal(t, false, p.Required)
	assert.Equal(t, 0, p.MinLength)
	assert.Equal(t, 0, p.MaxLength)
	assert.Equal(t, "^[a-z]+$", p.ValidationRegex)
}

func TestPropertyValidateWithRequiredInvalid(t *testing.T) {
	p := NewProperty("test", WithRequired())

	err := p.Validate("")

	assert.Equal(t, 1, len(err))
	assert.Equal(t, "property is required: value for property test", err[0].Error())
}

func TestPropertyValidateWithRequiredValid(t *testing.T) {
	p := NewProperty("test", WithRequired())

	err := p.Validate("test")

	assert.Equal(t, 0, len(err))
}

func TestPropertyValidateWithMinLengthInvalid(t *testing.T) {
	p := NewProperty("test", WithMinLength(5))

	err := p.Validate("test")

	assert.Equal(t, 1, len(err))
	assert.Equal(t, "value is less than minimum length: value test for property test is less than 5", err[0].Error())
}

func TestPropertyValidateWithMinLengthValid(t *testing.T) {
	p := NewProperty("test", WithMinLength(5))

	err := p.Validate("testing")

	assert.Equal(t, 0, len(err))
}

func TestPropertyValidateWithMaxLengthInvalid(t *testing.T) {
	p := NewProperty("test", WithMaxLength(5))

	err := p.Validate("testing")

	assert.Equal(t, 1, len(err))
	assert.Equal(t, "value is greater than maximum length: value testing for property test is greater than 5", err[0].Error())
}

func TestPropertyValidateWithMaxLengthValid(t *testing.T) {
	p := NewProperty("test", WithMaxLength(5))

	err := p.Validate("test")

	assert.Equal(t, 0, len(err))
}

func TestPropertyValidateWithValidationRegexInvalid(t *testing.T) {
	p := NewProperty("test", WithValidationRegex("^[a-z]+$"))

	err := p.Validate("123")

	assert.Equal(t, 1, len(err))
	assert.Equal(t, "value does not match regex: value 123 for property test does not match ^[a-z]+$", err[0].Error())
}

func TestPropertyValidateWithValidationRegexValid(t *testing.T) {
	p := NewProperty("test", WithValidationRegex("^[a-z]+$"))

	err := p.Validate("abc")

	assert.Equal(t, 0, len(err))
}

func TestPropertyExcludeFromTags(t *testing.T) {
	p := NewProperty("test", WithExcludeFromTags())

	actual := p.IncludeInTags
	assert.Equal(t, false, actual)
}
