package engine

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/systems"
	"math/rand"
)

type Bootstrap struct {
	World     *ecs.World
	Scheduler *ecs.SchedulerWithContext
	Player    ecs.Entity
	Render    *systems.Render
	MapRender *systems.MapRender
}

func New(r *rand.Rand) Bootstrap {
	w := ecs.NewWorld(r)
	render := &systems.Render{}
	mapRender := &systems.MapRender{}

	// Create player first
	p := w.Create()
	ecs.Add(w, p, components.Player{})
	ecs.Add(w, p, components.Position{})
	ecs.Add(w, p, components.Camera{Width: 200 - 30, Height: 80 - 5})
	ecs.Add(w, p, components.Renderable{Glyph: '@', TileType: components.TileStar})
	ecs.Add(w, p, components.Input{})
	ecs.Add(w, p, components.Velocity{})
	ecs.Add(w, p, components.FuelTank{Current: 100})

	// Create camera system with player as target
	camera := &systems.CameraSystem{Target: p}

	reg := ecs.SystemRegistry{
		UniversalSystems: []ecs.System{systems.InputSystem{}, systems.Tick{}, camera, systems.LevelManager{}, mapRender, render},
		SpaceSystems:     []ecs.System{systems.SpaceMovement{}, systems.FuelSystem{}, systems.PlanetApproachSystem{}, systems.PlanetSelection{}},
		SurfaceSystems:   []ecs.System{systems.SurfaceHeartbeat{}, systems.TerrainGen{}, systems.SurfaceMovement{}, systems.DepthProgression{}, systems.WeatherTick{}, systems.RiverFlow{}, systems.TradeRoutePatrols{}, systems.WildlifeSpawn{}, systems.KingdomGuards{}, systems.QuestSystem{}},
	}
	s := ecs.NewSchedulerWithContext(reg)
	ecs.Add(w, 1, components.WorldInfo{Width: 200, Height: 80})
	ecs.SetWorldContext(w, ecs.WorldContext{CurrentLayer: ecs.LayerSpace, QuestProgress: ecs.QuestProgress{ContractsNeeded: 5}})
	return Bootstrap{World: w, Scheduler: s, Player: p, Render: render, MapRender: mapRender}
}
