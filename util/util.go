package util

func DefaultIfEmpty[T comparable](value, defaultValue T) T {
	var zeroValue T
	if value == zeroValue {
		return defaultValue
	}
	return value
}
