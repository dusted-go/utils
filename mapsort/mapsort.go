package mapsort

import (
	"sort"
	"strconv"
)

// Keys sorts the keys of a given map in alphabetical order.
func Keys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// KeysByValue sorts the keys of a given map by their value in alphabetical order.
func KeysByValue[K comparable, V any](m map[K]V, getStringValue func(V) string) []K {

	// Invert the map and make values keys and keys values:
	inverted := map[string]K{}
	for k, v := range m {

		// Force each value to be unique:
		value := getStringValue(v)
		i := 0
		for {
			if _, exists := inverted[value]; !exists {
				break
			}
			value = value + strconv.Itoa(i)
			i++
		}

		inverted[value] = k
	}

	// Sort the inverted map's keys:
	sorted := Keys(inverted)

	// Get original keys in order of sorted values:
	keys := make([]K, 0, len(m))
	for _, v := range sorted {
		k := inverted[v]
		keys = append(keys, k)
	}

	return keys
}
