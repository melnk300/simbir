package utils

func IsEmptyString(s string) any {
	if len(s) == 0 {
		return s
	} else {
		return nil
	}
}
