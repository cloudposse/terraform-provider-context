package model

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/cloudposse/terraform-provider-context/pkg/cases"
	"github.com/cloudposse/terraform-provider-context/pkg/slice"
	"github.com/cloudposse/terraform-provider-context/pkg/stringHelpers"
)

var ErrLabelTooLong = errors.New("label exceeds maximum length")

type ProviderConfig struct {
	delimiter         string
	enabled           bool
	properties        []Property
	propertyOrder     []string
	replaceCharsRegex string
	tagsKeyCase       cases.Case
	tagsValueCase     cases.Case
	values            map[string]string
}

type DelmitedLabelOptions struct {
	Delimiter  *string
	Properties []string
	Values     map[string]string
}

// GetDelimiter returns the delimiter from the context.
func (c *ProviderConfig) GetDelimiter() string {
	return c.delimiter
}

// IsEnabled returns a boolean indicating whether the context is enabled.
func (c *ProviderConfig) IsEnabled() bool {
	return c.enabled
}

// GetProperties returns the properties from the context.
func (c *ProviderConfig) GetProperties() []Property {
	return c.properties
}

func (c *ProviderConfig) ValidateProperties(values map[string]string) []error {
	errors := []error{}
	for _, p := range c.properties {
		errors = append(errors, p.Validate(values[p.Name])...)
	}
	return errors
}

// GetPropertyNames returns the names of the properties from the context.
func (c *ProviderConfig) GetPropertyNames([]Property) []string {
	names := []string{}
	for _, p := range c.properties {
		names = append(names, p.Name)
	}
	return names
}

// GetMergedDelimiter merges the delimiter from the context with the delimiter passed in to the function. Used when
// creating a label.
func (c *ProviderConfig) GetMergedDelimiter(delimiter *string) string {
	mergedDelimiter := c.delimiter
	if delimiter != nil {
		mergedDelimiter = *delimiter
	}

	return mergedDelimiter
}

// GetMergedReplaceCharsRegex merges the replaceCharsRegex from the context with the replaceCharsRegex passed in to the
// function. Used when creating a label.
func (c *ProviderConfig) GetMergedReplaceCharsRegex(regex *string) string {
	mergedRegex := c.replaceCharsRegex
	if regex != nil {
		mergedRegex = *regex
	}

	return mergedRegex
}

// GetMergedPropertyNames returns either the names of the properties from the context or the names of the properties
// passed in to the function to derive the properties to use for creating a label.
func (c *ProviderConfig) GetMergedPropertyNames(propertyNames []string) []string {
	if len(propertyNames) > 0 {
		return propertyNames
	}
	return c.GetPropertyNames(c.properties)
}

// GetPropertyOrder returns the propertyOrder from the context.
func (c *ProviderConfig) GetPropertyOrder() []string {
	return c.propertyOrder
}

// getMergedPropertyOrder returns either the names in order of the propertyOrder from the context or the propertyOrder
// passed in to the function to derive the order of properties to use for creating a label.
func (c *ProviderConfig) GetMergedPropertyOrder(propertyOrder []string) []string {
	if len(propertyOrder) > 0 {
		return propertyOrder
	}
	return c.propertyOrder
}

// GetReplaceCharsRegex returns the replaceCharsRegex from the context.
func (c *ProviderConfig) GetReplaceCharsRegex() string {
	return c.replaceCharsRegex
}

// GetTagsKeyCase returns the tagsKeyCase from the context.
func (c *ProviderConfig) GetTagsKeyCase() string {
	return c.tagsKeyCase.String()
}

// GetTagsValueCase returns the tagsValueCase from the context.
func (c *ProviderConfig) GetTagsValueCase() string {
	return c.tagsValueCase.String()
}

// GetTagsKeyCase returns the tagsKeyCase from the context or the keyCase passed in to the function.
func (c *ProviderConfig) GetMergedTagsKeyCase(keyCase *cases.Case) cases.Case {
	if keyCase != nil {
		return *keyCase
	}
	return c.tagsKeyCase
}

