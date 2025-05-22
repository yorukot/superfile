package rendering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentRendererBasic(t *testing.T) {
	t.Run("Basic test", func(t *testing.T) {
		r := NewContentRenderer(6, 5, PlainTruncateRight)
		r.AddLines("123456")
		r.AddLines("12345\n12345", "123")
		assert.Equal(t, 4, r.CntLines())
		r.AddLineWithCustomTruncate("123456", TailsTruncateRight)
		r.AddLines("\t1234")
		// Should be ignored
		r.AddLines("1234")

		res := r.Render()
		expected := "12345\n" +
			"12345\n" +
			"12345\n" +
			"123\n" +
			"12...\n" +
			"    1"
		assert.Equal(t, expected, res, "Basic truncation, and adding lines")

		r.ClearLines()
		assert.Zero(t, r.CntLines(), "ClearLines should remove all content")

		r.AddLines("\x00\x11\x1babc")
		assert.Equal(t, "\x1babc", r.Render())

		r.sanitizeContent = false
		r.ClearLines()

		r.AddLines("\x00\x11\x1babc")

		assert.Equal(t, "\x00\x11\x1babc", r.Render())

		r = NewContentRenderer(0, 0, PlainTruncateRight)
		r.AddLines("L1")
		r.AddLines("L2")
		assert.Equal(t, "", r.Render())
	})
}
