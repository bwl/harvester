package rendering

import "harvester/pkg/components"

type Layer int

const (
	LayerGame Layer = iota
	LayerUI
	LayerMenu
	LayerTVFrame
)

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

type Color struct {
	R uint8
	G uint8
	B uint8
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

type Glyph struct {
	Char       rune
	Foreground Color
	Background Color
	Style      Style
	Alpha      float64              // 0.0 = transparent, 1.0 = opaque
	BlendMode  components.BlendMode // How this alpha should blend
}

type RenderableContent interface {
	GetLayer() Layer
	GetZ() int
	GetPosition() Position
	GetBounds() Bounds
	GetGlyphs() [][]Glyph
}
