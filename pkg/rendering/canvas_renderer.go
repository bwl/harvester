package rendering

import (
	"github.com/charmbracelet/lipgloss/v2"
)

// LinePatch represents a line change in the rendering output
type LinePatch struct {
	Y    int
	Line string
}

// CanvasRenderer uses lipgloss v2 Canvas for pure Layer-based rendering
type CanvasRenderer struct {
	width, height int
	canvas        *lipgloss.Canvas
	layerContent  map[Layer][]LayerContent
}

// ViewRenderer is an alias for CanvasRenderer for backward compatibility
type ViewRenderer = CanvasRenderer

func NewCanvasRenderer(w, h int) *CanvasRenderer {
	return &CanvasRenderer{
		width:        w,
		height:       h,
		canvas:       lipgloss.NewCanvas(),
		layerContent: make(map[Layer][]LayerContent),
	}
}

// NewViewRenderer creates a new ViewRenderer (alias for CanvasRenderer)
func NewViewRenderer(w, h int) *ViewRenderer {
	return NewCanvasRenderer(w, h)
}

func (cr *CanvasRenderer) SetDimensions(w, h int) {
	cr.width, cr.height = w, h
	cr.canvas = lipgloss.NewCanvas()
}

func (cr *CanvasRenderer) GetDimensions() (int, int) {
	return cr.width, cr.height
}

// RegisterContent adds Layer-based content to the renderer
func (cr *CanvasRenderer) RegisterContent(c LayerContent) {
	l := c.GetLayer()
	cr.layerContent[l] = append(cr.layerContent[l], c)
}

func (cr *CanvasRenderer) UnregisterAll() {
	for k := range cr.layerContent {
		cr.layerContent[k] = cr.layerContent[k][:0]
	}
	cr.canvas = lipgloss.NewCanvas()
}

func (cr *CanvasRenderer) Render() string {
	cr.canvas = lipgloss.NewCanvas()

	// Collect all Layer-based content
	var allLayers []*lipgloss.Layer
	for _, slice := range cr.layerContent {
		for _, content := range slice {
			layer := content.ToLipglossLayer()
			allLayers = append(allLayers, layer)
		}
	}

	// Add all layers to canvas
	if len(allLayers) > 0 {
		cr.canvas.AddLayers(allLayers...)
	}

	return cr.canvas.Render()
}

// For compatibility with existing ViewRenderer interface
func (cr *CanvasRenderer) Matrix() *GlyphMatrix {
	// This is deprecated - Layer system doesn't use matrices
	// Return empty matrix for compatibility
	return NewGlyphMatrix(cr.width, cr.height)
}

func (cr *CanvasRenderer) MarkDirty(x, y, w, h int) {
	// Canvas handles dirty regions automatically
}

func (cr *CanvasRenderer) MarkDirtyAll() {
	// Canvas handles dirty regions automatically
}

func (cr *CanvasRenderer) DirtyRegions() []Rect {
	// Canvas handles dirty regions automatically
	return []Rect{{0, 0, cr.width, cr.height}}
}

func (cr *CanvasRenderer) RenderPatch() []LinePatch {
	// For now, return full render as single patch
	// Canvas doesn't expose line-level patches directly
	content := cr.Render()
	lines := make([]LinePatch, 0)

	// Split content into lines
	y := 0
	for _, line := range splitLines(content) {
		lines = append(lines, LinePatch{Y: y, Line: line})
		y++
	}

	return lines
}

func splitLines(content string) []string {
	var lines []string
	current := ""

	for _, r := range content {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}