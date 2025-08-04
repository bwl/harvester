package testharness

import (
	"encoding/json"

	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/systems"
)

type Controller struct {
	World  *ecs.World
	Sched  *ecs.Scheduler
	Player ecs.Entity
	Render *systems.Render
}

type Options struct {
	Seed   int64
	Width  int
	Height int
}

func NewController(opt Options) *Controller {
	w := ecs.NewWorld(nil)
	r := &systems.Render{}
	cam := &systems.CameraSystem{}
	s := ecs.NewScheduler(systems.InputSystem{}, systems.Movement{}, cam, systems.Harvest{}, systems.Combat{}, systems.Tick{}, r)
	c := &Controller{World: w, Sched: s, Render: r}
	// player
	p := w.Create()
	ecs.Add(w, p, components.Position{})
	ecs.Add(w, p, components.Renderable{Glyph: '@'})
	ecs.Add(w, p, components.Input{})
	ecs.Add(w, p, components.Velocity{})
	ecs.Add(w, p, components.PlayerStats{Fuel: 100, Hull: 100, Drive: 1})
	ecs.Add(w, p, components.Camera{Width: opt.Width, Height: opt.Height})
	cam.Target = p
	c.Player = p
	// world info
	if opt.Width == 0 {
		opt.Width = 200
	}
	if opt.Height == 0 {
		opt.Height = 80
	}
	ecs.Add(w, 1, components.WorldInfo{Width: opt.Width, Height: opt.Height})
	// simple sparse starfield
	for y := 0; y < opt.Height; y++ {
		for x := 0; x < opt.Width; x++ {
			if (x+y)%11 == 0 {
				e := w.Create()
				ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
				ecs.Add(w, e, components.Tile{Glyph: '*', Type: components.TileStar})
			}
		}
	}
	return c
}

func (c *Controller) InjectKey(key string) { systems.SetPlayerInput(c.World, c.Player, key) }

func (c *Controller) Tick(n int, dt float64) {
	for i := 0; i < n; i++ {
		c.Sched.Update(dt, c.World)
	}
}

type Snapshot struct {
	Player struct {
		X     int `json:"x"`
		Y     int `json:"y"`
		Fuel  int `json:"fuel"`
		Hull  int `json:"hull"`
		Drive int `json:"drive"`
	} `json:"player"`
	Camera struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"w"`
		Height int `json:"h"`
	} `json:"camera"`
	Tick int64 `json:"tick"`
}

func (c *Controller) Snapshot() ([]byte, error) {
	pos, _ := ecs.Get[components.Position](c.World, c.Player)
	ps, _ := ecs.Get[components.PlayerStats](c.World, c.Player)
	wi, _ := ecs.Get[components.WorldInfo](c.World, 1)
	cam, _ := ecs.Get[components.Camera](c.World, c.Player)
	s := Snapshot{}
	s.Player.X, s.Player.Y = int(pos.X), int(pos.Y)
	s.Player.Fuel, s.Player.Hull, s.Player.Drive = ps.Fuel, ps.Hull, ps.Drive
	s.Camera.X, s.Camera.Y = cam.X, cam.Y
	s.Camera.Width, s.Camera.Height = cam.Width, cam.Height
	s.Tick = wi.Tick
	return json.MarshalIndent(s, "", "  ")
}
