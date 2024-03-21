package configHelpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/cloudposse/terraform-provider-context/internal/client"
)

func sortJSON(jsonData []byte) ([]byte, error) {
	// Unmarshal JSON into a map
	var jsonObj map[string]interface{}
	err := json.Unmarshal(jsonData, &jsonObj)
	if err != nil {
		return nil, err
	}

	// Extract keys from map and sort them
	keys := make([]string, 0, len(jsonObj))
	for key := range jsonObj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create a new map with sorted keys
	sortedJSON := make(map[string]interface{})
	for _, key := range keys {
		sortedJSON[key] = jsonObj[key]
	}

	// Marshal the sorted map back to JSON
	sortedData, err := json.Marshal(sortedJSON)
	if err != nil {
		return nil, err
	}

	return sortedData, nil
}

func HashConfig(delimiter string, enabled bool, properties []client.Property, values map[string]string) (string, error) {
	type config struct {
		Delimiter  string            `json:"delimiter"`
		Enabled    bool              `json:"enabled"`
		Properties []client.Property `json:"properties"`
		Values     map[string]string `json:"values"`
	}

	c := config{
		Delimiter:  delimiter,
		Enabled:    enabled,
		Properties: properties,
		Values:     values,
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	sortedJSON, err := sortJSON(jsonData)
	if err != nil {
		return "", err
	}

	// Hash the string using SHA256
	hasher := sha256.New()
	hasher.Write(sortedJSON)
	hashBytes := hasher.Sum(nil)

	// Convert the hash bytes to a hex string
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
