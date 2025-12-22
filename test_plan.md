# Test Plan for Window Resize & Footer Toggle Fixes

## Context
Recent changes refactored window resizing and footer toggle logic:
- Extracted setHeightValues() and updateComponentDimensions() methods
- Fixed footer toggle to properly update metadata and process bar dimensions
- Separated resize logic into discrete steps for better maintainability

## Critical Test Scenarios

### 1. Window Resize State Consistency
**File**: src/internal/model_test.go  
**Test Name**: TestWindowResizeDimensionConsistency

Test that after window resize:
- Main panel height is correctly calculated
- Footer height adjusts based on terminal size breakpoints
- All file panels receive correct HandleResize() calls
- Metadata and process bar dimensions are updated
- File preview panel dimensions update when visible

**Scenarios**:
- Resize to each height breakpoint (HeightBreakA-D)
- Resize with footer enabled/disabled
- Resize with file preview open/closed

### 2. Footer Toggle Component Updates
**File**: src/internal/model_test.go  
**Test Name**: TestFooterToggleComponentUpdates

Test that toggling footer:
- Updates metadata model dimensions via setMetadataModelSize()
- Updates process bar dimensions via setProcessBarModelSize()  
- Triggers file preview re-render when visible
- Correctly recalculates main panel height

**Critical Assertions**:
- Verify m.metadataModel.width and height after toggle
- Verify m.processBarModel.width and height after toggle
- Verify preview re-render command is returned

### 3. File Panel Search Bar Width
**File**: src/internal/ui/filepanel/navigation_test.go (extend existing)
**Test Name**: TestFilePanelSearchBarResize

Test that search bar width updates correctly:
- After window resize
- With different numbers of panels (1-3)
- With sidebar and preview panel combinations

### 4. Height Breakpoint Calculations
**File**: src/internal/model_test.go  
**Test Name**: TestSetHeightValuesBreakpoints

Test footer height calculation at boundaries:
- Just below/at/above each HeightBreakA-D constant
- With toggleFooter = false (should be 0)
- Edge case: very small terminal (< HeightBreakA)

## Tests NOT Needed (Low Value)

1. **Simple getter/setter tests** - The SetDimensions methods are trivial
2. **UI rendering tests** - Visual output is hard to test and changes frequently
3. **Mouse event handling** - Not affected by recent changes
4. **Help menu sizing** - Separate concern, stable code

## Implementation Notes

### Test Helper Functions Needed
```go
// assertDimensionsConsistent checks all component dimensions match expectations
func assertDimensionsConsistent(t *testing.T, m *model, expectedFooterHeight int) {
    // Check metadata, processbar, file panels, etc.
}

// simulateResize triggers resize and returns any commands
func simulateResize(m *model, width, height int) tea.Cmd {
    return TeaUpdate(m, tea.WindowSizeMsg{Width: width, Height: height})
}
```

### Testing Patterns to Follow
- Use table-driven tests for breakpoint scenarios
- Test state before AND after operations
- Use existing defaultTestModel() helper
- Keep tests focused on behavior, not implementation

## Estimated Test Count
- 4 main test functions
- ~15-20 sub-tests total (table-driven)
- Should add ~200-300 lines of test code

## Priority Order
1. **TestFooterToggleComponentUpdates** - Directly tests the bug fix
2. **TestWindowResizeDimensionConsistency** - Core functionality
3. **TestSetHeightValuesBreakpoints** - Edge cases
4. **TestFilePanelSearchBarResize** - Nice to have