// GetTagsValueCase returns the tagsValueCase from the context or the valueCase passed in to the function.
func (c *ProviderConfig) GetMergedTagsValueCase(valueCase *cases.Case) cases.Case {
	if valueCase != nil {
		return *valueCase
	}
	return c.tagsValueCase
}

// GetValues returns the values from the context.
func (c *ProviderConfig) GetValues() map[string]string {
	return c.values
}

// getMergedValues merges the values from the context with the values passed in to the function to derive the values to
// use when creating a label.
func (c *ProviderConfig) GetMergedValues(values map[string]string) map[string]string {
	mergedValues := make(map[string]string, len(c.values))
	for key, value := range c.values {
		mergedValues[key] = value
	}

	for key, value := range values {
		mergedValues[key] = value
	}

	return mergedValues
}

// getOrderedValues returns the values in the order of the propertyOrder for use in creating a delimited label.
func (c *ProviderConfig) getOrderedValues(propertyOrder []string, values map[string]string) []string {
	orderedValues := []string{}
	for _, prop := range propertyOrder {
		if values[prop] != "" {
			orderedValues = append(orderedValues, values[prop])
		}
	}
	return orderedValues
}

func getRedactedLabel(label string, regex string) (string, error) {
	if regex == "" {
		return label, nil
	}

	compiledRegex, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	replaced := compiledRegex.ReplaceAllString(label, "")
	return replaced, nil
}

//nolint:revive
func (c *ProviderConfig) GetDelimitedLabel(delimiter *string, properties []string, propertyOrder []string, values map[string]string, replaceCharsRegex *string, maxLength int, truncateIfExceedsMaxLength bool) (string, []error) {
	mergedValues := c.GetMergedValues(values)
	regex := c.GetMergedReplaceCharsRegex(replaceCharsRegex)
	validationErrors := c.ValidateProperties(mergedValues)
	if len(validationErrors) > 0 {
		return "", validationErrors
	}

	mergedDelimiter := c.GetMergedDelimiter(delimiter)
	mergedProperties := c.GetMergedPropertyNames(properties)
	mergedPropertyOrder := c.GetMergedPropertyOrder(propertyOrder)
	filteredPropertyOrder := []string{}
	for _, prop := range mergedPropertyOrder {
		if slice.Contains(mergedProperties, prop) {
			filteredPropertyOrder = append(filteredPropertyOrder, prop)
		}
	}
	orderedValues := c.getOrderedValues(filteredPropertyOrder, mergedValues)

	label := strings.Join(orderedValues, mergedDelimiter)

	redactedLabel, err := getRedactedLabel(label, regex)
	if err != nil {
		return "", []error{err}
	}

	if maxLength > 0 && len(redactedLabel) > maxLength {
		if !truncateIfExceedsMaxLength {
			return "", []error{fmt.Errorf("%w: %s (max: %d)", ErrLabelTooLong, redactedLabel, maxLength)}
		}
		return stringHelpers.TruncateWithHash(redactedLabel, maxLength), nil
	}

	return redactedLabel, nil
}

// GetTemplatedLabel returns a label from the template string and based on the properties and values in the context and
// overridden by the delimiter, properties and values passed into the function.
func (c *ProviderConfig) GetTemplatedLabel(templateString string, values map[string]string, replaceCharsRegex *string, maxLength int, truncateIfExceedsMaxLength bool) (string, []error) {
	mergedValues := c.GetMergedValues(values)
	regex := c.GetMergedReplaceCharsRegex(replaceCharsRegex)
	validationErrors := c.ValidateProperties(mergedValues)
	if len(validationErrors) > 0 {
		return "", validationErrors
	}

	tmpl, err := template.New("label").Parse(templateString)
	if err != nil {
		return "", []error{err}
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, mergedValues)
	if err != nil {
		return "", []error{err}
	}

	label := result.String()
	redactedLabel, err := getRedactedLabel(label, regex)
	if err != nil {
		return "", []error{err}
	}

	if maxLength > 0 && len(redactedLabel) > maxLength {
		if !truncateIfExceedsMaxLength {
			return "", []error{fmt.Errorf("%w: %s (max: %d)", ErrLabelTooLong, redactedLabel, maxLength)}
		}
		return stringHelpers.TruncateWithHash(redactedLabel, maxLength), nil
	}

	return redactedLabel, nil
}

