package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
	"math/rand"
	"strconv"
	"strings"
)

// LayerBackground creates a procedural background using Layer-based rendering
type LayerBackground struct {
	w, h int
	rng  *rand.Rand
}

func NewLayerBackground(w, h int) *LayerBackground {
	return &LayerBackground{w: w, h: h, rng: rand.New(rand.NewSource(42))}
}

func (lb *LayerBackground) GetLayer() rendering.Layer { return rendering.LayerGame }
func (lb *LayerBackground) GetZ() int                 { return rendering.ZBackground }

func (lb *LayerBackground) ToLipglossLayer() *lipgloss.Layer {
	// Create procedural background content
	var content strings.Builder
	
	for y := 0; y < lb.h; y++ {
		if y > 0 {
			content.WriteString("\n")
		}
		for x := 0; x < lb.w; x++ {
			// Create subtle variation in the background
			v := lb.rng.Intn(32) // 0..31 dark grey
			
			// Create hex color string
			hex := strconv.FormatInt(int64(v), 16)
			if len(hex) == 1 {
				hex = "0" + hex
			}
			colorStr := "#" + hex + hex + hex
			
			// Style the character with the color
			styledChar := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorStr)).
				Background(lipgloss.Color(colorStr)).
				Render("█")
			
			content.WriteString(styledChar)
		}
	}
	
	return lipgloss.NewLayer(content.String()).
		X(0).
		Y(0).
		Z(lb.GetZ()).
		ID("background")
}

// For backward compatibility, convert old function
func newBackgroundLayer(w, h int) rendering.LayerContent {
	return NewLayerBackground(w, h)
}