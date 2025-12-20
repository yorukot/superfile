package utils

// TODO : replace usage of this with slices.contains
func ArrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
