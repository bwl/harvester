package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"bubbleRouge/pkg/engine"
)

type Model struct {
	GS *engine.GameState
	Width, Height int
	log []string
}

func NewModel(gs *engine.GameState) Model { return Model{GS: gs, log: gs.Log} }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "h", "left":
			m.GS.Step(engine.ActMoveLeft)
		case "l", "right":
			m.GS.Step(engine.ActMoveRight)
		case "k", "up":
			m.GS.Step(engine.ActMoveUp)
		case "j", "down":
			m.GS.Step(engine.ActMoveDown)
		case "g":
			m.GS.Step(engine.ActHarvest)
		}
	}
	return m, nil
}

var (
	spaceStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
	galStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render
	playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Render
)

func (m Model) View() string {
	if m.GS == nil { return "" }
	w, h := m.Width, m.Height
	if w == 0 { w = 120 }
	if h == 0 { h = 40 }
	// reserve panels
	mapW, mapH := w-30, h-5
	if mapW < 10 { mapW = 10 }
	if mapH < 10 { mapH = 10 }
	mx0 := m.GS.Player.Pos.X - mapW/2
	my0 := m.GS.Player.Pos.Y - mapH/2
	var b strings.Builder
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			gx, gy := mx0+x, my0+y
			ch := " "
			if gx == m.GS.Player.Pos.X && gy == m.GS.Player.Pos.Y {
				ch = playerStyle("@")
			} else if gx >= 0 && gy >= 0 && gx < m.GS.Map.Width && gy < m.GS.Map.Height {
				k := m.GS.Map.Tiles[gy][gx].Kind
				if k == engine.Galaxy {
					ch = galStyle("*")
				} else {
					ch = spaceStyle(".")
				}
			}
			b.WriteString(ch)
		}
		b.WriteByte('\n')
	}
	status := fmt.Sprintf("tick:%d size:%dx%d fuel:%d hull:%d drive:%d", m.GS.Tick, m.GS.Map.Width, m.GS.Map.Height, m.GS.Player.Fuel, m.GS.Player.Hull, m.GS.Player.DriveLevel)
	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(mapW).Height(mapH).Render(b.String()),
		lipgloss.NewStyle().Width(30).Render(status+"\n\n"+strings.Join(m.GS.Log, "\n")),
	)
}
