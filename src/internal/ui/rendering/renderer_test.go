package rendering

import (
	"flag"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/utils"
)

const (
	sectionStr = "<SECTION>"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	m.Run()
}

func getDefaultTestRendererConfig(totalHeight int, totalWidth int, borderRequired bool,
	truncateHeight bool) RendererConfig {
	cfg := DefaultRendererConfig(totalHeight, totalWidth)
	if borderRequired {
		cfg.BorderRequired = true
		cfg.Border = lipgloss.Border{
			Top:    "─",
			Bottom: "─",
			Left:   "│",
			Right:  "│",

			TopLeft:  "╭",
			TopRight: "╮",

			BottomLeft:  "╰",
			BottomRight: "╯",

			MiddleLeft:  "├",
			MiddleRight: "┤",
		}
	}
	cfg.TruncateHeight = truncateHeight
	return cfg
}

func getDefaultTestRenderer(totalHeight int, totalWidth int, borderRequired bool) *Renderer {
	r, _ := NewRenderer(getDefaultTestRendererConfig(totalHeight, totalWidth, borderRequired, false))
	return r
}

func TestRendererBasic(t *testing.T) {
	t.Run("Basic test", func(t *testing.T) {
		r := getDefaultTestRenderer(4, 4, true)
		r.AddLines("L1")
		r.AddLines("L2--Extra line should truncated")
		r.AddLines("L3--Extra line should not be added")
		res := r.Render()
		expected := "" +
			"╭──╮\n" +
			"│L1│\n" +
			"│L2│\n" +
			"╰──╯"
		assert.Equal(t, expected, res)
	})

	t.Run("Empty Renderer", func(t *testing.T) {
		r := getDefaultTestRenderer(0, 0, false)
		r.AddLines("L1")
		r.AddLines("L2--Extra line should truncated")
		r.AddLines("L3--Extra line should not be added")
		res := r.Render()
		expected := ""
		assert.Equal(t, expected, res)
	})

	t.Run("Invalid config Renderer", func(t *testing.T) {
		cfg := getDefaultTestRendererConfig(0, 0, true, false)
		r, err := NewRenderer(cfg)
		assert.Nil(t, r)
		require.Error(t, err)
	})
}

func TestSections(t *testing.T) {
	sectionTests := []struct {
		name           string
		totalHeight    int
		totalWidth     int
		borderRequired bool
		// Test expects only single line strings.
		lines         []string
		trucateheight bool
		expected      string
	}{
		{
			name:           "Basic Sections",
			totalHeight:    7,
			totalWidth:     4,
			borderRequired: true,
			lines:          []string{"L1", sectionStr, "L2", sectionStr, sectionStr, "L3", sectionStr},
			trucateheight:  false,
			expected: "" +
				"╭──╮\n" +
				"│L1│\n" +
				"├──┤\n" +
				"│L2│\n" +
				"├──┤\n" +
				"├──┤\n" +
				"╰──╯",
		},
		{
			name:           "Only Sections, with empty lines",
			totalHeight:    7,
			totalWidth:     4,
			borderRequired: true,
			lines:          []string{sectionStr, sectionStr, "", sectionStr, sectionStr},
			trucateheight:  false,
			expected: "" +
				"╭──╮\n" +
				"├──┤\n" +
				"├──┤\n" +
				"│  │\n" +
				"├──┤\n" +
				"├──┤\n" +
				"╰──╯",
		},
		{
			name:           "Single line at the end",
			totalHeight:    7,
			totalWidth:     4,
			borderRequired: true,
			lines:          []string{sectionStr, sectionStr, sectionStr, sectionStr, "L1"},
			trucateheight:  false,
			expected: "" +
				"╭──╮\n" +
				"├──┤\n" +
				"├──┤\n" +
				"├──┤\n" +
				"├──┤\n" +
				"│L1│\n" +
				"╰──╯",
		},
		{
			name:           "Only sections",
			totalHeight:    3,
			totalWidth:     4,
			borderRequired: true,
			lines:          []string{sectionStr},
			trucateheight:  false,
			expected: "" +
				"╭──╮\n" +
				"├──┤\n" +
				"╰──╯",
		},
		{
			name:           "Minimal width",
			totalHeight:    4,
			totalWidth:     2,
			borderRequired: true,
			lines:          []string{sectionStr, "L1", sectionStr, sectionStr},
			trucateheight:  false,
			expected: "" +
				"╭╮\n" +
				"├┤\n" +
				"││\n" +
				"╰╯",
		},
		{
			name:           "Minimal height",
			totalHeight:    2,
			totalWidth:     8,
			borderRequired: true,
			lines:          []string{sectionStr, "L1", sectionStr, sectionStr},
			trucateheight:  false,
			expected: "" +
				"╭──────╮\n" +
				"│      │",
			// Border breaks here, because lipgloss creates a 3 line string, and
			// our renderer, than manually adjusts it.
		},
		{
			name:           "Minimal heightBorderless",
			totalHeight:    0,
			totalWidth:     8,
			borderRequired: false,
			lines:          []string{sectionStr, "L1", sectionStr, sectionStr},
			trucateheight:  false,
			expected:       "",
		},
		{
			name:           "No Border",
			totalHeight:    4,
			totalWidth:     4,
			borderRequired: false,
			lines:          []string{sectionStr, "L1", sectionStr},
			trucateheight:  false,
			expected: "" +
				"    \n" +
				"L1  \n" +
				"    \n" +
				"    ",
		},
	}

	for _, tt := range sectionTests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewRenderer(getDefaultTestRendererConfig(
				tt.totalHeight, tt.totalWidth, tt.borderRequired, tt.trucateheight))
			// maxL := r.contentWidth
			// if i >= maxL, check for errors here
			for _, l := range tt.lines {
				if l == sectionStr {
					r.AddSection()
				} else {
					r.AddLines(l)
				}
			}
			assert.Equal(t, tt.expected, r.Render())
		})
	}
}

