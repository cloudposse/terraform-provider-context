package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDefaultProviderConfig(t *testing.T, enabled bool) *ProviderConfig {
	properties := []Property{*NewProperty("foo"), *NewProperty("bar"), *NewProperty("baz")}
	values := map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"}

	options := []func(*ProviderConfig){}
	if !enabled {
		options = append(options, WithEnabled(false))
	}
	c, err := NewProviderConfig(properties, []string{}, values, options...)
	assert.NoError(t, err)

	return c
}

func TestProviderConfigIsEnabled(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	assert.Equal(t, true, c.IsEnabled())
}

func TestProviderConfigInitIsNotEnabled(t *testing.T) {
	c := getDefaultProviderConfig(t, false)
	assert.Equal(t, false, c.IsEnabled())
}

func TestProviderConfigGetDelimitedLabelWithNoProperties(t *testing.T) {
	properties := []Property{}
	values := map[string]string{}
	c, err := NewProviderConfig(properties, []string{}, values)
	assert.NoError(t, err)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "", actual)
}

func TestProviderConfigGetDelimitedLabelWithNoLocalOverrides(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo-bar-baz", actual)
}

func TestProviderConfigGetDelimitedLabelWithTruncation(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, nil, nil, 10, true)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo-ba6094", actual)
}

func TestProviderConfigGetDelimitedLabelWithTruncationError(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	_, errs := c.GetDelimitedLabel(nil, nil, nil, nil, nil, 10, false)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "label foo-bar-baz exceeds maximum length of 10", errs[0].Error())
}

func TestProviderConfigGetDelimitedLabelWithLocalDelimiter(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	delim := ":"
	actual, errs := c.GetDelimitedLabel(&delim, nil, nil, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo:bar:baz", actual)
}

func TestProviderConfigGetDelimitedLabelWithLocalProperties(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	actual, errs := c.GetDelimitedLabel(nil, []string{"foo", "bar"}, nil, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo-bar", actual)
}

func TestProviderConfigGetDelimitedLabelWithLocalPropertyOrder(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, []string{"bar", "foo"}, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "bar-foo", actual)
}

func TestProviderConfigGetDelimitedLabelWithLocalValues(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, map[string]string{"foo": "bar", "bar": "foo"}, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "bar-foo-baz", actual)
}

func TestProviderConfigGetDelimitedLabelWithLocalCharsReplace(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	regex := "o|a"
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, map[string]string{"foo": "bar", "bar": "foo"}, &regex, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "br-f-bz", actual)
}

func TestProviderConfigGetDelimitedLabelWithAllLocalParams(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	delim := ":"
	actual, errs := c.GetDelimitedLabel(&delim, []string{"foo", "bar"}, []string{}, map[string]string{"foo": "bar", "bar": "foo"}, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "bar:foo", actual)
}

func TestProviderConfigGetTemplatedLabelWithNoLocalValues(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	actual, errs := c.GetTemplatedLabel("{{.foo}}~~{{.bar}}", nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo~~bar", actual)
}

func TestProviderConfigGetTemplatedLabelWithLocalValues(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	actual, errs := c.GetTemplatedLabel("{{.foo}}~~{{.bar}}", map[string]string{"foo": "baz", "bar": "bat"}, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "baz~~bat", actual)
}

func TestProviderConfigGetTemplatedLabelWithLocalCharsReplace(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	regex := "o|a"
	actual, errs := c.GetTemplatedLabel("{{.foo}}~~{{.bar}}", nil, &regex, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "f~~br", actual)
}

func TestProviderConfigGetTags(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	tags, errs := c.GetTags(map[string]string{"foo": "bar", "bar": "baz"}, nil, nil)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, map[string]string{"Foo": "bar", "Bar": "baz", "Baz": "baz"}, tags)
}

func TestProviderConfigGetTagsAsList(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	tags, errs := c.GetTagsAsList(map[string]string{"foo": "bar", "bar": "baz"}, nil, nil)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, []map[string]string{{"Key": "Bar", "Value": "baz"}, {"Key": "Baz", "Value": "baz"}, {"Key": "Foo", "Value": "bar"}}, tags)
}

func TestProviderConfigGetTagsDefaultCase(t *testing.T) {
	c := getDefaultProviderConfig(t, true)
	tags, errs := c.GetTags(map[string]string{"foo": "bar", "bar": "baz"}, nil, nil)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, map[string]string{"Foo": "bar", "Bar": "baz", "Baz": "baz"}, tags)
}
