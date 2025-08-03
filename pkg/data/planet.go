package data

type BiomeType int

const (
	BiomeToftForest BiomeType = iota
	BiomeVolcanic
	BiomeIce
)

type Planet struct {
	ID              int
	Name            string
	Biome           BiomeType
	MaxDepth        int
	GravityModifier float64
	Seed            int64
}

type DepthLayer struct {
	MinDepth       int
	MaxDepth       int
	DifficultyMod  float64
	UniqueFeatures []string
}

type PlanetGenerator struct {
	Seed     int64
	Biome    BiomeType
	MaxDepth int
}

func (pg *PlanetGenerator) GenerateToft() *Planet {
	return &Planet{ID: 1, Name: "Toft", Biome: BiomeToftForest, MaxDepth: 120, GravityModifier: 1.0, Seed: pg.Seed}
}

func ToftDepthLayers() []DepthLayer {
	return []DepthLayer{
		{MinDepth: 0, MaxDepth: 10, DifficultyMod: 1.0, UniqueFeatures: []string{"villages", "forest_edge"}},
		{MinDepth: 11, MaxDepth: 60, DifficultyMod: 2.0, UniqueFeatures: []string{"rivers", "trade_routes", "highlands"}},
		{MinDepth: 61, MaxDepth: 120, DifficultyMod: 3.5, UniqueFeatures: []string{"deep_wilds", "kingdom_borders", "mountain_passes"}},
	}
}
