package ui

import (
	"harvester/pkg/components"
	"harvester/pkg/rendering"
)

// CRTShutdownOverlay creates an alpha-masked overlay for CRT shutdown effect
type CRTShutdownOverlay struct {
	width, height int
	progress      float64 // 0.0 to 1.0 animation progress
}

func NewCRTShutdownOverlay(w, h int, progress float64) *CRTShutdownOverlay {
	return &CRTShutdownOverlay{
		width:    w,
		height:   h,
		progress: progress,
	}
}

func (c *CRTShutdownOverlay) GetLayer() rendering.Layer {
	return rendering.LayerTVFrame // Highest layer for global effects
}

func (c *CRTShutdownOverlay) GetZ() int {
	return rendering.ZFrame + 100 // Above TV frame
}

func (c *CRTShutdownOverlay) GetPosition() rendering.Position {
	return rendering.Position{
		Horizontal: rendering.Left,
		Vertical:   rendering.Top,
	}
}

func (c *CRTShutdownOverlay) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: c.width, Height: c.height}
}

func (c *CRTShutdownOverlay) GetGlyphs() [][]rendering.Glyph {
	glyphs := make([][]rendering.Glyph, c.height)
	
	// Calculate visible area (center portion that shrinks)
	totalHeight := c.height
	visibleHeight := int(float64(totalHeight) * (1.0 - c.progress))
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	
	// Calculate margins
	topMargin := (totalHeight - visibleHeight) / 2
	bottomMargin := totalHeight - visibleHeight - topMargin
	
	for y := 0; y < c.height; y++ {
		row := make([]rendering.Glyph, c.width)
		
		// Determine if this row should be masked (hidden)
		var alpha float64 = 0.0 // Default: fully transparent (hidden)
		
		if y >= topMargin && y < totalHeight-bottomMargin {
			// This row is in the visible area
			alpha = 1.0 // Fully visible
		}
		
		// Create masking glyphs
		for x := 0; x < c.width; x++ {
			row[x] = rendering.Glyph{
				Char:       ' ',
				Foreground: rendering.Color{R: 0, G: 0, B: 0},
				Background: rendering.Color{R: 0, G: 0, B: 0},
				Style:      rendering.StyleNone,
				Alpha:      1.0 - alpha, // Invert: 1.0 = hide, 0.0 = show
				BlendMode:  components.BlendNormal,
			}
		}
		
		glyphs[y] = row
	}
	
	return glyphs
}

// CRTOpeningOverlay creates an alpha-masked overlay for CRT opening effect
type CRTOpeningOverlay struct {
	width, height int
	progress      float64 // 0.0 to 1.0 animation progress
}

func NewCRTOpeningOverlay(w, h int, progress float64) *CRTOpeningOverlay {
	return &CRTOpeningOverlay{
		width:    w,
		height:   h,
		progress: progress,
	}
}

func (c *CRTOpeningOverlay) GetLayer() rendering.Layer {
	return rendering.LayerTVFrame // Highest layer for global effects
}

func (c *CRTOpeningOverlay) GetZ() int {
	return rendering.ZFrame + 100 // Above TV frame
}

func (c *CRTOpeningOverlay) GetPosition() rendering.Position {
	return rendering.Position{
		Horizontal: rendering.Left,
		Vertical:   rendering.Top,
	}
}

func (c *CRTOpeningOverlay) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: c.width, Height: c.height}
}

func (c *CRTOpeningOverlay) GetGlyphs() [][]rendering.Glyph {
	glyphs := make([][]rendering.Glyph, c.height)
	
	// Calculate visible area (expands from center)
	totalHeight := c.height
	visibleHeight := int(float64(totalHeight) * c.progress)
	if visibleHeight < 0 {
		visibleHeight = 0
	}
	
	// Calculate margins
	topMargin := (totalHeight - visibleHeight) / 2
	bottomMargin := totalHeight - visibleHeight - topMargin
	
	for y := 0; y < c.height; y++ {
		row := make([]rendering.Glyph, c.width)
		
		// Determine if this row should be masked (hidden)
		var alpha float64 = 0.0 // Default: fully transparent (show content)
		
		if y < topMargin || y >= totalHeight-bottomMargin {
			// This row is outside the visible area - hide it
			alpha = 1.0 // Fully opaque mask
		}
		
		// Create masking glyphs
		for x := 0; x < c.width; x++ {
			row[x] = rendering.Glyph{
				Char:       ' ',
				Foreground: rendering.Color{R: 0, G: 0, B: 0},
				Background: rendering.Color{R: 0, G: 0, B: 0},
				Style:      rendering.StyleNone,
				Alpha:      alpha, // 1.0 = hide, 0.0 = show
				BlendMode:  components.BlendNormal,
			}
		}
		
		glyphs[y] = row
	}
	
	return glyphs
}