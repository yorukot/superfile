package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

// Assuming both a and b are string in format of
// vX.Y.Z.W .... (any number of digits allowed)
// Return 1 if a > b
// Return 0 if a == b
// Return -1 if a < b
// Return non-nil error if string are not correctly formated
func versionCompare(a string, b string) (int, error) {
	res := 0
	if len(a) < 2 || len(b) < 2 || a[0] != 'v' || b[0] != 'v' {
		return res, fmt.Errorf("Invalid version strings %v and %v", a, b)
	}

	a_parts := strings.Split(strings.TrimPrefix(a, "v"), ".")
	b_parts := strings.Split(strings.TrimPrefix(b, "v"), ".")
	curIdx := 0
	for curIdx < len(a_parts) && curIdx < len(b_parts) {
		aVal, bVal := 0, 0
		aVal, err := strconv.Atoi(a_parts[curIdx])
		if err != nil || aVal < 0 {
			return res, fmt.Errorf("Non positive integer %v in version : %w", a_parts[curIdx], err)
		}
		bVal, err = strconv.Atoi(b_parts[curIdx])
		if err != nil || bVal < 0 {
			return res, fmt.Errorf("Non positive integer %v in version : %w", b_parts[curIdx], err)
		}
		if aVal > bVal {
			return 1, nil
		} else if aVal < bVal {
			return -1, nil
		}
		// Otherwise continue iteration
		curIdx++
	}

	if curIdx < len(a_parts) {
		// some parts of a are still left, while b is completely iterated
		// Just make sure they are all integers
		for curIdx < len(a_parts) {
			if aVal, err := strconv.Atoi(a_parts[curIdx]); err != nil || aVal < 0 {
				return res, fmt.Errorf("Non integer part %v in version : %w", a_parts[curIdx], err)
			}
			curIdx++
		}
		return 1, nil
	}

	if curIdx < len(b_parts) {
		// some parts of b are still left, while a is completely iterated
		// Just make sure they are all integers
		for curIdx < len(b_parts) {
			if bVal, err := strconv.Atoi(b_parts[curIdx]); err != nil || bVal < 0 {
				return res, fmt.Errorf("Non integer part %v in version : %w", b_parts[curIdx], err)
			}
			curIdx++
		}
		return -1, nil
	}

	// Both a and b are completely iterated
	return 0, nil

}
