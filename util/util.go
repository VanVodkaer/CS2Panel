package util

// DefaultIfEmpty 检查值是否为零值，如果是，则返回默认值；如果值不是零值，则返回原值
func DefaultIfEmpty[T comparable](value, defaultValue T) T {
	var zeroValue T
	if value == zeroValue {
		return defaultValue
	}
	return value
}
