package zoxide

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateModalDimensions(t *testing.T) {
	testcases := []struct {
		name           string
		fullWidth      int
		fullHeight     int
		widthPercent   int
		heightPercent  int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "default 50 percent sizing",
			fullWidth:      100,
			fullHeight:     40,
			widthPercent:   50,
			heightPercent:  50,
			expectedWidth:  50,
			expectedHeight: 20,
		},
		{
			name:           "custom 70 percent width and 60 percent height",
			fullWidth:      100,
			fullHeight:     40,
			widthPercent:   70,
			heightPercent:  60,
			expectedWidth:  70,
			expectedHeight: 24,
		},
		{
			name:           "minimum valid percentages 10 percent",
			fullWidth:      200,
			fullHeight:     100,
			widthPercent:   10,
			heightPercent:  10,
			expectedWidth:  20,
			expectedHeight: 10,
		},
		{
			name:           "maximum valid percentages 100 percent",
			fullWidth:      100,
			fullHeight:     50,
			widthPercent:   100,
			heightPercent:  100,
			expectedWidth:  100,
			expectedHeight: 50,
		},
		{
			name:           "zero percentage falls back to default 50 percent",
			fullWidth:      100,
			fullHeight:     40,
			widthPercent:   0,
			heightPercent:  0,
			expectedWidth:  50,
			expectedHeight: 20,
		},
		{
			name:           "negative percentage falls back to default 50 percent",
			fullWidth:      100,
			fullHeight:     40,
			widthPercent:   -20,
			heightPercent:  -10,
			expectedWidth:  50,
			expectedHeight: 20,
		},
		{
			name:           "over 100 percentage falls back to default 50 percent",
			fullWidth:      100,
			fullHeight:     40,
			widthPercent:   120,
			heightPercent:  150,
			expectedWidth:  50,
			expectedHeight: 20,
		},
		{
			name:           "below minimum width clamps to ZoxideMinWidth",
			fullWidth:      20,
			fullHeight:     20,
			widthPercent:   10,
			heightPercent:  10,
			expectedWidth:  ZoxideMinWidth,
			expectedHeight: ZoxideMinHeight,
		},
		{
			name:           "tiny terminal clamps to minimums",
			fullWidth:      10,
			fullHeight:     2,
			widthPercent:   50,
			heightPercent:  50,
			expectedWidth:  ZoxideMinWidth,
			expectedHeight: ZoxideMinHeight,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w, h := CalculateModalDimensions(tc.fullWidth, tc.fullHeight, tc.widthPercent, tc.heightPercent)
			assert.Equal(t, tc.expectedWidth, w, "modal width mismatch")
			assert.Equal(t, tc.expectedHeight, h, "modal height mismatch")
		})
	}
}
