package jsonata

import (
	"errors"
	"reflect"
	"sort"
)

type visitedMap map[uintptr]bool

func makeDeterministic(input interface{}, visited visitedMap) (interface{}, error) {
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
			return nil, errors.New("circular dependency detected in output - please investigate")
		}
		visited[value.Pointer()] = true

		return makeDeterministicMap(value, visited)
	case reflect.Slice:
		return makeDeterministicArray(value, visited)
	default:
		return input, nil
	}
}

func makeDeterministicMap(input reflect.Value, visited visitedMap) (map[string]interface{}, error) {
	keys := make([]string, 0, input.Len())
	for _, key := range input.MapKeys() {
		keys = append(keys, key.String())
	}
	sort.Strings(keys)

	deterministicMap := make(map[string]interface{})
	for _, key := range keys {
		value, err := makeDeterministic(input.MapIndex(reflect.ValueOf(key)).Interface(), visited)
		if err != nil {
			return nil, err
		}

		if value != nil {
			deterministicMap[key] = value
		}
	}

	return deterministicMap, nil
}

func makeDeterministicArray(input reflect.Value, visited visitedMap) ([]interface{}, error) {
	deterministicArray := make([]interface{}, input.Len())
	var err error
	for i := 0; i < input.Len(); i++ {
		deterministicArray[i], err = makeDeterministic(input.Index(i).Interface(), visited)
		if err != nil {
			return nil, err
		}
	}

	return deterministicArray, nil
}