func getCasedTag(key string, value string, keyCase cases.Case, valueCase cases.Case) (string, string) {
	keyValue := keyCase.Apply(key)
	valueValue := valueCase.Apply(value)

	return keyValue, valueValue
}

func (c *ProviderConfig) GetTags(values map[string]string, tagsKeyCase *cases.Case, tagsValueCase *cases.Case) (map[string]string, []error) {
	tags := map[string]string{}
	mergedValues := c.GetMergedValues(values)
	validationErrors := c.ValidateProperties(mergedValues)
	mergedTagsKeyCase := c.GetMergedTagsKeyCase(tagsKeyCase)
	mergedTagsValueCase := c.GetMergedTagsValueCase(tagsValueCase)

	if len(validationErrors) > 0 {
		return tags, validationErrors
	}

	for _, p := range c.properties {
		if p.IncludeInTags {
			// Use property-specific casing if available, otherwise use provider-level casing
			keyCase := mergedTagsKeyCase
			if p.TagsKeyCase != nil {
				keyCase = *p.TagsKeyCase
			}
			valueCase := mergedTagsValueCase
			if p.TagsValueCase != nil {
				valueCase = *p.TagsValueCase
			}
			key, value := getCasedTag(p.Name, mergedValues[p.Name], keyCase, valueCase)
			if value != "" {
				tags[key] = value
			}
		}
	}
	return tags, nil
}

func (c *ProviderConfig) GetTagsAsList(values map[string]string, tagsKeyCase *cases.Case, tagsValueCase *cases.Case) ([]map[string]string, []error) {
	tags, err := c.GetTags(values, tagsKeyCase, tagsValueCase)
	if err != nil {
		return nil, err
	}

	tagsList := []map[string]string{}

	// Sort the tags by key
	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tagsList = append(tagsList, map[string]string{"Key": k, "Value": tags[k]})
	}
	return tagsList, nil
}

// NewProviderConfig is the factory for creating a new provider config.
func NewProviderConfig(properties []Property, propertyOrder []string, values map[string]string, options ...func(*ProviderConfig)) (*ProviderConfig, error) {
	cc := &ProviderConfig{
		delimiter:         "-",
		enabled:           true,
		properties:        properties,
		replaceCharsRegex: "",
		tagsKeyCase:       cases.TitleCase,
		tagsValueCase:     cases.None,
		values:            values,
	}

	cc.propertyOrder = cc.GetMergedPropertyOrder(cc.GetPropertyNames(properties))
	cc.propertyOrder = cc.GetMergedPropertyOrder(propertyOrder)

	for _, option := range options {
		option(cc)
	}

	return cc, nil
}

// WithProperties is a functional option for setting the properties in the context when creating a new provider config.
func WithEnabled(enabled bool) func(*ProviderConfig) {
	return func(obj *ProviderConfig) {
		obj.enabled = enabled
	}
}

// WithProperties is a functional option for setting the properties in the context when creating a new provider config.
func WithDelimiter(delimiter string) func(*ProviderConfig) {
	return func(obj *ProviderConfig) {
		obj.delimiter = delimiter
	}
}

// WithReplaceCharsRegex is a functional option for setting the properties in the context when creating a new provider config.
func WithReplaceCharsRegex(regex string) func(*ProviderConfig) {
	return func(obj *ProviderConfig) {
		obj.replaceCharsRegex = regex
	}
}

func WithTagsKeyCase(keyCase cases.Case) func(*ProviderConfig) {
	return func(obj *ProviderConfig) {
		obj.tagsKeyCase = keyCase
	}
}

func WithTagsValueCase(valueCase cases.Case) func(*ProviderConfig) {
	return func(obj *ProviderConfig) {
		obj.tagsValueCase = valueCase
	}
}
