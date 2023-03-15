package utils

func Contains[T comparable](value T, values []T) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}

	return false
}
