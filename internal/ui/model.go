package ui

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/systems"
)

func itoa(i int) string { return strconv.Itoa(i) }

var mapRender *systems.MapRender

type Model struct {
	Width, Height int
	log           []string
	rng           *rand.Rand

	world     *ecs.World
	scheduler interface {
		Update(dt float64, w *ecs.World)
	}
	render *systems.Render
	player ecs.Entity
}

func NewModel(gs any) Model { return NewModelWithRNG(rand.New(rand.NewSource(1))) }

func NewModelWithRNG(r *rand.Rand) Model {
	w := ecs.NewWorld(r)
	render := &systems.Render{}
	mapRender = &systems.MapRender{}
	camera := &systems.CameraSystem{}
	reg := ecs.SystemRegistry{
		UniversalSystems: []ecs.System{systems.InputSystem{}, systems.Tick{}, camera, systems.LevelManager{}, mapRender, render},
		SpaceSystems:     []ecs.System{systems.SpaceMovement{}, systems.FuelSystem{}, systems.PlanetApproachSystem{}, systems.PlanetSelection{}},
		SurfaceSystems:   []ecs.System{systems.SurfaceHeartbeat{}, systems.TerrainGen{}, systems.SurfaceMovement{}, systems.DepthProgression{}, systems.WeatherTick{}, systems.RiverFlow{}, systems.TradeRoutePatrols{}, systems.WildlifeSpawn{}, systems.KingdomGuards{}, systems.QuestSystem{}},
	}
	s := ecs.NewSchedulerWithContext(reg)
	m := Model{rng: r, world: w, scheduler: s, render: render}
	m.player = w.Create()
	ecs.Add(w, m.player, components.Position{})
	ecs.Add(w, m.player, components.Camera{Width: 200 - 30, Height: 80 - 5})
	ecs.Add(w, m.player, components.Renderable{Glyph: '@', TileType: components.TileStar})
	ecs.SetWorldContext(w, ecs.WorldContext{CurrentLayer: ecs.LayerSpace})
	ecs.Add(w, m.player, components.Input{})
	ecs.Add(w, m.player, components.Velocity{})
	ecs.Add(w, m.player, systems.FuelTank{Current: 100})
	ecs.Add(w, m.player, systems.Velocity{})
	ecs.Add(w, 1, components.WorldInfo{Width: 200, Height: 80})
	for y := 0; y < 80; y++ {
		for x := 0; x < 200; x++ {
			if (x+y)%11 == 0 {
				e := w.Create()
				ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
				ecs.Add(w, e, components.Tile{Glyph: '*', Type: components.TileStar})
				ecs.Add(w, e, components.Renderable{Glyph: '*', TileType: components.TileStar, StyleMod: &components.ColorModifier{Special: components.EffectTwinkling}})
			}
		}
	}
	ecs.Add(w, m.player, components.PlayerStats{Fuel: 100, Hull: 100, Drive: 1})
	return m
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "h", "left":
			systems.SetPlayerInput(m.world, m.player, "left")
			systems.ApplyDirectionalVelocity(m.world, m.player, -1, 0)
		case "l", "right":
			systems.SetPlayerInput(m.world, m.player, "right")
			systems.ApplyDirectionalVelocity(m.world, m.player, 1, 0)
		case "k", "up":
			systems.SetPlayerInput(m.world, m.player, "up")
			systems.ApplyDirectionalVelocity(m.world, m.player, 0, -1)
		case "j", "down":
			systems.SetPlayerInput(m.world, m.player, "down")
			systems.ApplyDirectionalVelocity(m.world, m.player, 0, 1)
		case ">":
			systems.SetPlayerInput(m.world, m.player, "enter")
		case "ctrl+s":
			m.save()
		case "ctrl+o":
			m.load()
		case "ctrl+shift+s":
			m.saveCompressed()
		case "ctrl+shift+o":
			m.loadCompressed()
		case "1":
			m.saveSlot(1)
		case "2":
			m.saveSlot(2)
		case "3":
			m.saveSlot(3)
		case "shift+1":
			m.loadSlot(1)
		case "shift+2":
			m.loadSlot(2)
		case "shift+3":
			m.loadSlot(3)
		}
	}
	prev := ecs.GetWorldContext(m.world)
	m.scheduler.Update(1.0, m.world)
	next := ecs.GetWorldContext(m.world)
	if prev.CurrentLayer != next.CurrentLayer {
		m.log = append(m.log, "Layer: "+layerName(next.CurrentLayer))
	}
	return m, nil
}

var (
	spaceStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
	galStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render
	playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Render
)

