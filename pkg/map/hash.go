package mapHelpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashMap(m map[string]string) string {
	// Convert the map to a string representation
	var str string
	for k, v := range m {
		str += k + ":" + v + ";"
	}

	// Hash the string using SHA256
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hashBytes := hasher.Sum(nil)

	// Convert the hash bytes to a hex string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
