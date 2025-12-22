# Phase 1 Detailed Implementation Plan - Filepanel Dimension Management

## Step 1: Add Dimension Fields to Model

### 1.1 Update Model struct
**File:** `src/internal/ui/filepanel/types.go`
**Line:** 14-31 (inside Model struct)

**Current:**
```go
type Model struct {
	Cursor      int
	RenderIndex int
	IsFocused   bool
	Location    string
	// ... other fields ...
}
```

**Change to:**
```go
type Model struct {
	Cursor      int
	RenderIndex int
	IsFocused   bool
	Location    string
	// Dimension fields
	width       int // Total width including borders
	height      int // Total height including borders
	// ... other fields ...
}
```

## Step 2: Add Dimension Constants

### 2.1 Add constants for minimum dimensions
**File:** `src/internal/common/ui_consts.go`
**Line:** After line 15 (after DefaultFilePanelWidth)

**Add:**
```go
FilePanelMinWidth  = 20  // Minimum width for file panel
FilePanelMinHeight = 5   // Minimum height for file panel
```

## Step 3: Add Dimension Methods

### 3.1 Create dimension management methods
**File:** `src/internal/ui/filepanel/model_utils.go`
**Location:** Add at the end of file

**Add these methods:**
```go
// UpdateDimensions sets the panel dimensions with validation
func (m *Model) UpdateDimensions(width, height int) {
	if width < common.FilePanelMinWidth {
		width = common.FilePanelMinWidth
	}
	if height < common.FilePanelMinHeight {
		height = common.FilePanelMinHeight
	}
	m.width = width
	m.height = height
}

// GetWidth returns the total panel width
func (m *Model) GetWidth() int {
	return m.width
}

// GetHeight returns the total panel height
func (m *Model) GetHeight() int {
	return m.height
}

// GetMainPanelHeight returns content height (total height minus borders)
func (m *Model) GetMainPanelHeight() int {
	return m.height - common.BorderPadding
}

// GetContentWidth returns content width (total width minus borders)
func (m *Model) GetContentWidth() int {
	return m.width - common.BorderPadding
}
```

## Step 4: Convert panelElementHeight to Receiver Method

### 4.1 Delete standalone function and create receiver method
**File:** `src/internal/ui/filepanel/misc_utils.go`

**Delete (lines 108-110):**
```go
func panelElementHeight(mainPanelHeight int) int {
	return mainPanelHeight - common.PanelPadding
}
```

**File:** `src/internal/ui/filepanel/model_utils.go`
**Add after GetContentWidth method:**
```go
// PanelElementHeight calculates the number of visible elements in content area
func (m *Model) PanelElementHeight() int {
	contentHeight := m.GetMainPanelHeight()
	if m.SearchBar != nil {
		contentHeight -= 2  // Search bar takes 2 lines
	}
	return contentHeight - common.PanelPadding
}
```

### 4.2 Update getPageScrollSize to use receiver
**File:** `src/internal/ui/filepanel/misc_utils.go`
**Line:** ~155-162 (getPageScrollSize function)

**Current:**
```go
func getPageScrollSize(mainPanelHeight int) int {
	scrollSize := common.Config.PageScrollSize
	if scrollSize == 0 {
		scrollSize = panelElementHeight(mainPanelHeight)
	}
	return scrollSize
}
```

**Change to:**
```go
func (m *Model) getPageScrollSize() int {
	scrollSize := common.Config.PageScrollSize
	if scrollSize == 0 {
		scrollSize = m.PanelElementHeight()
	}
	return scrollSize
}
```

## Step 5: Update Render Methods

### 5.1 Update main Render method
**File:** `src/internal/ui/filepanel/render.go`
**Line:** 17-25

**Current:**
```go
func (m *Model) Render(mainPanelHeight int, filePanelWidth int, focussed bool) string {
	r := ui.FilePanelRenderer(mainPanelHeight+common.BorderPadding, filePanelWidth+common.BorderPadding, focussed)

	m.renderTopBar(r, filePanelWidth)
	m.renderSearchBar(r)
	m.renderFooter(r, uint(len(m.Selected)))
	m.renderFileEntries(r, mainPanelHeight, filePanelWidth)

	return r.Render()
}
```

**Change to:**
```go
func (m *Model) Render(focussed bool) string {
	r := ui.FilePanelRenderer(m.height, m.width, focussed)

	m.renderTopBar(r)
	m.renderSearchBar(r)
	m.renderFooter(r, uint(len(m.Selected)))
	m.renderFileEntries(r)

	return r.Render()
}
```

### 5.2 Update renderTopBar method
**File:** `src/internal/ui/filepanel/render.go`
**Line:** 28

**Current:**
```go
func (m *Model) renderTopBar(r *rendering.Renderer, filePanelWidth int)
```

**Change to:**
```go
func (m *Model) renderTopBar(r *rendering.Renderer)
```

**Inside the method:** Replace all uses of `filePanelWidth` with `m.GetContentWidth()`

### 5.3 Update renderFileEntries method
**File:** `src/internal/ui/filepanel/render.go`
**Line:** 64-65

**Current:**
```go
func (m *Model) renderFileEntries(r *rendering.Renderer, mainPanelHeight, filePanelWidth int) {
	end := min(m.RenderIndex+panelElementHeight(mainPanelHeight), len(m.Element))
```

