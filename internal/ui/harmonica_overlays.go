package ui

import (
	"github.com/charmbracelet/harmonica"
	"harvester/pkg/components"
	"harvester/pkg/rendering"
	"harvester/pkg/timing"
	"math"
)

// FadeOverlay creates a simple alpha fade effect using harmonica
type FadeOverlay struct {
	width, height int
	progress      float64 // 0.0 to 1.0 animation progress
	fadeIn        bool    // true for fade in, false for fade out
}

func NewFadeOverlay(w, h int, progress float64, fadeIn bool) *FadeOverlay {
	return &FadeOverlay{
		width:    w,
		height:   h,
		progress: progress,
		fadeIn:   fadeIn,
	}
}

func (f *FadeOverlay) GetLayer() rendering.Layer { return rendering.LayerTVFrame }
func (f *FadeOverlay) GetZ() int                 { return rendering.ZFrame + 50 }
func (f *FadeOverlay) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (f *FadeOverlay) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: f.width, Height: f.height}
}

func (f *FadeOverlay) GetGlyphs() [][]rendering.Glyph {
	glyphs := make([][]rendering.Glyph, f.height)

	// Use harmonica spring for smooth fade with subtle overshoot
	spring := harmonica.NewSpring(timing.HarmonicaFPS, 6.0, 0.3)
	target := 1.0
	if f.fadeIn {
		target = 0.0 // Fade in means reducing the overlay alpha
	}

	pos := f.progress
	if !f.fadeIn {
		pos = 1.0 - f.progress // Invert for fade out
	}
	vel := 0.0
	easedAlpha, _ := spring.Update(pos, vel, target)

	// Clamp alpha to valid range
	if easedAlpha < 0.0 {
		easedAlpha = 0.0
	} else if easedAlpha > 1.0 {
		easedAlpha = 1.0
	}

	for y := 0; y < f.height; y++ {
		row := make([]rendering.Glyph, f.width)
		for x := 0; x < f.width; x++ {
			row[x] = rendering.Glyph{
				Char:       ' ',
				Foreground: rendering.Color{R: 24, G: 24, B: 28},
				Background: rendering.Color{R: 24, G: 24, B: 28},
				Style:      rendering.StyleNone,
				Alpha:      easedAlpha,
				BlendMode:  components.BlendNormal,
			}
		}
		glyphs[y] = row
	}

	return glyphs
}

// PulseOverlay creates a pulsing alpha effect using harmonica oscillator
type PulseOverlay struct {
	width, height int
	time          float64 // Time in seconds for oscillation
	intensity     float64 // 0.0 to 1.0 pulse intensity
}

func NewPulseOverlay(w, h int, time, intensity float64) *PulseOverlay {
	return &PulseOverlay{
		width:     w,
		height:    h,
		time:      time,
		intensity: intensity,
	}
}

func (p *PulseOverlay) GetLayer() rendering.Layer { return rendering.LayerTVFrame }
func (p *PulseOverlay) GetZ() int                 { return rendering.ZFrame + 25 }
func (p *PulseOverlay) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (p *PulseOverlay) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: p.width, Height: p.height}
}

func (p *PulseOverlay) GetGlyphs() [][]rendering.Glyph {
	glyphs := make([][]rendering.Glyph, p.height)

	// Use simple sine wave for pulsing (harmonica doesn't have oscillators in this version)
	pulseValue := math.Sin(p.time * 2.0 * math.Pi) // 1 Hz pulse

	// Convert sine output (-1 to 1) to alpha (0 to intensity)
	alpha := (pulseValue + 1.0) * 0.5 * p.intensity

	for y := 0; y < p.height; y++ {
		row := make([]rendering.Glyph, p.width)
		for x := 0; x < p.width; x++ {
			row[x] = rendering.Glyph{
				Char:       ' ',
				Foreground: rendering.Color{R: 32, G: 32, B: 64}, // Subtle blue tint
				Background: rendering.Color{R: 16, G: 16, B: 32}, // Dark blue background
				Style:      rendering.StyleNone,
				Alpha:      alpha,
				BlendMode:  components.BlendAdditive, // Additive for glow effect
			}
		}
		glyphs[y] = row
	}

	return glyphs
}

