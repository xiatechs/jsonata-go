package jsonata

import (
	"sort"
)

func sortMap(inputMap map[string]interface{}) map[string]interface{} {
	sortedMap := make(map[string]interface{})

	// Extract and sort keys
	keys := make([]string, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Recursively sort nested maps and arrays
	for _, key := range keys {
		value := inputMap[key]

		switch typedValue := value.(type) {
		case map[string]interface{}:
			sortedMap[key] = sortMap(typedValue)
		case []interface{}:
			sortedMap[key] = sortArray(typedValue)
		default:
			sortedMap[key] = value
		}
	}

	return sortedMap
}

func sortArray(inputArray []interface{}) []interface{} {
	sortedArray := make([]interface{}, len(inputArray))

	// Recursively sort elements of the array
	for i, element := range inputArray {
		switch typedElement := element.(type) {
		case map[string]interface{}:
			sortedArray[i] = sortMap(typedElement)
		case []interface{}:
			sortedArray[i] = sortArray(typedElement)
		default:
			sortedArray[i] = element
		}
	}

	return sortedArray
}
