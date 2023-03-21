package tools

// Map returns a new slice where each element is the result of fn for the corresponding element in the original slice
func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, t := range slice {
		result[i] = fn(t)
	}

	return result
}

// Contains returns true if find appears in slice
func Contains[T comparable](slice []T, find T) bool {
	for _, t := range slice {
		if t == find {
			return true
		}
	}

	return false
}

// IndexOf returns the index of find if it appears in slice. If find is not in slice, -1 will be returned.
func IndexOf[T comparable](slice []T, find T) int {
	for i, t := range slice {
		if t == find {
			return i
		}
	}

	return -1
}

// GroupBy returns a map that is keyed by keySelector and contains a slice of elements returned by valSelector
func GroupBy[T any, K comparable, V any](slice []T, keySelector func(T) K, valSelector func(T) V) map[K][]V {
	grouping := make(map[K][]V)
	for _, t := range slice {
		key := keySelector(t)
		grouping[key] = append(grouping[key], valSelector(t))
	}

	return grouping
}

// ToSet returns a map keyed by keySelector and contains a value of an empty struct
func ToSet[T any, K comparable](slice []T, keySelector func(T) K) *Set[K] {
	set := NewSetWithCapacity[K](len(slice))
	for _, t := range slice {
		set.Add(keySelector(t))
	}

	return set
}

// ToMap return a map that is keyed keySelector and has the value of valSelector for each element in slice.
// If multiple elements return the same key the element that appears later in slice will be chosen.
func ToMap[T any, K comparable, V any](slice []T, keySelector func(T) K, valSelector func(T) V) map[K]V {
	m := make(map[K]V)
	for _, t := range slice {
		m[keySelector(t)] = valSelector(t)
	}

	return m
}
