# Component & Systems Roadmap

## Persisted Baseline Components (Phase 1)
- Position: world-space location (float64 X,Y).
- Velocity: per-tick delta (VX,VY).
- Camera: viewport origin (X,Y) and size (W,H) centered on target.
- PlayerStats: Fuel, Hull, Drive; extendable for core resources.
- WorldInfo: Tick counter, world dimensions, seed.
- Input: last directional intent (Left/Right/Up/Down).
- Inventory: map[string]int for simple resource counts.
- Resource: Kind, Amount; placed on tiles.
- Tile: static glyph (later biome/material).
- Renderable: entity glyph; later style layers.
- Health: HP/Max; later regen and resistances.

## Biological / Survival Systems (Future)
- Hunger/Satiation: depletes over time; eating adds satiation; penalties at low levels.
- Thirst/Hydration: faster decay than hunger; requires water sources.
- Stamina/Exhaustion: consumed by actions; regen at rest; gates sprint/attacks.
- Temperature/Thermoregulation: ambient + clothing + activity; hypothermia/hyperthermia effects.
- Sleep/Fatigue: circadian cycle; sleep restores stamina/mood; penalties when deprived.
- Mood/MentalState: modifiers from events/environment; affects abilities and UI feedback.
- Buffs/Debuffs: timed modifiers to stats/actions; stack rules.
- Disease/Immunity: contagion, incubation, immunity buildup; cures/medicine.

## Environmental Interaction (Future)
- Weather: world-level state influencing temperature, visibility, movement.
- DayNightCycle: lighting changes, NPC schedules, spawn tables.
- Seasons: long-term modifiers to weather, flora/fauna, resources.
- Lighting: dynamic light sources; shadow casting per tile/entity; affects vision/AI.
- Physics: gravity, friction, momentum; simple integrator for entities.
- Fluid: water tiles with flow; swimming mechanics; buoyancy.
- Fire: spread between tiles; consumes fuel; heat/light emission.
- Sound: propagation with falloff/occlusion; AI hearing and alerts.

## World Systems (Future)
- ChunkLoading/Streaming: partition world into chunks; load/unload by camera.
- ProceduralGeneration: seeded generation for terrain, dungeons, resources.
- Biomes: biome map influencing tiles, spawns, weather.
- Ecosystem: wildlife AI loops, plant growth/decay; resource cycles.
- Decay/Erosion: gradual changes to terrain/structures.
- NPCSchedules: daily routines; tasks and locations; reactions to player/factions.
- Economy/Trading: prices, supply/demand; vendors; loot tables.
- Faction/Reputation: standings changed by actions; unlocks/hostility.

## Crafting / Building (Future)
- Recipe: inputs -> outputs; skills/tools modify yields.
- Blueprint: placeable multi-tile structures; construction stages.
- Durability: item wear and repair; affects performance.
- Quality/Rarity: item tiers with stat ranges and affixes.
- Enchantment/Upgrades: modifiers applied via resources/stations.
- Power/Energy: networks, generators, storage; machine operation costs.
- Automation/Logistics: inserters/belts/pipes; scripted jobs for NPCs.

## Combat / Action (Future)
- Damage/DamageTypes: physical/elemental; resistances; DOTs.
- Armor/Defense: mitigation, block, dodge; location-based hits.
- Skills/Abilities: active/passive abilities; cooldowns and resource costs.
- Cooldowns: per-ability timers; GCD handling.
- Projectile: kinematics; collision with tiles/entities; spread.
- MeleeHitbox: short-range arcs; stagger/knockback.
- Stealth/Detection: visibility/sound vs. AI perception; cover and light levels.

## Notes
- Each system to be modeled as ECS components + systems with clear data flow.
- Determinism: fixed scheduler order; fixed dt; seeded RNG in World.
- Persistence: extend Save/Load to include new components as they ship.
