package ui

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
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

func (m *Model) World() *ecs.World { return m.world }

func NewModel(gs any) Model { return NewModelWithRNG(rand.New(rand.NewSource(1))) }

func NewModelWithRNG(r *rand.Rand) Model {
	bs := engine.New(r)
	w := bs.World
	m := Model{rng: r, world: w, scheduler: bs.Scheduler, render: bs.Render}
	mapRender = bs.MapRender
	m.player = bs.Player
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

func (m *Model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return t })
}

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
		case "ctrl+shift+s":
			m.saveCompressed()
		case "1":
			m.saveSlot(1)
		case "2":
			m.saveSlot(2)
		case "3":
			m.saveSlot(3)
		}
	case time.Time:
		prev := ecs.GetWorldContext(m.world)
		start := time.Now()
		m.scheduler.Update(1.0/20.0, m.world)
		wi, _ := ecs.Get[components.WorldInfo](m.world, 1)
		if int(wi.Tick)%100 == 0 {
			m.saveCompressed()
		}
		dur := time.Since(start)
		if os.Getenv("DEBUG_TICK") == "1" {
			m.log = append(m.log, "engine dt:0.05 ui dt:0.0167 tick:"+dur.String())
			if len(m.log) > 5 {
				m.log = m.log[len(m.log)-5:]
			}
		}
		next := ecs.GetWorldContext(m.world)
		if prev.CurrentLayer != next.CurrentLayer {
			m.log = append(m.log, "Layer: "+layerName(next.CurrentLayer))
		}
		return m, tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return t })
	}
	return m, nil
}

var (
	spaceStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
	galStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render
	playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Render
)

const (
	rightPanelWidth   = 30
	bottomPanelHeight = 5
	minMapW           = 10
	minMapH           = 10
)

func (m *Model) renderStatusBar(w int) string {
	wi, _ := ecs.Get[components.WorldInfo](m.world, 1)
	ps, _ := ecs.Get[components.PlayerStats](m.world, m.player)
	ctx := ecs.GetWorldContext(m.world)
	return lipgloss.NewStyle().Width(w).Render(strings.Join([]string{
		"Layer " + layerName(ctx.CurrentLayer),
		"Planet " + itoa(ctx.PlanetID),
		"Depth " + itoa(ctx.Depth),
		"Tick " + itoa(int(wi.Tick)),
		"Fuel " + itoa(ps.Fuel) + "  Hull " + itoa(ps.Hull) + "  Drive " + itoa(ps.Drive),
	}, "  |  "))
}

func (m *Model) renderMap(mapW, mapH int) string {
	m.scheduler.Update(0, m.world)
	cam, _ := ecs.Get[components.Camera](m.world, m.player)
	mx0, my0 := cam.X, cam.Y
	canvas := make([][]rune, mapH)
	for i := range canvas {
		canvas[i] = make([]rune, mapW)
	}
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			canvas[y][x] = '.'
		}
	}
	for _, d := range mapRender.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < mapW && y < mapH {
			canvas[y][x] = d.Glyph
		}
	}
	for _, d := range m.render.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < mapW && y < mapH {
			canvas[y][x] = d.Glyph
		}
	}
	var b strings.Builder
	styled := make(map[[2]int]string, len(mapRender.Output))
	for _, d := range mapRender.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < mapW && y < mapH {
			styled[[2]int{x, y}] = d.Style.Render(string(d.Glyph))
		}
	}
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			if s, ok := styled[[2]int{x, y}]; ok {
				b.WriteString(s)
				continue
			}
			ch := string(canvas[y][x])
			if ch == "." {
				ch = spaceStyle(".")
			}
			b.WriteString(ch)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (m *Model) renderUI(w, h int, mapStr, rightStr, statusStr, logStr string) string {
	main := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(w-rightPanelWidth).Height(h-bottomPanelHeight).Render(mapStr),
		lipgloss.NewStyle().Width(rightPanelWidth).Render(rightStr),
	)
	return lipgloss.JoinVertical(lipgloss.Left,
		statusStr,
		main,
		lipgloss.NewStyle().Width(w).Render(logStr),
	)
}

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
	mapW, mapH := w-rightPanelWidth, h-bottomPanelHeight
	if mapW < minMapW {
		mapW = minMapW
	}
	if mapH < minMapH {
		mapH = minMapH
	}
	mapStr := m.renderMap(mapW, mapH)
	right := strings.Join([]string{
		"Quest:", royalCharterStatus(ecs.GetWorldContext(m.world).QuestProgress),
		"",
		"Keys:", "h/j/k/l move", "> enter", "q quit",
	}, "\n")
	status := m.renderStatusBar(w)
	log := strings.Join(m.log, "\n")
	return m.renderUI(w, h, mapStr, right, status, log)
}
