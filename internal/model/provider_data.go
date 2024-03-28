package model

// ProviderData is returned from the provider's Configure method and is passed to each resource and data source in their
// Configure methods.
type ProviderData struct {
	ProviderConfig *ProviderConfig
}
