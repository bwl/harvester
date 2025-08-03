package ui

import (
	"github.com/charmbracelet/harmonica"
	"harvester/pkg/components"
	"harvester/pkg/rendering"
	"harvester/pkg/timing"
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
	
	// Use harmonica spring for smooth CRT shutdown effect
	// Spring parameters: stiffness=8.0, damping=0.25 for realistic TV shutdown
	spring := harmonica.NewSpring(timing.HarmonicaFPS, 8.0, 0.25)
	pos := c.progress
	vel := 0.0
	easedProgress, _ := spring.Update(pos, vel, 1.0)
	
	// TV frame padding - don't mask the TV frame itself
	const tvPadding = 3
	
	// Calculate visible area (center portion that shrinks with spring easing)
	// Only affect the inner content area, not the TV frame
	innerHeight := c.height - (2 * tvPadding)
	if innerHeight < 1 {
		innerHeight = 1
	}
	visibleHeight := int(float64(innerHeight) * (1.0 - easedProgress))
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	
	// Calculate margins within the inner content area
	topMargin := tvPadding + (innerHeight-visibleHeight)/2
	bottomMargin := c.height - tvPadding - (innerHeight-visibleHeight)/2
	
	for y := 0; y < c.height; y++ {
		row := make([]rendering.Glyph, c.width)
		
		// Don't mask the TV frame border areas
		if y < tvPadding || y >= c.height-tvPadding {
			// TV frame border - don't add any overlay
			glyphs[y] = row // Empty row (transparent)
			continue
		}
		
		// Determine if this inner content row should be masked (hidden)
		var alpha float64 = 1.0 // Default: hide content
		
		if y >= topMargin && y < bottomMargin {
			// This row is in the visible area
			alpha = 0.0 // Don't hide (transparent overlay)
		}
		
		// Create masking glyphs only for inner content area
		for x := 0; x < c.width; x++ {
			// Don't mask the TV frame side borders
			if x < tvPadding || x >= c.width-tvPadding {
				row[x] = rendering.Glyph{} // Transparent
			} else {
				row[x] = rendering.Glyph{
					Char:       ' ',
					Foreground: rendering.Color{R: 0, G: 0, B: 0},
					Background: rendering.Color{R: 0, G: 0, B: 0},
					Style:      rendering.StyleNone,
					Alpha:      alpha, // 1.0 = hide content, 0.0 = show content
					BlendMode:  components.BlendNormal,
				}
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
	
	// Use harmonica spring for smooth CRT opening effect
	// Different spring parameters: higher stiffness for snappier opening, lower damping for slight overshoot
	spring := harmonica.NewSpring(timing.HarmonicaFPS, 12.0, 0.15)
	pos := 1.0 - c.progress // Start from 1.0 and ease toward 0.0
	vel := 0.0
	maskedProgress, _ := spring.Update(pos, vel, 0.0)
	
	// TV frame padding - don't mask the TV frame itself
	const tvPadding = 3
	
	// Calculate visible area (expands from center with spring easing)
	// Only affect the inner content area, not the TV frame
	innerHeight := c.height - (2 * tvPadding)
	if innerHeight < 1 {
		innerHeight = 1
	}
	visibleHeight := int(float64(innerHeight) * (1.0 - maskedProgress))
	if visibleHeight < 0 {
		visibleHeight = 0
	}
	
	// Calculate margins within the inner content area
	topMargin := tvPadding + (innerHeight-visibleHeight)/2
	bottomMargin := c.height - tvPadding - (innerHeight-visibleHeight)/2
	
	for y := 0; y < c.height; y++ {
		row := make([]rendering.Glyph, c.width)
		
		// Don't mask the TV frame border areas
		if y < tvPadding || y >= c.height-tvPadding {
			// TV frame border - don't add any overlay
			glyphs[y] = row // Empty row (transparent)
			continue
		}
		
		// Determine if this inner content row should be masked (hidden)
		var alpha float64 = 0.0 // Default: show content (transparent overlay)
		
		if y < topMargin || y >= bottomMargin {
			// This row is outside the visible area - hide it
			alpha = 1.0 // Hide content (opaque overlay)
		}
		
		// Create masking glyphs only for inner content area
		for x := 0; x < c.width; x++ {
			// Don't mask the TV frame side borders
			if x < tvPadding || x >= c.width-tvPadding {
				row[x] = rendering.Glyph{} // Transparent
			} else {
				row[x] = rendering.Glyph{
					Char:       ' ',
					Foreground: rendering.Color{R: 0, G: 0, B: 0},
					Background: rendering.Color{R: 0, G: 0, B: 0},
					Style:      rendering.StyleNone,
					Alpha:      alpha, // 1.0 = hide content, 0.0 = show content
					BlendMode:  components.BlendNormal,
				}
			}
		}
		
		glyphs[y] = row
	}
	
	return glyphs
}