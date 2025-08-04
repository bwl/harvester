# Game Vision: Multi-Planet Deep Exploration

## Core Game Vision

Harvest of Stars is a **multi-layered space exploration roguelike** that combines rocket ship navigation, planetary landing, and deep underground exploration. Players begin in space with a rocket ship and must strategically choose which planets to explore, knowing that once they land, they're committed to that world until they complete the necessary quests to re-enable space travel.

## Game Loop Structure

### **1. Space Layer - Navigation & Choice**
- Player spawns in a rocket ship with limited fuel
- **3 planets** available for exploration, each with unique characteristics
- **Fuel management** creates tension - players can keep flying but risk running out
- **Strategic decision**: Which planet offers the best risk/reward ratio?
- Planets are visible with basic information (biome type, difficulty hints)

### **2. Planet Landing - Commitment Point**
- Landing is **irreversible** - player is now committed to this planet
- Rocket ship becomes inoperable until quest conditions are met
- Planet surface reveals the chosen biome and initial exploration area
- **"You're stuck here now"** - creates urgency and investment

### **3. Deep Exploration - The Real Game**
- Each planet can be **hundreds of levels deep**
- **Angband-style depth progression** with increasing difficulty and rewards
- **Biome-specific systems** create unique gameplay on each planet
- Quest objectives scattered throughout the depths
- Must find items, complete objectives, and solve planet-specific challenges

### **4. Quest Gates - Escape Condition**
- Each planet has **unique escape requirements**
- Examples: Craft heat shield, repair ancient technology, tame local wildlife
- Completing requirements **re-enables space travel**
- Creates structured progression within open exploration

### **5. Multi-Planet Progression**
- Knowledge and items from previous planets may help on new ones
- **Cross-planet synergies** reward exploration of multiple worlds
- Each planet offers different technological/magical advancement paths

## Planet & Biome Design

### **Planet Archetypes**

#### **Vulcanus - Volcanic World**
- **Surface**: Lava flows, geysers, volcanic glass formations
- **Deep**: Magma chambers, rare heat-resistant crystals, molten core access
- **Systems**: Heat damage, lava flow dynamics, thermal vents
- **Escape Quest**: Forge heat shield from volcanic materials
- **Unique Mechanics**: Temperature management, lava surfing, crystal formation

#### **Glacialis - Ice World**
- **Surface**: Frozen tundra, ice caves, aurora phenomena
- **Deep**: Permafrost layers, ancient frozen specimens, ice crystal cores
- **Systems**: Hypothermia, ice physics, frozen water mechanics
- **Escape Quest**: Revive ancient ice-locked technology
- **Unique Mechanics**: Ice formation/melting, preservation systems, cryogenic puzzles

#### **Mechanicus - Ancient Ruins World**
- **Surface**: Abandoned structures, broken technology, archaeological sites
- **Deep**: Ancient data cores, robotic guardians, technological mysteries
- **Systems**: Digital archaeology, ancient AI interactions, tech restoration
- **Escape Quest**: Rebuild planetary defense system to enable safe departure
- **Unique Mechanics**: Technology puzzles, AI negotiations, ancient programming

### **Biome System Features**

#### **Environmental Hazards**
- **Vulcanus**: Lava eruptions, toxic gas vents, extreme heat
- **Glacialis**: Blizzards, ice collapses, freezing winds
- **Mechanicus**: Radiation, malfunctioning defenses, data corruption

#### **Resource Types**
- **Vulcanus**: Volcanic glass, heat crystals, molten metals
- **Glacialis**: Ice cores, preserved organics, crystalline water
- **Mechanicus**: Ancient circuits, data fragments, exotic alloys

#### **Depth Progression**
Each planet follows depth-based difficulty scaling:
- **Surface (0-10 levels)**: Tutorial area, basic resources, introduction to biome
- **Shallow (11-50 levels)**: Core biome mechanics, intermediate challenges
- **Deep (51-200+ levels)**: Rare resources, extreme hazards, quest objectives

## ECS Architecture Integration

### **Multi-Layer World Management**

```go
type GameLayer int
const (
    LayerSpace GameLayer = iota
    LayerPlanetSurface
    LayerPlanetDeep
)

type WorldContext struct {
    CurrentLayer GameLayer
    PlanetID     int
    Depth        int
    BiomeType    BiomeType
}

type LayerTransition struct {
    TargetLayer GameLayer
    TargetDepth int
    RequiredItems []string
    QuestComplete bool
}
```

