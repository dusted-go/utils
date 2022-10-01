package mapsort

import (
	"sort"
)

type Sortable interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | string
}

// Keys sorts the keys of a given map in ascending order.
func Keys[K Sortable, T any](m map[K]T) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// KeysDescending sorts the keys of a given map in descending order.
func KeysDescending[K Sortable, T any](m map[K]T) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return keys
}

// KeysByValue sorts the keys of a given map by their sortable value in ascending order.
// The getSortable function is used to extract a sortable value from a map's item.
func KeysByValue[K comparable, V any, S Sortable](m map[K]V, getSortable func(V) S) []K {

	// Invert the map and make values keys and keys values:
	inverted := []struct {
		Key   S
		Value K
	}{}

	for k, v := range m {
		inverted = append(inverted, struct {
			Key   S
			Value K
		}{
			Key:   getSortable(v),
			Value: k,
		})
	}

	// Sort the inverted array:
	sort.Slice(inverted, func(i, j int) bool {
		return inverted[i].Key < inverted[j].Key
	})

	// Get original keys in order of sorted values:
	keys := make([]K, 0, len(m))
	for _, i := range inverted {
		k := i.Value
		keys = append(keys, k)
	}

	return keys
}
