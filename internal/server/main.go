package server

import "strings"

func validateField(field string) bool {
	if len(strings.TrimSpace(field)) == 0 {
		return false
	}
	return true
}
