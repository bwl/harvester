package components

type Transparency struct {
	Alpha     float64   // 0.0 = fully transparent, 1.0 = fully opaque
	BlendMode BlendMode // How this alpha should blend
}

type BlendMode int

const (
	BlendNormal BlendMode = iota // Standard alpha blending
	BlendAdditive                // Add colors together
	BlendMultiply                // Multiply colors
	BlendScreen                  // Screen blend mode
)

// Helper constructors
func NewTransparency(alpha float64) Transparency {
	return Transparency{
		Alpha:     alpha,
		BlendMode: BlendNormal,
	}
}

func NewTransparencyWithBlend(alpha float64, mode BlendMode) Transparency {
	return Transparency{
		Alpha:     alpha,
		BlendMode: mode,
	}
}