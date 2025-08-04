package rendering

// Compatibility stubs for components that haven't been converted yet

// Basic types for compatibility
type Color struct {
	R, G, B uint8
}

type Style uint32

const (
	StyleNone Style = 0
	StyleBold Style = 1 << iota
	StyleItalic
	StyleUnderline
	StyleDim
	StyleReverse
)

type Position struct {
	Horizontal HorizontalAlign
	Vertical   VerticalAlign
	OffsetX    int
	OffsetY    int
}

type Bounds struct {
	Width  int
	Height int
}

type HorizontalAlign int
type VerticalAlign int

const (
	Left HorizontalAlign = iota
	CenterH
	Right
)

const (
	Top VerticalAlign = iota
	CenterV
	Bottom
)

type Glyph struct {
	Char       rune
	Foreground Color
	Background Color
	Style      Style
	Alpha      float64
	BlendMode  int // For compatibility with components.BlendMode
}

// Legacy RenderableContent interface for components not yet converted
type RenderableContent interface {
	GetLayer() Layer
	GetZ() int
	GetPosition() Position
	GetBounds() Bounds
	GetGlyphs() [][]Glyph
}

// Function to calculate position (simplified version)
func CalculatePosition(pos Position, bounds Bounds, containerW, containerH int) (int, int) {
	x, y := 0, 0
	
	switch pos.Horizontal {
	case CenterH:
		x = (containerW - bounds.Width) / 2
	case Right:
		x = containerW - bounds.Width
	}
	x += pos.OffsetX
	
	switch pos.Vertical {
	case CenterV:
		y = (containerH - bounds.Height) / 2
	case Bottom:
		y = containerH - bounds.Height
	}
	y += pos.OffsetY
	
	return x, y
}