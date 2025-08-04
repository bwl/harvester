package ui

import (
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/components"
	"harvester/pkg/debug"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
	"harvester/pkg/systems"
	"harvester/pkg/timing"
)

// Removed mapRender - now using unified render system

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
	debug.Info("model", "Initializing new model with RNG")
	bs := engine.New(r)
	w := bs.World
	m := Model{
		rng:           r,
		world:         w,
		scheduler:     bs.Scheduler,
		render:        bs.Render,
		layoutManager: NewLayoutManager(120, 40),
	}
	m.player = bs.Player
	debug.Infof("model", "Model initialized with player entity %d", m.player)
	for y := 0; y < 80; y++ {
		for x := 0; x < 200; x++ {
			n := (x*73856093 ^ y*19349663) % 100
			if n < 6 {
				e := w.Create()
				ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
				tt := components.TileStar
				glyph := '*'
				switch {
				case n < 1:
					tt, glyph = components.TileGalaxyCore, '¤'
				case n < 2:
					tt, glyph = components.TileNebula, '≈'
				case n < 4:
					tt, glyph = components.TileAsteroid, '·'
				default:
					tt, glyph = components.TileStar, '*'
				}
				ecs.Add(w, e, components.Tile{Glyph: glyph, Type: tt})
				ecs.Add(w, e, components.Renderable{Glyph: glyph, TileType: tt, StyleMod: &components.ColorModifier{Special: components.EffectTwinkling}})
			}
			// rare comets
			if n == 42 {
				e := w.Create()
				ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
				ecs.Add(w, e, components.Tile{Glyph: '⤳', Type: components.TileComet})
				ecs.Add(w, e, components.Renderable{Glyph: '⤳', TileType: components.TileComet})
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
		debug.Debugf("ui", "Window resized to %dx%d", msg.Width, msg.Height)
		m.Width, m.Height = msg.Width, msg.Height
		m.layoutManager.Update(msg.Width, msg.Height)
	case tea.KeyMsg:
		// key handling moved to unified input router
	case time.Time:
		frameTimer := debug.StartSystemTimer("frame")
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

		// Record metrics
		debug.RecordFrame()
		debug.SetEntityCount(m.world.EntityCount())
		debug.UpdateMemoryStats()

		updateTimer := debug.StartSystemTimer("scheduler")
		m.scheduler.Update(1.0/20.0, m.world)
		updateTimer.Stop()

		m.frame++ // Increment frame counter for animations

		// Autosave is now handled at a higher level
		dur := time.Since(start)
		debug.RecordUpdateTime(dur)

		if os.Getenv("DEBUG_TICK") == "1" {
			debug.Debugf("performance", "Engine tick took %v", dur)
			m.log = append(m.log, "engine dt:0.05 ui dt:0.0167 tick:"+dur.String())
			if len(m.log) > 5 {
				m.log = m.log[len(m.log)-5:]
			}
		}

		next := ecs.GetWorldContext(m.world)
		if prev.CurrentLayer != next.CurrentLayer {
			debug.Infof("game", "Layer changed from %s to %s", layerName(prev.CurrentLayer), layerName(next.CurrentLayer))
			m.log = append(m.log, "Layer: "+layerName(next.CurrentLayer))
		}

		frameTimer.Stop()
		return m, tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return t })
	}
	return m, nil
}

/* moved to styles.go and layout.go */

func (m *Model) renderStatusBar(w int) string {
	// Get world dimensions for rendering
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
		Tick: int(timing.Tick()),
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

/*
	ctx := ecs.GetWorldContext(m.world)
	layout := m.layoutManager.GetLayout()
	dims := layout.Calculate()

	// Dynamic quest panel with state-aware styling
	questData := QuestPanelData{
		Status: fmt.Sprintf("%d/%d contracts", ctx.QuestProgress.ContractsCollected, ctx.QuestProgress.ContractsNeeded),
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
*/

func (m *Model) renderMap(mapW, mapH int) string {
	renderTimer := debug.StartSystemTimer("map_render")
	defer renderTimer.Stop()

	m.scheduler.Update(0, m.world)
	cam, _ := ecs.Get[components.Camera](m.world, m.player)
	debug.Debugf("render", "Rendering map %dx%d, camera at (%d, %d)", mapW, mapH, cam.X, cam.Y)
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
	// Use unified render system for all drawables (tiles and entities)
	for _, d := range m.render.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < mapW && y < mapH {
			canvas[y][x] = d.Glyph
		}
	}
	var b strings.Builder
	styled := make(map[[2]int]string, len(m.render.Output))
	for _, d := range m.render.Output {
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
	rightStr := ""
	status := ""

	// Convert log to LogMessage format for enhanced styling
	var logMessages []LogMessage
	for _, logLine := range m.log {
		logMessages = append(logMessages, LogMessage{
			Text: logLine,
			Type: LogInfo,
		})
	}
	_ = logMessages
	logStr := ""
	return m.renderUI(mapStr, rightStr, status, logStr)
}
