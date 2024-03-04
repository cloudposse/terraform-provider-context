package provider

import "github.com/cloudposse/terraform-provider-context/internal/client"

// providerData is returned from the provider's Configure method and is passed to each resource and data source in their
// Configure methods.
type providerData struct {
	contextClient *client.Client
}
