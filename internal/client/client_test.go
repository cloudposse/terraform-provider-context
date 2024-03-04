package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDefaultClient(t *testing.T, enabled bool) *Client {
	properties := []Property{*NewProperty("foo"), *NewProperty("bar"), *NewProperty("baz")}
	values := map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"}

	options := []func(*Client){}
	if !enabled {
		options = append(options, WithEnabled(false))
	}
	c, err := NewClient(properties, []string{}, values, options...)
	assert.NoError(t, err)

	return c
}
func TestContextClientIsEnabled(t *testing.T) {
	c := getDefaultClient(t, true)
	assert.Equal(t, true, c.IsEnabled())
}

func TestContextClientInitIsNotEnabled(t *testing.T) {
	c := getDefaultClient(t, false)
	assert.Equal(t, false, c.IsEnabled())
}

func TestContextClientGetDelimitedLabelWithNoProperties(t *testing.T) {
	properties := []Property{}
	values := map[string]string{}
	c, err := NewClient(properties, []string{}, values)
	assert.NoError(t, err)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "", actual)
}

func TestContextClientGetDelimitedLabelWithNoLocalOverrides(t *testing.T) {
	c := getDefaultClient(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo-bar-baz", actual)
}

func TestContextClientGetDelimitedLabelWithTruncation(t *testing.T) {
	c := getDefaultClient(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, nil, 10, true)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo-ba6094", actual)
}

func TestContextClientGetDelimitedLabelWithTruncationError(t *testing.T) {
	c := getDefaultClient(t, true)
	_, errs := c.GetDelimitedLabel(nil, nil, nil, nil, 10, false)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "label foo-bar-baz exceeds maximum length of 10", errs[0].Error())
}

func TestContextClientGetDelimitedLabelWithLocalDelimiter(t *testing.T) {
	c := getDefaultClient(t, true)
	delim := ":"
	actual, errs := c.GetDelimitedLabel(&delim, nil, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo:bar:baz", actual)
}

func TestContextClientGetDelimitedLabelWithLocalProperties(t *testing.T) {
	c := getDefaultClient(t, true)
	actual, errs := c.GetDelimitedLabel(nil, []string{"foo", "bar"}, nil, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo-bar", actual)
}

func TestContextClientGetDelimitedLabelWithLocalPropertyOrder(t *testing.T) {
	c := getDefaultClient(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, []string{"bar", "foo"}, nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "bar-foo", actual)
}

func TestContextClientGetDelimitedLabelWithLocalValues(t *testing.T) {
	c := getDefaultClient(t, true)
	actual, errs := c.GetDelimitedLabel(nil, nil, nil, map[string]string{"foo": "bar", "bar": "foo"}, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "bar-foo-baz", actual)
}

func TestContextClientGetDelimitedLabelWithAllLocalParams(t *testing.T) {
	c := getDefaultClient(t, true)
	delim := ":"
	actual, errs := c.GetDelimitedLabel(&delim, []string{"foo", "bar"}, []string{}, map[string]string{"foo": "bar", "bar": "foo"}, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "bar:foo", actual)
}

func TestContextClientGetTemplatedLabelWithNoLocalValues(t *testing.T) {
	c := getDefaultClient(t, true)
	actual, errs := c.GetTemplatedLabel("{{.foo}}~~{{.bar}}", nil, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "foo~~bar", actual)
}

func TestContextClientGetTemplatedLabelWithLocalValues(t *testing.T) {
	c := getDefaultClient(t, true)
	actual, errs := c.GetTemplatedLabel("{{.foo}}~~{{.bar}}", map[string]string{"foo": "baz", "bar": "bat"}, 0, false)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "baz~~bat", actual)
}

func TestContextClientGetTags(t *testing.T) {
	c := getDefaultClient(t, true)
	tags, errs := c.GetTags(map[string]string{"foo": "bar", "bar": "baz"})
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, map[string]string{"foo": "bar", "bar": "baz", "baz": "baz"}, tags)
}

func TestContextClientGetTagsAsList(t *testing.T) {
	c := getDefaultClient(t, true)
	tags, errs := c.GetTagsAsList(map[string]string{"foo": "bar", "bar": "baz"})
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, []map[string]string{{"Key": "bar", "Value": "baz"}, {"Key": "baz", "Value": "baz"}, {"Key": "foo", "Value": "bar"}}, tags)
}
