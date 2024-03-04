package client

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/cloudposse/terraform-provider-context/pkg/slice"
	"github.com/cloudposse/terraform-provider-context/pkg/stringHelpers"
)

type Client struct {
	delimiter     string
	enabled       bool
	properties    []Property
	propertyOrder []string
	values        map[string]string
}

type DelmitedLabelOptions struct {
	Delimiter  *string
	Properties []string
	Values     map[string]string
}

// GetDelimiter returns the delimiter from the context.
func (c *Client) GetDelimiter() string {
	return c.delimiter
}

// GetMergedDelmitier merges the delimiter from the context with the delimiter passed in to the function. Used when
// creating a label.
func (c *Client) GetMergedDelmitier(delimiter *string) string {
	mergedDelimiter := c.delimiter
	if delimiter != nil {
		mergedDelimiter = *delimiter
	}

	return mergedDelimiter
}

// IsEnabled returns a boolean indicating whether the context is enabled.
func (c *Client) IsEnabled() bool {
	return c.enabled
}

// GetProperties returns the properties from the context.
func (c *Client) GetProperties() []Property {
	return c.properties
}

func (c *Client) ValidateProperties(values map[string]string) []error {
	errors := []error{}
	for _, p := range c.properties {
		errors = append(errors, p.Validate(values[p.Name])...)
	}
	return errors
}

// GetPropertyNames returns the names of the properties from the context.
func (c *Client) GetPropertyNames([]Property) []string {
	names := []string{}
	for _, p := range c.properties {
		names = append(names, p.Name)
	}
	return names
}

// GetMergedPropertyNames returns either the names of the properties from the context or the names of the properties
// passed in to the function to derive the properties to use for creating a label.
func (c *Client) GetMergedPropertyNames(propertyNames []string) []string {
	if len(propertyNames) > 0 {
		return propertyNames
	}
	return c.GetPropertyNames(c.properties)
}

// GetPropertyOrder returns the propertyOrder from the context.
func (c *Client) GetPropertyOrder() []string {
	return c.propertyOrder
}

// getMergedPropertyOrder returns either the names in order of the propertyOrder from the context or the propertyOrder
// passed in to the function to derive the order of properties to use for creating a label.
func (c *Client) GetMergedPropertyOrder(propertyOrder []string) []string {
	if len(propertyOrder) > 0 {
		return propertyOrder
	}
	return c.propertyOrder
}

// GetValues returns the values from the context.
func (c *Client) GetValues() map[string]string {
	return c.values
}

// getMergedValues merges the values from the context with the values passed in to the function to derive the values to
// use when creating a label.
func (c *Client) GetMergedValues(values map[string]string) map[string]string {
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
func (c *Client) getOrderedValues(propertyOrder []string, values map[string]string) []string {
	orderedValues := []string{}
	for _, prop := range propertyOrder {
		orderedValues = append(orderedValues, values[prop])
	}
	return orderedValues
}

func getTruncatedLabel(label string, maxLength int, truncateIfExceedsMaxLength bool) (string, []error) {
	if maxLength > 0 && len(label) > maxLength {
		if truncateIfExceedsMaxLength {
			label = stringHelpers.TruncateWithHash(label, maxLength)
		} else {
			return "", []error{fmt.Errorf("label %s exceeds maximum length of %d", label, maxLength)}
		}
	}
	return label, nil
}

// GetDelimitedLabel returns a delimited label based on the properties and values in the context and overridden by the
// delimiter, properties and values passed into the function.
func (c *Client) GetDelimitedLabel(delimiter *string, properties []string, propertyOrder []string, values map[string]string, maxLength int, truncateIfExceedsMaxLength bool) (string, []error) {
	mergedValues := c.GetMergedValues(values)

	validationErrors := c.ValidateProperties(mergedValues)
	if len(validationErrors) > 0 {
		return "", validationErrors
	}

	mergedDelimiter := c.GetMergedDelmitier(delimiter)
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

	return getTruncatedLabel(label, maxLength, truncateIfExceedsMaxLength)
}

// GetTemplatedLabel returns a label from the template string and based on the properties and values in the context and
// overridden by the delimiter, properties and values passed into the function.
func (c *Client) GetTemplatedLabel(templateString string, values map[string]string, maxLength int, truncateIfExceedsMaxLength bool) (string, []error) {
	mergedValues := c.GetMergedValues(values)
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
	return getTruncatedLabel(label, maxLength, truncateIfExceedsMaxLength)
}

func (c *Client) GetTags(values map[string]string) (map[string]string, []error) {
	tags := map[string]string{}
	mergedValues := c.GetMergedValues(values)
	validationErrors := c.ValidateProperties(mergedValues)
	if len(validationErrors) > 0 {
		return tags, validationErrors
	}

	for _, p := range c.properties {
		if p.IncludeInTags {
			tags[p.Name] = mergedValues[p.Name]
		}
	}
	return tags, nil
}

func (c *Client) GetTagsAsList(values map[string]string) ([]map[string]string, []error) {
	tags, err := c.GetTags(values)
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

// NewClient is the factory for creating a new context client.
func NewClient(properties []Property, propertyOrder []string, values map[string]string, options ...func(*Client)) (*Client, error) {
	cc := &Client{
		delimiter:  "-",
		enabled:    true,
		properties: properties,
		values:     values,
	}

	cc.propertyOrder = cc.GetMergedPropertyOrder(cc.GetPropertyNames(properties))
	cc.propertyOrder = cc.GetMergedPropertyOrder(propertyOrder)

	for _, option := range options {
		option(cc)
	}

	return cc, nil
}

// WithProperties is a functional option for setting the properties in the context when creating a new context client.
func WithEnabled(enabled bool) func(*Client) {
	return func(obj *Client) {
		obj.enabled = enabled
	}
}

// WithProperties is a functional option for setting the properties in the context when creating a new context client.
func WithDelimiter(delimiter string) func(*Client) {
	return func(obj *Client) {
		obj.delimiter = delimiter
	}
}
