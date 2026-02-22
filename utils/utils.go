package utils

func SetPtr[T any](value T) *T {
	result := new(T)

	*result = value

	return result
}

func GetPtr[T any](value *T) T {
	if value == nil {
		var zeroValue T
		return zeroValue
	}

	return *value
}