func TestDynamicHeight(t *testing.T) {
	dynamicHeightTests := []struct {
		name          string
		totalHeight   int
		lines         []string
		trucateheight bool
		expected      string
	}{
		{
			name:          "No truncate",
			totalHeight:   5,
			lines:         []string{"L1"},
			trucateheight: false,
			expected: "" +
				"╭──╮\n" +
				"│L1│\n" +
				"│  │\n" +
				"│  │\n" +
				"╰──╯",
		},
		{
			name:          "Basic truncate",
			totalHeight:   7,
			lines:         []string{"L1", ""},
			trucateheight: true,
			expected: "" +
				"╭──╮\n" +
				"│L1│\n" +
				"│  │\n" +
				"╰──╯",
		},
		{
			name:          "Basic truncate with Sections",
			totalHeight:   100,
			lines:         []string{"L1", "", sectionStr, "L2", "", "L3"},
			trucateheight: true,
			expected: "" +
				"╭──╮\n" +
				"│L1│\n" +
				"│  │\n" +
				"├──┤\n" +
				"│L2│\n" +
				"│  │\n" +
				"│L3│\n" +
				"╰──╯",
		},
	}

	for _, tt := range dynamicHeightTests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewRenderer(getDefaultTestRendererConfig(
				tt.totalHeight, 4, true, tt.trucateheight))
			for _, l := range tt.lines {
				if l == sectionStr {
					r.AddSection()
				} else {
					r.AddLines(l)
				}
			}
			assert.Equal(t, tt.expected, r.Render())
		})
	}
}
func TestBorders(t *testing.T) {
	t.Run("Basic test", func(t *testing.T) {
		r := getDefaultTestRenderer(4, 10, true)
		r.AddLines("L1")
		r.AddLines("L2")
		r.SetBorderTitle("Title")
		res := r.Render()
		expected := "" +
			"╭┤ Titl ├╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰────────╯"
		assert.False(t, r.AreInfoItemsTruncated())
		assert.Equal(t, expected, res, "No margin if title is too big")
		r.SetBorderTitle("T")

		res = r.Render()
		expected = "" +
			"╭─┤ T ├──╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰────────╯"
		assert.Equal(t, expected, res, "Margin should be there if title fits well")

		r.border.SetInfoItems("A", "B")
		assert.False(t, r.AreInfoItemsTruncated())
		res = r.Render()
		expected = "" +
			"╭─┤ T ├──╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰┤A├─┤B├─╯"
		assert.Equal(t, expected, res)

		r.border.SetInfoItems("A1", "B2")
		assert.True(t, r.AreInfoItemsTruncated())
		res = r.Render()
		expected = "" +
			"╭─┤ T ├──╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰┤A├─┤B├─╯"
		assert.Equal(t, expected, res)

		r.border.SetInfoItems("A12345")
		assert.True(t, r.AreInfoItemsTruncated())
		res = r.Render()
		expected = "" +
			"╭─┤ T ├──╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰┤A1234├─╯"
		assert.Equal(t, expected, res, "Info Items Truncation")

		r.SetBorderTitle("✅1✅2✅3")
		r.SetBorderInfoItems()
		res = r.Render()
		expected = "" +
			"╭┤ ✅1 ├─╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰────────╯"
		assert.Equal(t, expected, res, "Double terminal width characters in Title")

		testStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
		title := testStyle.Render("Title")

		r.SetBorderTitle(title)
		res = r.Render()
		expected = "" +
			"╭┤ Titl ├╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰────────╯"

		assert.Equal(t, expected, res, "Ansi escapes are not preserved")

		r.SetBorderTitle("")
		r.SetBorderInfoItems("A", "")

		res = r.Render()
		expected = "" +
			"╭────────╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰─┤A├─┤├─╯"

		assert.Equal(t, expected, res, "Empty title is ignored, but not empty infoitems")

		r.SetBorderInfoItems("AA", "")

		res = r.Render()
		expected = "" +
			"╭────────╮\n" +
			"│L1      │\n" +
			"│L2      │\n" +
			"╰─┤A├─┤├─╯"
		assert.True(t, r.AreInfoItemsTruncated())
		assert.Equal(t, expected, res, "Truncated even if there was enough space because one item was too big")
	})

	t.Run("Different Border", func(t *testing.T) {
		cfg := getDefaultTestRendererConfig(6, 10, true, false)
		cfg.Border = lipgloss.Border{
			Top:    "─",
			Bottom: "*",
			Left:   "+",
			Right:  "│",

			TopLeft:  "╭",
			TopRight: "╮",

			BottomLeft:  "╰",
			BottomRight: "╯",

			MiddleLeft:  "├",
			MiddleRight: "┤",
		}

		r, _ := NewRenderer(cfg)
		r.SetBorderTitle("Title")
		r.SetBorderInfoItems("A")
		r.AddLines("L1")
		r.AddSection()
		r.AddLines("")
		r.AddLines("L2")

		res := r.Render()
		expected := "" +
			"╭┤ Titl ├╮\n" +
			"+L1      │\n" +
			"├────────┤\n" +
			"+        │\n" +
			"+L2      │\n" +
			"╰****┤A├*╯"

		assert.Equal(t, expected, res, "Ansi escape is preserved")
	})
}
