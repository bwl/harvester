# Status & Next Steps

## Current Status
- Context-aware layers implemented: Space, Surface, Deep enums; WorldContext stored in ECS; contextual scheduler wired.
- UI shows current layer in HUD; logs transition events.
- Space systems stubbed and wired: input→velocity, movement, fuel burn, planet approach (switch to Surface on X>50), planet selection scaffold.
- Surface systems wired: heartbeat tick, Toft biome stubs (WeatherTick, RiverFlow, TradeRoutePatrols, WildlifeSpawn, KingdomGuards).
- Data models added: Planet (Toft), Biome enum, DepthLayer presets.
- Build, fmt, vet pass locally.

## Immediate Next Steps
1. Depth progression
   - Track ctx.Depth and increment/decrement on descent/ascend actions.
   - Apply DepthLayer difficulty mods to spawn rates and hazards.
   - Expose depth in HUD and logs.
2. Quest system (Royal Charter)
   - Define QuestGate, Quest, Objective structs.
   - Implement basic tracker and persistence; HUD/Modal to display progress.
   - Hook simple triggers (trade contracts collected) to mark progress; unlock escape.
3. Surface terrain generation (Toft)
   - Generate forest/rivers/mountains tiles deterministically from Seed and DepthLayer.
   - Spawn POIs: villages, trade routes, bridges/fords; ensure camera bounds.
   - Render tiles via MapRender with clear glyphs.
4. Planet selection UI (Space)
   - Render three planet cards with name/biome/depth; basic input to choose.
   - On select, set WorldContext PlanetID and seed; initiate approach/landing.
5. Biome systems MVP behaviors
   - WeatherTick: movement penalty when raining; visual hint in HUD.
   - RiverFlow: mark river tiles and crossing cost; spawn bridges/fords.
   - TradeRoutePatrols: periodic guard spawn along routes (non-hostile unless provoked).
   - WildlifeSpawn: spawn tables by region/depth; simple hostile/neutral logic.
   - Faction/Guards: placeholder rep on events; guards react to aggression.
   - Rain affects trade routes: per-tile rain flags increase route cost/slow patrols; reroute logic; render wet routes.
6. Save/Load extensions
   - Persist WorldContext, Planet state, Depth, quest progress, faction rep.

## Nice-to-Haves (after MVP)
- Layer transition animations; landing/descent/escape sequences.
- Performance: viewport culling and dirty-region updates.
- VHS demo scripts for core loop: space→landing→surface roam→quest progress.

## Testing Plan
- Deterministic seed for Toft gen; snapshot tests for terrain at fixed seeds.
- System integration tests: layer switching, quest progression flags.
- Perf sanity checks with large terrain and spawn counts.
