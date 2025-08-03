package ui

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
	"harvester/pkg/systems"
)

var mapRender *systems.MapRender

type Model struct {
	Width, Height int
	log           []string
	rng           *rand.Rand

	world     *ecs.World
	scheduler interface {
		Update(dt float64, w *ecs.World)
	}
	render        *systems.Render
	player        ecs.Entity
	layoutManager *LayoutManager
	frame         int
	prevStats     PlayerStatsData
}

func (m *Model) World() *ecs.World { return m.world }

func NewModel(gs any) Model { return NewModelWithRNG(rand.New(rand.NewSource(1))) }

func NewModelWithRNG(r *rand.Rand) Model {
	bs := engine.New(r)
	w := bs.World
	m := Model{
		rng:           r,
		world:         w,
		scheduler:     bs.Scheduler,
		render:        bs.Render,
		layoutManager: NewLayoutManager(120, 40),
	}
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
		m.layoutManager.Update(msg.Width, msg.Height)
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
		
		// Store previous stats for trend calculation
		if ps, exists := ecs.Get[components.PlayerStats](m.world, m.player); exists {
			m.prevStats = PlayerStatsData{
				Fuel:  ps.Fuel,
				Hull:  ps.Hull,
				Drive: ps.Drive,
			}
		}
		
		m.scheduler.Update(1.0/20.0, m.world)
		m.frame++ // Increment frame counter for animations
		
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

/* moved to styles.go and layout.go */

func (m *Model) renderStatusBar(w int) string {
	wi, _ := ecs.Get[components.WorldInfo](m.world, 1)
	ps, _ := ecs.Get[components.PlayerStats](m.world, m.player)
	ctx := ecs.GetWorldContext(m.world)

	location := LocationData{
		Layer:  layerName(ctx.CurrentLayer),
		Planet: ctx.PlanetID,
		Depth:  ctx.Depth,
	}

	currentStats := PlayerStatsData{
		Fuel:  ps.Fuel,
		Hull:  ps.Hull,
		Drive: ps.Drive,
	}

	info := GameInfoData{
		Tick: int(wi.Tick),
	}

	// Use advanced stats with trends and animations
	left := lipgloss.JoinHorizontal(lipgloss.Center, 
		LocationComponent(location), 
		Muted("  |  "), 
		GameInfoComponent(info))
	
	right := AdvancedPlayerStatsComponent(currentStats, m.prevStats, m.frame)
	
	return NewStyleBuilder().
		Width(w).
		Background(GetCurrentTheme().Bg).
		Foreground(GetCurrentTheme().Text).
		PaddingHorizontal(1).
		Render(lipgloss.JoinHorizontal(lipgloss.Center,
			left,
			strings.Repeat(" ", max(0, w-lipgloss.Width(left)-lipgloss.Width(right)-4)),
			right,
		))
}



func (m *Model) renderRightPanel() string {
	ctx := ecs.GetWorldContext(m.world)
	layout := m.layoutManager.GetLayout()
	dims := layout.Calculate()

	// Dynamic quest panel with state-aware styling
	questData := QuestPanelData{
		Status: royalCharterStatus(ctx.QuestProgress),
	}
	
	// Determine quest state based on progress
	questState := StateNormal
	if ctx.QuestProgress.RoyalCharterComplete {
		questState = StateSuccess
	} else if ctx.QuestProgress.ContractsCollected == 0 {
		questState = StateWarning
	}

	controlGroups := []ControlsGroup{
		{
			Title: "Movement",
			Items: []ControlItem{
				{"h j k l", "or arrows"},
			},
		},
		{
			Title: "Actions",
			Items: []ControlItem{
				{"> ", "enter/descend"},
				{"q ", "quit"},
			},
		},
		{
			Title: "Save",
			Items: []ControlItem{
				{"Ctrl+S", "manual save"},
				{"1-3", "save slots"},
			},
		},
	}

	// Use dynamic quest panel and responsive controls
	quest := DynamicQuestPanel(questData, questState)
	controls := ResponsiveControlsPanel(controlGroups, dims.RightWidth)

	return lipgloss.JoinVertical(lipgloss.Left, quest, controls)
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
				ch = Muted(".")
			}
			b.WriteString(ch)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (m *Model) renderUI(mapStr, rightStr, statusStr, logStr string) string {
	return m.layoutManager.RenderWithLayout(mapStr, rightStr, statusStr, logStr)
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
	// Update layout manager if dimensions changed
	w, h := m.Width, m.Height
	if w == 0 {
		w = 120
	}
	if h == 0 {
		h = 40
	}
	
	// Get layout dimensions
	layout := m.layoutManager.GetLayout()
	dims := layout.Calculate()
	
	// Render components
	mapStr := m.renderMap(dims.MapWidth, dims.MapHeight)
	rightStr := m.renderRightPanel()
	status := m.renderStatusBar(dims.ContentWidth)
	
	// Convert log to LogMessage format for enhanced styling
	var logMessages []LogMessage
	for _, logLine := range m.log {
		logMessages = append(logMessages, LogMessage{
			Text: logLine,
			Type: LogInfo,
		})
	}
	logStr := LogPanel(logMessages, dims.LogWidth)
	
	return m.renderUI(mapStr, rightStr, status, logStr)
}
