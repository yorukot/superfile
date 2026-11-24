package utils

// IsASCIIDigit reports whether b is an ASCII digit ('0'-'9').
func IsASCIIDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// LeadingASCIIDigits returns the leading ASCII digit prefix of s.
func LeadingASCIIDigits(s string) string {
	for i := range len(s) {
		if !IsASCIIDigit(s[i]) {
			return s[:i]
		}
	}
	return s
}
