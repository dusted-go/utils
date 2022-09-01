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
func KeysByValue[T comparable](m map[T]string) []T {

	// Invert the map and make values keys and keys values:
	inverted := map[string]T{}
	for k, v := range m {

		// Force each value to be unique:
		value := v
		i := 0
		for {
			if _, exists := inverted[v]; !exists {
				break
			}
			v = value + strconv.Itoa(i)
			i++
		}

		inverted[v] = k
	}

	// Sort the inverted map's keys:
	sorted := Keys(inverted)

	// Get original keys in order of sorted values:
	keys := make([]T, 0, len(m))
	for _, v := range sorted {
		k := inverted[v]
		keys = append(keys, k)
	}

	return keys
}
