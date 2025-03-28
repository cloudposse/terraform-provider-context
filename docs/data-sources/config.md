---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "context_config Data Source - terraform-provider-context"
subcategory: ""
description: |-
  Context Config data source
---

# context_config (Data Source)

Context Config data source



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) Config identifier

### Read-Only

- `delimiter` (String) Delimiter to use when creating the label from properties. Conflicts with `template`.
- `enabled` (Boolean) Flag to indicate if the config is enabled.
- `properties` (Attributes Map) A map of properties to use for labels created by the provider. (see [below for nested schema](#nestedatt--properties))
- `property_order` (List of String) A list of properties to use for labels created by the provider.
- `replace_chars_regex` (String) Regex to use for replacing characters in labels created by the provider.
- `tags_key_case` (String) Case to use for keys in tags created by the provider.
- `tags_value_case` (String) Case to use for values in tags created by the provider.
- `values` (Map of String) A map of values to use for labels created by the provider.

<a id="nestedatt--properties"></a>
### Nested Schema for `properties`

Optional:

- `include_in_tags` (Boolean) A flag to indicate if the property should be included in tags.
- `max_length` (Number) The maximum length of the property.
- `min_length` (Number) The minimum length of the property.
- `required` (Boolean) A flag to indicate if the property is required.
- `tags_key_case` (String) The case to use for the key of this property in tags.
- `tags_value_case` (String) The case to use for the value of this property in tags.
- `validation_regex` (String) A regular expression to validate the property.