### **Context-Aware System Activation**

```go
type SystemRegistry struct {
    SpaceSystems     []System  // Fuel management, navigation, planet approach
    SurfaceSystems   []System  // Biome weather, surface exploration
    DeepSystems      []System  // Mining, depth pressure, cave-ins
    UniversalSystems []System  // Movement, rendering, save/load
}

func (s *Scheduler) UpdateForContext(ctx WorldContext, world *World) {
    activeSystems := s.getSystemsForLayer(ctx.CurrentLayer)
    for _, system := range activeSystems {
        system.Update(dt, world)
    }
}
```

### **Planet Generation Components**

```go
type Planet struct {
    ID              int
    Name            string
    BiomeType       BiomeType
    MaxDepth        int
    GravityModifier float64
    AtmosphereType  string
    QuestSeed       int64
    EscapeCondition QuestGate
}

type BiomeType int
const (
    BiomeVolcanic BiomeType = iota
    BiomeIce
    BiomeAncientRuins
    BiomeDesert
    BiomeOcean
    BiomeCrystalline
    BiomeGaseous
)

type BiomeData struct {
    Temperature     float64
    Pressure        float64
    ResourceTypes   []ResourceType
    HazardTypes     []HazardType
    UniqueRules     []string  // "lava_flows", "zero_gravity", "ancient_tech"
    SystemOverrides []string  // Systems to enable/disable for this biome
}
```

### **Depth-Based Progression**

```go
type DepthLayer struct {
    MinDepth        int
    MaxDepth        int
    DifficultyMod   float64
    ResourceRarity  map[ResourceType]float64
    UniqueFeatures  []string
    HazardIntensity float64
}

// Example: Vulcanus depth progression
vulcanusLayers := []DepthLayer{
    {0, 10, 1.0, lowRarity, []string{"surface_geysers"}, 1.0},
    {11, 50, 2.0, medRarity, []string{"lava_tubes"}, 2.0},
    {51, 200, 4.0, highRarity, []string{"molten_core", "rare_crystals"}, 4.0},
}
```

## System Design Examples

### **Space Layer Systems**

#### **Fuel Management System**
```go
type FuelSystem struct{}
func (s FuelSystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.CurrentLayer != LayerSpace { return }
    
    for entity, fuel, velocity := range Query2[FuelTank, Velocity](w) {
        burnRate := calculateFuelBurn(velocity)
        fuel.Current -= burnRate * dt
        
        if fuel.Current <= 0 {
            // Emergency landing or game over
            triggerFuelCrisis(w, entity)
        }
    }
}
```

#### **Planet Approach System**
```go
type PlanetApproachSystem struct{}
func (s PlanetApproachSystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.CurrentLayer != LayerSpace { return }
    
    for entity, pos, ship := range Query2[Position, Spaceship](w) {
        nearbyPlanet := detectNearbyPlanet(pos)
        if nearbyPlanet != nil && ship.LandingSequence {
            initiatePlanetLanding(w, entity, nearbyPlanet)
        }
    }
}
```

### **Planet-Specific Systems**

#### **Volcanic Biome Systems**
```go
type VolcanicWeatherSystem struct{}
func (s VolcanicWeatherSystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.BiomeType != BiomeVolcanic { return }
    
    // Handle lava eruptions
    if shouldErupt(w.RNG) {
        pos := selectEruptionSite(w)
        createLavaFlow(w, pos)
    }
    
    // Apply heat damage to entities
    for entity, pos, health := range Query2[Position, Health](w) {
        heatLevel := calculateHeatAtPosition(pos)
        if heatLevel > SAFE_TEMPERATURE {
            applyHeatDamage(health, heatLevel)
        }
    }
}

type LavaFlowSystem struct{}
func (s LavaFlowSystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.BiomeType != BiomeVolcanic { return }
    
    for entity, flow := range Query1[LavaFlow](w) {
        flow.update(dt)
        affectedTiles := flow.getAffectedTiles()
        
        for _, tile := range affectedTiles {
            convertToLavaTile(w, tile)
            damageEntitiesOnTile(w, tile)
        }
    }
}
```

