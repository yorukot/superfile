package internal

import (
	"testing"
)

// Test KeyMessage type functionality
func TestKeyMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Basic key", "a", "a"},
		{"Arrow key", "↑", "↑"},
		{"Special key", "ctrl+c", "ctrl+c"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyMsg := NewKeyMessage(tt.input)
			if keyMsg.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, keyMsg.String())
			}
		})
	}
}

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

// Helper function to check if direction counts are equal
func assertEqualDirectionCounts(t *testing.T, upCount, downCount int) {
	if upCount != downCount {
		t.Error("Should have equal number of up and down navigation types")
	}
}

// Helper function to validate up direction parity
func validateUpDirections(t *testing.T, directions []NavigationType) {
	for _, upType := range directions {
		if int(upType)%2 != 0 {
			t.Errorf("Up navigation type %v should be even number", upType)
		}
	}
}

// Helper function to validate down direction parity
func validateDownDirections(t *testing.T, directions []NavigationType) {
	for _, downType := range directions {
		if int(downType)%2 != 1 {
			t.Errorf("Down navigation type %v should be odd number", downType)
		}
	}
}

// Test to verify the navigation strategy documentation matches implementation
func TestNavigationStrategyConsistency(t *testing.T) {
	// This test ensures that our documented navigation strategy is consistent
	// with the actual implementation by checking that up/down directions are paired

	upDirections := []NavigationType{NavigateUp, NavigatePageUp}
	downDirections := []NavigationType{NavigateDown, NavigatePageDown}

	// Verify we have equal number of up and down directions
	assertEqualDirectionCounts(t, len(upDirections), len(downDirections))

	// Verify up directions are even numbers, down directions are odd
	validateUpDirections(t, upDirections)
	validateDownDirections(t, downDirections)
}