**Change to:**
```go
func (m *Model) renderFileEntries(r *rendering.Renderer) {
	end := min(m.RenderIndex+m.PanelElementHeight(), len(m.Element))
```

**Inside the method:** 
- Replace all uses of `mainPanelHeight` with `m.GetMainPanelHeight()`
- Replace all uses of `filePanelWidth` with `m.GetContentWidth()`

## Step 6: Update Navigation Methods

### 6.1 Update PgUp and PgDown
**File:** `src/internal/ui/filepanel/navigation.go`

**Lines ~40-60 (PgUp method):**
```go
// Current:
func (m *Model) PgUp(mainPanelHeight int) {
	scrollSize := getPageScrollSize(mainPanelHeight)
	m.moveCursorBy(-scrollSize, mainPanelHeight)
}

// Change to:
func (m *Model) PgUp() {
	scrollSize := m.getPageScrollSize()
	m.moveCursorBy(-scrollSize)
}
```

**Similarly for PgDown:**
```go
// Current:
func (m *Model) PgDown(mainPanelHeight int) {
	scrollSize := getPageScrollSize(mainPanelHeight)
	m.moveCursorBy(scrollSize, mainPanelHeight)
}

// Change to:
func (m *Model) PgDown() {
	scrollSize := m.getPageScrollSize()
	m.moveCursorBy(scrollSize)
}
```

### 6.2 Update scrollToCursor
**File:** `src/internal/ui/filepanel/navigation.go`
**Lines ~100-120**

**Current signature:**
```go
func (m *Model) scrollToCursor(mainPanelHeight int)
```

**Change to:**
```go
func (m *Model) scrollToCursor()
```

**Inside method:** Replace `panelElementHeight(mainPanelHeight)` with `m.PanelElementHeight()`

### 6.3 Update HandleResize
**File:** `src/internal/ui/filepanel/navigation.go`

**Current:**
```go
func (m *Model) HandleResize(mainPanelHeight int)
```

**Change to:**
```go
func (m *Model) HandleResize()
```

**Inside method:** Use `m.scrollToCursor()` instead of `m.scrollToCursor(mainPanelHeight)`

### 6.4 Update moveCursorBy
**File:** `src/internal/ui/filepanel/navigation.go`

**Current:**
```go
func (m *Model) moveCursorBy(delta int, mainPanelHeight int)
```

**Change to:**
```go
func (m *Model) moveCursorBy(delta int)
```

**Inside method:** Use `m.scrollToCursor()` instead of `m.scrollToCursor(mainPanelHeight)`

### 6.5 Update ItemSelectUp and ItemSelectDown
**File:** `src/internal/ui/filepanel/navigation.go`

**Current:**
```go
func (m *Model) ItemSelectUp(mainPanelHeight int) {
	// ... code ...
	m.moveCursorBy(-1, mainPanelHeight)
}

func (m *Model) ItemSelectDown(mainPanelHeight int) {
	// ... code ...
	m.moveCursorBy(1, mainPanelHeight)
}
```

**Change to:**
```go
func (m *Model) ItemSelectUp() {
	// ... code ...
	m.moveCursorBy(-1)
}

func (m *Model) ItemSelectDown() {
	// ... code ...
	m.moveCursorBy(1)
}
```

## Step 7: Initialize Dimensions in Constructor

Since there's no explicit constructor for filepanel.Model in the filepanel package,
we need to ensure dimensions are set when panels are created.

**Note:** The actual panel creation happens in `src/internal/` files. We'll need to:
1. Find where `filepanel.Model{}` is created
2. Immediately call `UpdateDimensions` after creation
3. Set initial values to minimum dimensions

**Example pattern:**
```go
panel := filepanel.Model{
	// ... existing fields ...
}
panel.UpdateDimensions(common.FilePanelMinWidth, common.FilePanelMinHeight)
```

## Compilation Check Points

After each step, ensure compilation:
```bash
go build -o bin/spf ./src/cmd
```

## Key Changes Summary

1. **Add fields:** `width`, `height` to Model struct
2. **Add constants:** `FilePanelMinWidth`, `FilePanelMinHeight`
3. **Add methods:** `UpdateDimensions`, `GetWidth`, `GetHeight`, `GetMainPanelHeight`, `GetContentWidth`, `PanelElementHeight`
4. **Remove parameters:** All `mainPanelHeight` and `filePanelWidth` parameters from methods
5. **Convert to receiver:** `panelElementHeight` → `m.PanelElementHeight()`
6. **Use receiver methods:** Access dimensions via `m.GetWidth()`, `m.GetHeight()`, etc.

## Files to Modify

1. `src/internal/ui/filepanel/types.go` - Add dimension fields
2. `src/internal/common/ui_consts.go` - Add dimension constants
3. `src/internal/ui/filepanel/model_utils.go` - Add dimension methods
4. `src/internal/ui/filepanel/misc_utils.go` - Delete panelElementHeight, update getPageScrollSize
5. `src/internal/ui/filepanel/render.go` - Update all render methods
6. `src/internal/ui/filepanel/navigation.go` - Update all navigation methods

## Testing Requirements

After Phase 1 completion, all methods should:
1. Access dimensions from the Model struct
2. Not require dimension parameters
3. Use receiver methods for calculations
4. Validate dimensions through UpdateDimensions