#### **Ice Biome Systems**
```go
type FreezingSystem struct{}
func (s FreezingSystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.BiomeType != BiomeIce { return }
    
    for entity, pos, temp := range Query2[Position, BodyTemperature](w) {
        ambientTemp := getAmbientTemperature(pos, ctx.Depth)
        temp.Current += (ambientTemp - temp.Current) * THERMAL_TRANSFER_RATE * dt
        
        if temp.Current < HYPOTHERMIA_THRESHOLD {
            addComponent(w, entity, Hypothermia{Severity: calculateSeverity(temp.Current)})
        }
    }
}

type IcePhysicsSystem struct{}
func (s IcePhysicsSystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.BiomeType != BiomeIce { return }
    
    // Handle ice formation and melting
    for entity, pos, water := range Query2[Position, WaterTile](w) {
        if shouldFreeze(pos, ctx) {
            convertToIce(w, entity)
        }
    }
}
```

#### **Ancient Ruins Systems**
```go
type DigitalArchaeologySystem struct{}
func (s DigitalArchaeologySystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.BiomeType != BiomeAncientRuins { return }
    
    for entity, pos, tool := range Query2[Position, ArchaeologyTool](w) {
        nearbyArtifacts := findAncientArtifacts(pos)
        
        for _, artifact := range nearbyArtifacts {
            if canAnalyze(tool, artifact) {
                dataFragment := extractData(artifact, tool)
                addToKnowledgeBase(w, entity, dataFragment)
            }
        }
    }
}

type AncientAISystem struct{}
func (s AncientAISystem) Update(dt float64, w *World) {
    ctx := GetWorldContext(w)
    if ctx.BiomeType != BiomeAncientRuins { return }
    
    for entity, ai, interaction := range Query2[AncientAI, PlayerInteraction](w) {
        if ai.isAwakened() {
            response := ai.processPlayerInput(interaction.Input)
            createDialogEvent(w, entity, response)
            
            if ai.trustLevel > COOPERATION_THRESHOLD {
                revealSecretKnowledge(w, entity, ai.secretData)
            }
        }
    }
}
```

## Quest System Integration

### **Quest Gate Mechanics**
```go
type QuestGate struct {
    PlanetID         int
    RequiredItems    []ItemType
    RequiredKnowledge []string
    CompletionFlag   string
    UnlocksTravel    bool
}

type PlanetQuest struct {
    ID           string
    PlanetID     int
    Description  string
    Objectives   []QuestObjective
    Rewards      []Reward
    EscapeReward bool  // Enables space travel when completed
}

type QuestObjective struct {
    Type        ObjectiveType  // "collect", "craft", "discover", "defeat"
    Target      string
    Count       int
    Location    string
    Completed   bool
}
```

### **Planet-Specific Quest Examples**

#### **Vulcanus Escape Quests**
```go
vulcanusQuests := []PlanetQuest{
    {
        ID: "forge_heat_shield",
        Description: "Craft a heat shield capable of withstanding re-entry",
        Objectives: []QuestObjective{
            {Type: "collect", Target: "volcanic_glass", Count: 15},
            {Type: "collect", Target: "heat_crystals", Count: 5},
            {Type: "craft", Target: "heat_shield", Count: 1},
        },
        EscapeReward: true,
    },
    {
        ID: "tame_lava_wyrm",
        Description: "Befriend a lava wyrm to gain thermal immunity",
        Objectives: []QuestObjective{
            {Type: "discover", Target: "lava_wyrm_nest", Location: "depth_100+"},
            {Type: "collect", Target: "fire_flower", Count: 10},
            {Type: "complete_ritual", Target: "thermal_bonding"},
        },
        EscapeReward: true,
    },
}
```

#### **Glacialis Escape Quests**
```go
glacialisQuests := []PlanetQuest{
    {
        ID: "revive_cryo_core",
        Description: "Restore the planet's ancient cryogenic systems",
        Objectives: []QuestObjective{
            {Type: "discover", Target: "cryo_core", Location: "depth_150+"},
            {Type: "collect", Target: "ice_essence", Count: 20},
            {Type: "repair", Target: "cryo_core_systems"},
        },
        EscapeReward: true,
    },
}
```

