package mapHelpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

func HashMap(m interface{}) string {
	switch v := m.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		h := sha256.New()
		for _, k := range keys {
			h.Write([]byte(k))
			h.Write([]byte(HashMap(v[k])))
		}
		return hex.EncodeToString(h.Sum(nil))
	default:
		// For primitive types, convert to JSON and hash
		b, _ := json.Marshal(v)
		hashed := sha256.Sum256(b)
		return hex.EncodeToString(hashed[:])
	}
}
