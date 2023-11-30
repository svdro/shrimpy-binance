package defaults

import "sort"

// filterZeros is a helper function that filters out all zeros from a slice of
// int64 values This is useful for slices that are intialized with zeros.
func filterZeros(values []int64) []int64 {
	var filtered []int64
	for _, value := range values {
		if value != 0 {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

// mean is a helper function that calculates the mean of a slice of int64 values.
func mean(values []int64) (int64, bool) {
	n := len(values)
	if n == 0 {
		return 0, false
	}

	sum := int64(0)
	for _, value := range values {
		sum += value
	}
	return sum / int64(n), true
}

// median is a helper function that calculates the median of a slice of int64 values.
func median(values []int64) (int64, bool) {
	n := len(values)
	if n == 0 {
		return 0, false
	}

	valuesCopy := append([]int64(nil), values...)
	sort.Slice(valuesCopy, func(i, j int) bool { return valuesCopy[i] < valuesCopy[j] })

	mid := n / 2
	if n%2 == 0 {
		return (valuesCopy[mid-1] + valuesCopy[mid]) / 2, true
	}
	return valuesCopy[mid], true
}
