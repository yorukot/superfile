package internal

import (
	"testing"
)

// Test the NavigationType enum values
func TestNavigationType(t *testing.T) {
	tests := []struct {
		name     string
		navType  NavigationType
		expected int
	}{
		{"NavigateUp", NavigateUp, 0},
		{"NavigateDown", NavigateDown, 1},
		{"NavigatePageUp", NavigatePageUp, 2},
		{"NavigatePageDown", NavigatePageDown, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.navType) != tt.expected {
				t.Errorf("Expected %s to be %d, got %d", tt.name, tt.expected, int(tt.navType))
			}
		})
	}
}

// Test NavigationType constants are properly defined
func TestNavigationTypeConstants(t *testing.T) {
	// Ensure all navigation types are different
	types := []NavigationType{NavigateUp, NavigateDown, NavigatePageUp, NavigatePageDown}

	for i := 0; i < len(types); i++ {
		for j := i + 1; j < len(types); j++ {
			if types[i] == types[j] {
				t.Errorf("NavigationType constants must be unique: %v == %v", types[i], types[j])
			}
		}
	}
}

// Test that navigation types follow expected ordering
func TestNavigationTypeOrdering(t *testing.T) {
	if NavigateUp >= NavigateDown {
		t.Error("NavigateUp should be less than NavigateDown")
	}
	if NavigateDown >= NavigatePageUp {
		t.Error("NavigateDown should be less than NavigatePageUp")
	}
	if NavigatePageUp >= NavigatePageDown {
		t.Error("NavigatePageUp should be less than NavigatePageDown")
	}
}

// Mock test to verify the navigation strategy documentation matches implementation
func TestNavigationStrategyConsistency(t *testing.T) {
	// This test ensures that our documented navigation strategy is consistent
	// with the actual implementation by checking that up/down directions are paired

	upDirections := []NavigationType{NavigateUp, NavigatePageUp}
	downDirections := []NavigationType{NavigateDown, NavigatePageDown}

	// Verify we have equal number of up and down directions
	if len(upDirections) != len(downDirections) {
		t.Error("Should have equal number of up and down navigation types")
	}

	// Verify up directions are even numbers, down directions are odd
	for _, upType := range upDirections {
		if int(upType)%2 != 0 {
			t.Errorf("Up navigation type %v should be even number", upType)
		}
	}

	for _, downType := range downDirections {
		if int(downType)%2 != 1 {
			t.Errorf("Down navigation type %v should be odd number", downType)
		}
	}
}

