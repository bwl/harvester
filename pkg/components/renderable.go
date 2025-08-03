package components

import "github.com/charmbracelet/lipgloss"

type TileType int

const (
	TileUnknown TileType = iota
	TileGalaxy
	TileStar
	TilePlanet
	TileForest
	TileMountain
	TileRiver
	TileLava
	TileNebula
	TileGalaxyCore
	TileAsteroid
	TileComet
)

type SpecialEffect int

const (
	EffectNone SpecialEffect = iota
	EffectPulsing
	EffectTwinkling
	EffectBurning
	EffectFrozen
)

type ColorModifier struct {
	TintColor      *lipgloss.Color
	PulseRate      float64
	TemperatureHue float64
	Special        SpecialEffect
}

type Renderable struct {
	Glyph    rune
	TileType TileType
	StyleMod *ColorModifier
}
