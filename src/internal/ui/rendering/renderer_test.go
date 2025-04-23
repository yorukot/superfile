package rendering

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func getDefaultTestRenderer(totalHeight int, totalWidth int, borderRequired bool) Renderer {
	cfg := DefaultRendererConfig(totalHeight, totalWidth)
	if borderRequired {
		cfg.BorderRequired = true
		cfg.Border = lipgloss.Border{
			Top:          "─",
			Bottom:       "─",
			Left:         "│",
			Right:        "│",
			
			TopLeft:      "╭",
			TopRight:     "╮",
			
			BottomLeft:   "╰",
			BottomRight:  "╯",
			
			MiddleLeft:   "├",
			MiddleRight:  "┤",
		}
	}
	return NewRenderer(cfg)
}

func TestRenderer(t *testing.T) {
	t.Run("Basic test", func(t *testing.T){
		r := getDefaultTestRenderer(4,4, true)
		r.AddLines("L1")
		r.AddLines("L2--Extra line should truncated")
		r.AddLines("L3--Extra line should not be added")
		res := r.Render()
		expected := "╭──╮\n" + 
		            "│L1│\n" + 
		            "│L2│\n" + 
		            "╰──╯"
		assert.Equal(t, expected, res)
	})

	t.Run("Empty Renderer", func(t *testing.T){
		r := getDefaultTestRenderer(0,0,false)
		r.AddLines("L1")
		r.AddLines("L2--Extra line should truncated")
		r.AddLines("L3--Extra line should not be added")
		res := r.Render()
		expected := ""
		assert.Equal(t, expected, res)
	})

	t.Run("Invalid config Renderer", func(t *testing.T){
		r := getDefaultTestRenderer(0,0,true)
		r.AddLines("L1")
		r.AddLines("L2--Extra line should truncated")
		r.AddLines("L3--Extra line should not be added")
		res := r.Render()
		expected := "╭──╮\n" + 
		            "│L1│\n" + 
		            "│L2│\n" + 
		            "╰──╯"
		assert.Equal(t, expected, res)
	})

	t.Run("Section test", func(t *testing.T){
		r := getDefaultTestRenderer(7,4, true)
		r.AddLines("L1")
		r.AddSection()
		r.AddLines("L2")
		r.AddLines("L3 Should be ignored")
		r.AddSection()
		r.AddSection()
		r.AddLines("L4 Should be ignored")
		// Should be ignored
		r.AddSection()
		res := r.Render()
		expected := "╭──╮\n" + 
		            "│L1│\n" + 
		            "├──┤\n" + 
		            "│L2│\n" + 
		            "├──┤\n" + 
		            "├──┤\n" + 
		            "╰──╯"
		assert.Equal(t, expected, res)
	})
}

func TestBorders(t *testing.T) {
	t.Run("Basic test", func(t *testing.T){
		r := getDefaultTestRenderer(4,10, true)
		r.AddLines("L1")
		r.AddLines("L2")
		r.SetBorderTitle("Title")
		res := r.Render()
		expected := "╭┤ Titl ├╮\n" + 
		            "│L1      │\n" + 
		            "│L2      │\n" + 
		            "╰────────╯"
		assert.Equal(t, expected, res, "No margin if title is too big")
		r.SetBorderTitle("T")

		res = r.Render()
		expected =  "╭─┤ T ├──╮\n" + 
		            "│L1      │\n" + 
		            "│L2      │\n" + 
		            "╰────────╯"
		assert.Equal(t, expected, res, "Margin should be there if title fits well")

		r.border.SetInfoItems([]string{"A", "B"})

		res = r.Render()
		expected =  "╭─┤ T ├──╮\n" + 
		            "│L1      │\n" + 
		            "│L2      │\n" + 
		            "╰──┤A├┤B├╯"
		assert.Equal(t, expected, res, "Info Items")
	})
}