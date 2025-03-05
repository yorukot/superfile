package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_versionCompare(t *testing.T) {
	data := []struct {
		vA             string
		vB             string
		expectedError  bool
		expectedResult int
		description    string
	}{
		{"v1.2.0", "v1.3.0", false, -1, "Basic version comparison (a < b)"},
		{"v1.2.0", "v1.1.7.1", false, 1, "Version with more parts is greater"},
		{"v1.2.0", "v3", false, -1, "Single digit version comparison"},
		{"v1.7.1", "v1.7.1", false, 0, "Exact version match"},
		{"v4", "v4", false, 0, "Single digit version match"},
		{"v4", "v5", false, -1, "Single digit version comparison"},
		{"v5", "v4", false, 1, "Single digit version comparison (reverse)"},
		{"v5.1", "v5", false, 1, "Version with additional part is greater"},
		{"v5", "v5.1", false, -1, "Shorter version is lesser"},

		// Error cases
		{"1.7.1", "v1.7.1", true, 0, "Missing 'v' prefix (first version)"},
		{"v1.7a.1", "v1.7.1", true, 0, "Non-numeric part in version"},
		{"v1.7.1.", "v1.7.1", true, 0, "Trailing dot in version"},
		{"v1.7.1", "v1.7.1.", true, 0, "Trailing dot in second version"},
		{"v", "v1", true, 0, "Incomplete version string"},
		{"v1.-1.0", "v1.0.0", true, 0, "Negative number in version"},
		{"v1.2..3", "v1.2.3", true, 0, "Double dot in version"},

		{"v1.0.0", "v1.0.1", false, -1, "Smallest difference in last part"},
		{"v10.0.0", "v2.0.0", false, 1, "Multi-digit version comparison"},
		{"v1.2.3.4.5", "v1.2.3.4.6", false, -1, "Many version parts comparison"},
		{"v1.2.3.4.5", "v1.2.3.4.5", false, 0, "Many version parts exact match"},
		{"v0.1.0", "v1.0.0", false, -1, "Zero to non-zero version comparison"},
		{"v0", "v0.0.1", false, -1, "Zero version with additional parts"},
		{"v0.0.0", "v0", false, 1, "Multiple zero representations"},
		{"v01.2.3", "v1.2.3", false, 0, "Leading zero in version part"},
		{"v1.2.03", "v1.2.3", false, 0, "Leading zero in version part"},
	}

	for _, tt := range data {
		t.Run(tt.description, func(t *testing.T) {
			res, err := versionCompare(tt.vA, tt.vB)
			if tt.expectedError {
				assert.NotNil(t, err, "Error is expected for %s", tt.description)
			} else {
				assert.Nil(t, err, "Error should be Nil for %s", tt.description)
				assert.Equal(t, tt.expectedResult, res, "Result should be as expected for %s", tt.description)
			}
		})
	}
}