// WaveDistortionOverlay creates a wave-like distortion effect using harmonica
type WaveDistortionOverlay struct {
	width, height int
	time          float64 // Time in seconds for wave animation
	amplitude     float64 // Wave amplitude (0.0 to 1.0)
	frequency     float64 // Wave frequency in Hz
}

func NewWaveDistortionOverlay(w, h int, time, amplitude, frequency float64) *WaveDistortionOverlay {
	return &WaveDistortionOverlay{
		width:     w,
		height:    h,
		time:      time,
		amplitude: amplitude,
		frequency: frequency,
	}
}

func (w *WaveDistortionOverlay) GetLayer() rendering.Layer { return rendering.LayerTVFrame }
func (w *WaveDistortionOverlay) GetZ() int                 { return rendering.ZFrame + 75 }
func (w *WaveDistortionOverlay) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (w *WaveDistortionOverlay) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: w.width, Height: w.height}
}

func (w *WaveDistortionOverlay) GetGlyphs() [][]rendering.Glyph {
	glyphs := make([][]rendering.Glyph, w.height)

	for y := 0; y < w.height; y++ {
		row := make([]rendering.Glyph, w.width)

		// Calculate wave phase based on row position
		rowPhase := float64(y) / float64(w.height) * 2.0 // 0 to 2 across height
		waveValue := math.Sin((w.time + rowPhase) * w.frequency * 2.0 * math.Pi)

		// Convert wave to alpha modulation
		alphaModulation := (waveValue + 1.0) * 0.5 * w.amplitude

		for x := 0; x < w.width; x++ {
			// Create subtle distortion effect
			row[x] = rendering.Glyph{
				Char:       ' ',
				Foreground: rendering.Color{R: 24, G: 24, B: 28},
				Background: rendering.Color{R: 8, G: 8, B: 12}, // Very subtle blue
				Style:      rendering.StyleNone,
				Alpha:      alphaModulation,
				BlendMode:  components.BlendMultiply, // Multiply for subtle darkening
			}
		}
		glyphs[y] = row
	}

	return glyphs
}

// SpringBounceOverlay demonstrates harmonica spring physics with different parameters
type SpringBounceOverlay struct {
	width, height int
	progress      float64 // 0.0 to 1.0 animation progress
	bounceHeight  int     // Height of the bounce effect in characters
}

func NewSpringBounceOverlay(w, h int, progress float64, bounceHeight int) *SpringBounceOverlay {
	return &SpringBounceOverlay{
		width:        w,
		height:       h,
		progress:     progress,
		bounceHeight: bounceHeight,
	}
}

func (s *SpringBounceOverlay) GetLayer() rendering.Layer { return rendering.LayerTVFrame }
func (s *SpringBounceOverlay) GetZ() int                 { return rendering.ZFrame + 60 }
func (s *SpringBounceOverlay) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (s *SpringBounceOverlay) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: s.width, Height: s.height}
}

func (s *SpringBounceOverlay) GetGlyphs() [][]rendering.Glyph {
	glyphs := make([][]rendering.Glyph, s.height)

	// Use harmonica spring with high stiffness and low damping for bouncy effect
	spring := harmonica.NewSpring(timing.HarmonicaFPS, 15.0, 0.1) // Very bouncy!
	pos := s.progress
	vel := 0.0
	easedProgress, _ := spring.Update(pos, vel, 1.0)

	// Calculate bounce offset (can overshoot beyond target)
	bounceOffset := int(easedProgress * float64(s.bounceHeight))
	if bounceOffset > s.height {
		bounceOffset = s.height
	}

	for y := 0; y < s.height; y++ {
		row := make([]rendering.Glyph, s.width)

		// Create a moving line that bounces
		var alpha float64 = 0.0
		distanceFromBounce := abs(y - bounceOffset)
		if distanceFromBounce < 3 { // 3-line wide bounce effect
			alpha = 1.0 - float64(distanceFromBounce)/3.0
		}

		for x := 0; x < s.width; x++ {
			row[x] = rendering.Glyph{
				Char:       '▬',                                   // Block character for line effect
				Foreground: rendering.Color{R: 255, G: 128, B: 0}, // Orange
				Background: rendering.Color{R: 64, G: 32, B: 0},   // Dark orange
				Style:      rendering.StyleBold,
				Alpha:      alpha,
				BlendMode:  components.BlendNormal,
			}
		}
		glyphs[y] = row
	}

	return glyphs
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
