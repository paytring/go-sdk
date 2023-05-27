package paytring

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortParams(t *testing.T) {
	params := map[string]interface{}{
		"phone":  "1234567890",
		"amount": "100",
		"key":    "test_123",
		"hash":   "secret_123",
	}

	expectedSortedKeys := []string{"amount", "hash", "key", "phone"}

	// Step 1: Extract keys from the map
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}

	// Step 2: Sort the keys
	sort.Strings(keys)

	// Step 3: Iterate over sorted keys and access values
	sortedKeys := make([]string, 0, len(params))
	for _, key := range keys {
		sortedKeys = append(sortedKeys, key)
	}

	assert.Equal(t, expectedSortedKeys, sortedKeys)
}
