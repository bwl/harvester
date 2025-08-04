package rendering

import (
	"github.com/charmbracelet/lipgloss/v2"
	"image/color"
)

// LayerContent represents content that can be rendered as a lipgloss Layer
type LayerContent interface {
	GetLayer() Layer                  // Which rendering layer (Game, UI, Menu, etc.)
	GetZ() int                        // Z-index for depth sorting
	ToLipglossLayer() *lipgloss.Layer // Convert to lipgloss Layer
}

// StyledContent represents content with lipgloss styling
type StyledContent struct {
	layer   Layer
	z       int
	content string
	style   lipgloss.Style
	x, y    int
	id      string
}

// NewStyledContent creates a new styled content layer
func NewStyledContent(layer Layer, z int, content string) *StyledContent {
	return &StyledContent{
		layer:   layer,
		z:       z,
		content: content,
		style:   lipgloss.NewStyle(),
	}
}

func (sc *StyledContent) GetLayer() Layer { return sc.layer }
func (sc *StyledContent) GetZ() int       { return sc.z }

func (sc *StyledContent) ToLipglossLayer() *lipgloss.Layer {
	styledContent := sc.style.Render(sc.content)
	layer := lipgloss.NewLayer(styledContent).
		X(sc.x).
		Y(sc.y).
		Z(sc.z)

	if sc.id != "" {
		layer.ID(sc.id)
	}

	return layer
}

// Fluent interface methods
func (sc *StyledContent) WithStyle(style lipgloss.Style) *StyledContent {
	sc.style = style
	return sc
}

func (sc *StyledContent) WithPosition(x, y int) *StyledContent {
	sc.x, sc.y = x, y
	return sc
}

func (sc *StyledContent) WithID(id string) *StyledContent {
	sc.id = id
	return sc
}

func (sc *StyledContent) WithForeground(c color.Color) *StyledContent {
	sc.style = sc.style.Foreground(c)
	return sc
}

func (sc *StyledContent) WithBackground(c color.Color) *StyledContent {
	sc.style = sc.style.Background(c)
	return sc
}

func (sc *StyledContent) WithBorder(border lipgloss.Border) *StyledContent {
	sc.style = sc.style.Border(border)
	return sc
}

func (sc *StyledContent) WithBorderColor(c color.Color) *StyledContent {
	sc.style = sc.style.BorderForeground(c)
	return sc
}

func (sc *StyledContent) WithPadding(top, right, bottom, left int) *StyledContent {
	sc.style = sc.style.Padding(top, right, bottom, left)
	return sc
}

func (sc *StyledContent) WithMargin(top, right, bottom, left int) *StyledContent {
	sc.style = sc.style.Margin(top, right, bottom, left)
	return sc
}

func (sc *StyledContent) WithWidth(width int) *StyledContent {
	sc.style = sc.style.Width(width)
	return sc
}

func (sc *StyledContent) WithHeight(height int) *StyledContent {
	sc.style = sc.style.Height(height)
	return sc
}

func (sc *StyledContent) WithAlign(align lipgloss.Position) *StyledContent {
	sc.style = sc.style.Align(align)
	return sc
}

// CompositeContent represents multiple layers composed together
type CompositeContent struct {
	layer    Layer
	z        int
	children []LayerContent
}

func NewCompositeContent(layer Layer, z int) *CompositeContent {
	return &CompositeContent{
		layer:    layer,
		z:        z,
		children: make([]LayerContent, 0),
	}
}

func (cc *CompositeContent) GetLayer() Layer { return cc.layer }
func (cc *CompositeContent) GetZ() int       { return cc.z }

func (cc *CompositeContent) AddChild(child LayerContent) *CompositeContent {
	cc.children = append(cc.children, child)
	return cc
}

func (cc *CompositeContent) ToLipglossLayer() *lipgloss.Layer {
	// Create a canvas for compositing child layers
	canvas := lipgloss.NewCanvas()

	// Add all children to the canvas
	for _, child := range cc.children {
		childLayer := child.ToLipglossLayer()
		canvas.AddLayers(childLayer)
	}

	// Render the canvas to a string and create a layer from it
	compositeContent := canvas.Render()
	return lipgloss.NewLayer(compositeContent).Z(cc.z)
}