#### **Mechanicus Escape Quests**
```go
mechanicusQuests := []PlanetQuest{
    {
        ID: "rebuild_defense_grid",
        Description: "Restore planetary defenses to enable safe departure",
        Objectives: []QuestObjective{
            {Type: "collect", Target: "ancient_circuits", Count: 25},
            {Type: "discover", Target: "central_core", Location: "depth_200+"},
            {Type: "program", Target: "defense_protocols"},
            {Type: "defeat", Target: "corrupted_guardian"},
        },
        EscapeReward: true,
    },
}
```

## Procedural Generation Strategy

### **Planet Generation**
```go
type PlanetGenerator struct {
    Seed        int64
    BiomeType   BiomeType
    MaxDepth    int
    Difficulty  float64
}

func (pg *PlanetGenerator) GeneratePlanet() *Planet {
    rng := rand.New(rand.NewSource(pg.Seed))
    
    planet := &Planet{
        BiomeType: pg.BiomeType,
        MaxDepth:  pg.MaxDepth,
    }
    
    // Generate biome-specific features
    switch pg.BiomeType {
    case BiomeVolcanic:
        planet.addVolcanicFeatures(rng)
    case BiomeIce:
        planet.addIceFeatures(rng)
    case BiomeAncientRuins:
        planet.addAncientFeatures(rng)
    }
    
    return planet
}
```

### **Depth-Based Generation**
```go
type DepthGenerator struct {
    BiomeType    BiomeType
    CurrentDepth int
    PlanetSeed   int64
}

func (dg *DepthGenerator) GenerateLevel() *Level {
    // Combine planet seed with depth for deterministic but varied levels
    levelSeed := combineSeed(dg.PlanetSeed, dg.CurrentDepth)
    rng := rand.New(rand.NewSource(levelSeed))
    
    level := &Level{
        Depth: dg.CurrentDepth,
        BiomeType: dg.BiomeType,
    }
    
    // Apply depth-based difficulty scaling
    difficultyMod := calculateDifficultyMod(dg.CurrentDepth)
    level.applyDifficulty(difficultyMod)
    
    // Generate biome-specific level content
    dg.generateBiomeContent(level, rng)
    
    return level
}
```

## Implementation Roadmap

### **Phase 1: Core Layer System**
- Implement layer switching (Space ↔ Planet Surface ↔ Planet Deep)
- Basic rocket ship navigation
- Simple planet landing mechanics
- Context-aware system activation

### **Phase 2: Single Planet Implementation**
- Implement one complete planet (Vulcanus recommended)
- Volcanic biome systems (lava, heat, eruptions)
- Basic depth progression (surface to deep)
- Simple escape quest (heat shield crafting)

### **Phase 3: Multi-Planet System**
- Add remaining two planets (Glacialis, Mechanicus)
- Planet selection interface in space
- Biome-specific systems for each planet
- Unique escape quests for each world

### **Phase 4: Deep Progression**
- Extend depth to 200+ levels per planet
- Implement depth-based difficulty scaling
- Add rare resources and powerful items at extreme depths
- Complex quest chains spanning multiple depth layers

### **Phase 5: Cross-Planet Integration**
- Knowledge and items that carry between planets
- Meta-progression systems
- Advanced quest chains requiring multiple planet exploration
- End-game content that utilizes all planetary knowledge

## Technical Considerations

### **Save System Extensions**
- Multiple planet states must be preserved
- Player progress across different worlds
- Quest completion status per planet
- Cross-planet inventory and knowledge management

### **Performance Optimization**
- Only load active planet/layer data
- Unload distant planet data to memory
- Efficient depth-based streaming for deep exploration
- Optimize biome-specific system activation

### **Testing Strategy**
- Deterministic planet generation testing
- Biome system integration tests
- Quest completion validation
- Cross-layer transition testing
- Performance testing for deep exploration

## Vision Summary

This multi-planet deep exploration design transforms Harvest of Stars from a simple expanding universe game into a rich, **strategic space exploration experience**. The combination of:

- **High-stakes planet selection** (fuel management + commitment)
- **Biome-specific gameplay systems** (unique mechanics per world)
- **Deep vertical progression** (hundreds of levels of increasing challenge)
- **Quest-gated progression** (must earn your escape)
- **Cross-planet meta-progression** (knowledge and items carry forward)

Creates a game with enormous replay value where each planet offers a completely different gameplay experience, but the underlying ECS architecture remains consistent and expandable.

The vision leverages your excellent ECS foundation while providing clear direction for years of development, with each planet essentially serving as a full game's worth of content within the larger space exploration framework.