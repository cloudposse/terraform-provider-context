package mapHelpers

import "testing"

func TestHashMap(t *testing.T) {
	// Test case 1: Empty map
	emptyMap := make(map[string]string)
	expectedHashEmpty := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // SHA256 hash of empty string
	hashEmpty := HashMap(emptyMap)
	if hashEmpty != expectedHashEmpty {
		t.Errorf("Expected hash for empty map: %s, got: %s", expectedHashEmpty, hashEmpty)
	}

	// Test case 2: Non-empty map
	myMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	expectedHash := "bd232b1694ab90cba0824e5decd0259d8e322a8fc0b1630ec6dfcce8fcb0107d" // SHA256 hash of "key1:value1;key2:value2;key3:value3;"
	hash := HashMap(myMap)
	if hash != expectedHash {
		t.Errorf("Expected hash for non-empty map: %s, got: %s", expectedHash, hash)
	}
}