func (m *Model) save() {
	ss, _ := ecs.Save(m.world, nil)
	b, _ := json.Marshal(ss)
	_ = os.MkdirAll(".saves", 0o755)
	_ = os.WriteFile(filepath.Join(".saves", "autosave.json"), b, 0o644)
}
func (m *Model) load() {
	b, err := os.ReadFile(filepath.Join(".saves", "autosave.json"))
	if err != nil {
		return
	}
	var s ecs.Snapshot
	if json.Unmarshal(b, &s) != nil {
		return
	}
	_ = ecs.Load(m.world, &s, nil)
}

func (m *Model) saveCompressed() {
	ss, _ := ecs.Save(m.world, nil)
	b, _ := ecs.EncodeSnapshot(ss, ecs.SaveOptions{Compress: true})
	_ = os.MkdirAll(".saves", 0o755)
	_ = os.WriteFile(filepath.Join(".saves", "autosave.gz"), b, 0o644)
}

func (m *Model) loadCompressed() {
	b, err := os.ReadFile(filepath.Join(".saves", "autosave.gz"))
	if err != nil {
		return
	}
	s, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true})
	if err != nil {
		return
	}
	_ = ecs.Load(m.world, s, nil)
}

func (m *Model) saveSlot(n int) {
	ss, _ := ecs.Save(m.world, nil)
	b, _ := ecs.EncodeSnapshot(ss, ecs.SaveOptions{Compress: true})
	_ = os.MkdirAll(".saves", 0o755)
	_ = os.WriteFile(filepath.Join(".saves", "slot"+itoa(n)+".gz"), b, 0o644)
}

func (m *Model) loadSlot(n int) {
	b, err := os.ReadFile(filepath.Join(".saves", "slot"+itoa(n)+".gz"))
	if err != nil {
		return
	}
	s, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true})
	if err != nil {
		return
	}
	_ = ecs.Load(m.world, s, nil)
}

func (m Model) View() string {
	w, h := m.Width, m.Height
	if w == 0 {
		w = 120
	}
	if h == 0 {
		h = 40
	}
	// reserve panels
	mapW, mapH := w-30, h-5
	m.scheduler.Update(0, m.world)
	if mapW < 10 {
		mapW = 10
	}
	if mapH < 10 {
		mapH = 10
	}
	cam, _ := ecs.Get[components.Camera](m.world, m.player)
	mx0 := cam.X
	my0 := cam.Y
	canvas := make([][]rune, mapH)
	for i := range canvas {
		canvas[i] = make([]rune, mapW)
	}
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			canvas[y][x] = '.'
		}
	}
	// draw background tiles (with styles ignored in canvas)
	for _, d := range mapRender.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < mapW && y < mapH {
			canvas[y][x] = d.Glyph
		}
	}
	// draw entities
	for _, d := range m.render.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < mapW && y < mapH {
			canvas[y][x] = d.Glyph
		}
	}
	var b strings.Builder
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			written := false
			for _, d := range mapRender.Output {
				if d.X-mx0 == x && d.Y-my0 == y {
					b.WriteString(d.Style.Render(string(d.Glyph)))
					written = true
					break
				}
			}
			if !written {
				ch := string(canvas[y][x])
				if ch == "." {
					ch = spaceStyle(".")
				}
				b.WriteString(ch)
			}
		}
		b.WriteByte('\n')
	}
	wi, _ := ecs.Get[components.WorldInfo](m.world, 1)
	ps, _ := ecs.Get[components.PlayerStats](m.world, m.player)
	ctx := ecs.GetWorldContext(m.world)
	top := lipgloss.NewStyle().Width(w).Render(strings.Join([]string{
		"Layer " + layerName(ctx.CurrentLayer),
		"Planet " + itoa(ctx.PlanetID),
		"Depth " + itoa(ctx.Depth),
		"Tick " + itoa(int(wi.Tick)),
		"Fuel " + itoa(ps.Fuel) + "  Hull " + itoa(ps.Hull) + "  Drive " + itoa(ps.Drive),
	}, "  |  "))
	right := lipgloss.NewStyle().Width(30).Render(strings.Join([]string{
		"Quest:", royalCharterStatus(ctx.QuestProgress),
		"",
		"Keys:", "h/j/k/l move", "> enter", "q quit",
	}, "\n"))
	main := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(mapW).Height(mapH).Render(b.String()),
		right,
	)
	log := lipgloss.NewStyle().Width(w).Render(strings.Join(m.log, "\n"))
	frame := lipgloss.JoinVertical(lipgloss.Left,
		top,
		main,
		log,
	)
	return frame
}
