package jsonata

import (
	"reflect"
	"sort"
)

type visitedMap map[interface{}]bool

func makeDeterministic(input interface{}, visited visitedMap) interface{} {
	if visited == nil {
		visited = make(visitedMap)
	}

	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Map:
		if visited[value.Pointer()] {
			// here we kill circular dependencies i.e they don't appear in the JSON
			// but is this how we want to handle this? do we want to return an error
			return nil
		}
		visited[value.Pointer()] = true

		return makeDeterministicMap(value, visited)
	case reflect.Slice:
		return makeDeterministicArray(value, visited)
	default:
		return input
	}
}

func makeDeterministicMap(input reflect.Value, visited visitedMap) map[string]interface{} {
	keys := make([]string, 0, input.Len())
	for _, key := range input.MapKeys() {
		keys = append(keys, key.String())
	}
	sort.Strings(keys)

	deterministicMap := make(map[string]interface{})
	for _, key := range keys {
		value := makeDeterministic(input.MapIndex(reflect.ValueOf(key)).Interface(), visited)
		if value != nil {
			deterministicMap[key] = value
		}
	}

	return deterministicMap
}

func makeDeterministicArray(input reflect.Value, visited visitedMap) []interface{} {
	deterministicArray := make([]interface{}, input.Len())
	for i := 0; i < input.Len(); i++ {
		deterministicArray[i] = makeDeterministic(input.Index(i).Interface(), visited)
	}

	return deterministicArray
}