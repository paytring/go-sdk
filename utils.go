package paytring

import "sort"

func MergeMaps[K comparable, V any](m1, m2 map[K]V) map[K]V {
	merged := make(map[K]V)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}

func SortParams(inputMap map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	sortedMap := make(map[string]interface{})
	for _, key := range keys {
		sortedMap[key] = inputMap[key]
	}

	return sortedMap
}
