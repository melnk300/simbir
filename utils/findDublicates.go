package utils

func FindDublicates[T comparable](values []T) bool {
	seen := make(map[T]bool)
	for _, value := range values {
		if seen[value] {
			return true
		} else {
			seen[value] = true
		}
	}

	return false
}
