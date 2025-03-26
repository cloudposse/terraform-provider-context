package stringHelpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(s string) string {
	// Hash the string using SHA256
	hasher := sha256.New()
	hasher.Write([]byte(s))
	hashBytes := hasher.Sum(nil)

	// Convert the hash bytes to a hex string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
