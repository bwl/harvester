package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/debug"
)

type DebugTab int

const (
	TabLogs DebugTab = iota
	TabPerformance
	TabECS
	TabSystem
)

func (t DebugTab) String() string {
	switch t {
	case TabLogs:
		return "Logs"
	case TabPerformance:
		return "Performance"
	case TabECS:
		return "ECS"
	case TabSystem:
		return "System"
	default:
		return "Unknown"
	}
}

type DebugPanel struct {
	width, height  int
	ActiveTab      DebugTab
	scrollOffset   int
	logFilter      debug.LogLevel
	categoryFilter string
}

func NewDebugPanel() *DebugPanel {
	return &DebugPanel{
		ActiveTab: TabLogs,
		logFilter: debug.LogDebug,
	}
}

func (dp *DebugPanel) SetDimensions(w, h int) {
	dp.width = w
	dp.height = h
}

func (dp *DebugPanel) NextTab() {
	dp.ActiveTab = (dp.ActiveTab + 1) % 4
	dp.scrollOffset = 0
}

func (dp *DebugPanel) PrevTab() {
	dp.ActiveTab = (dp.ActiveTab + 3) % 4
	dp.scrollOffset = 0
}

func (dp *DebugPanel) ScrollUp() {
	if dp.scrollOffset > 0 {
		dp.scrollOffset--
	}
}

func (dp *DebugPanel) ScrollDown() {
	dp.scrollOffset++
}

func (dp *DebugPanel) CycleLogFilter() {
	switch dp.logFilter {
	case debug.LogDebug:
		dp.logFilter = debug.LogInfo
	case debug.LogInfo:
		dp.logFilter = debug.LogWarn
	case debug.LogWarn:
		dp.logFilter = debug.LogError
	case debug.LogError:
		dp.logFilter = debug.LogDebug
	}
}

func (dp *DebugPanel) Render() string {
	if dp.width == 0 || dp.height == 0 {
		return ""
	}

	// Main panel style - full screen overlay
	panelStyle := lipgloss.NewStyle().
		Width(dp.width).
		Height(dp.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("75")).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("250")).
		Padding(1)

	// Header with tabs
	header := dp.renderTabs()

	// Content area dimensions (accounting for header, footer, padding, and border)
	contentHeight := dp.height - 8 // border + padding + header + footer
	contentWidth := dp.width - 6   // border + padding + content padding

	// Content based on active tab (render with inner dimensions)
	innerWidth := contentWidth - 2   // account for content padding
	innerHeight := contentHeight - 2 // account for content padding

	var content string
	switch dp.ActiveTab {
	case TabLogs:
		content = dp.renderLogsTab(innerWidth, innerHeight)
	case TabPerformance:
		content = dp.renderPerformanceTab(innerWidth, innerHeight)
	case TabECS:
		content = dp.renderECSTab(innerWidth, innerHeight)
	case TabSystem:
		content = dp.renderSystemTab(innerWidth, innerHeight)
	}

	// Wrap content with consistent styling
	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("250")).
		Padding(1)

	styledContent := contentStyle.Render(content)

	// Footer with help
	footer := dp.renderFooter(contentWidth)

	// Combine all parts
	fullContent := lipgloss.JoinVertical(lipgloss.Left,
		header,
		styledContent,
		footer,
	)

	return panelStyle.Render(fullContent)
}

func (dp *DebugPanel) renderTabs() string {
	var tabs []string
	for i := TabLogs; i <= TabSystem; i++ {
		tabStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("246"))

		if i == dp.ActiveTab {
			tabStyle = tabStyle.
				Background(lipgloss.Color("75")).
				Foreground(lipgloss.Color("254")).
				Bold(true)
		}

		tabs = append(tabs, tabStyle.Render(i.String()))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (dp *DebugPanel) renderLogsTab(width, height int) string {
	entries := debug.GetEntriesFiltered(dp.logFilter, dp.categoryFilter)

	// Title with filter info
	title := fmt.Sprintf("Debug Logs (Level: %s, Total: %d)", dp.logFilter.String(), len(entries))
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true).
		MarginBottom(1)

	if len(entries) == 0 {
		noLogsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Italic(true)
		return lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(title),
			noLogsStyle.Render("No log entries match current filter"),
		)
	}

	// Calculate visible range
	maxVisible := height - 3 // account for title and spacing
	startIdx := len(entries) - maxVisible - dp.scrollOffset
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + maxVisible
	if endIdx > len(entries) {
		endIdx = len(entries)
	}

	// Render visible entries
	var lines []string
	lines = append(lines, titleStyle.Render(title))

	for i := startIdx; i < endIdx; i++ {
		entry := entries[i]
		lineStyle := lipgloss.NewStyle().Foreground(dp.getLogLevelColor(entry.Level))

		// Truncate long lines
		line := entry.String()
		if len(line) > width-2 {
			line = line[:width-5] + "..."
		}

		lines = append(lines, lineStyle.Render(line))
	}

	// Add scroll indicator if needed
	if len(entries) > maxVisible {
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d (offset: %d)",
			startIdx+1, endIdx, len(entries), dp.scrollOffset)
		scrollStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Italic(true)
		lines = append(lines, scrollStyle.Render(scrollInfo))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (dp *DebugPanel) renderPerformanceTab(width, height int) string {
	metrics := debug.GetMetricsSnapshot()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("155"))

	var lines []string
	lines = append(lines, titleStyle.Render("Performance Metrics"))

	// Frame metrics
	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Frame Rate:"),
			valueStyle.Render(fmt.Sprintf("%.1f FPS", metrics.FrameRate)),
		),
	)

	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Frame Time:"),
			valueStyle.Render(fmt.Sprintf("%v", metrics.FrameTime)),
		),
	)

	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Frame Count:"),
			valueStyle.Render(fmt.Sprintf("%d", metrics.FrameCount)),
		),
	)

	// Memory metrics
	lines = append(lines, "")
	lines = append(lines, titleStyle.Render("Memory"))

	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Heap Alloc:"),
			valueStyle.Render(formatBytes(metrics.HeapAlloc)),
		),
	)

	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Heap Sys:"),
			valueStyle.Render(formatBytes(metrics.HeapSys)),
		),
	)

	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Goroutines:"),
			valueStyle.Render(fmt.Sprintf("%d", metrics.NumGoroutines)),
		),
	)

	// System timing
	if len(metrics.SystemTime) > 0 {
		lines = append(lines, "")
		lines = append(lines, titleStyle.Render("System Timing"))

		for system, duration := range metrics.SystemTime {
			lines = append(lines,
				lipgloss.JoinHorizontal(lipgloss.Left,
					labelStyle.Render(system+":"),
					valueStyle.Render(duration.String()),
				),
			)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (dp *DebugPanel) renderECSTab(width, height int) string {
	metrics := debug.GetMetricsSnapshot()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("155"))

	var lines []string
	lines = append(lines, titleStyle.Render("ECS Information"))

	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Entity Count:"),
			valueStyle.Render(fmt.Sprintf("%d", metrics.EntityCount)),
		),
	)

	// TODO: Add more ECS-specific debug info when available
	// - Component counts by type
	// - System execution order
	// - Entity composition analysis

	lines = append(lines, "")
	lines = append(lines,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Italic(true).
			Render("More ECS debug info coming soon..."),
	)

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (dp *DebugPanel) renderSystemTab(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("155"))

	var lines []string
	lines = append(lines, titleStyle.Render("System Information"))

	lines = append(lines,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Time:"),
			valueStyle.Render(time.Now().Format("15:04:05")),
		),
	)

	// Available categories
	categories := debug.GetCategories()
	if len(categories) > 0 {
		lines = append(lines, "")
		lines = append(lines, titleStyle.Render("Log Categories"))

		for _, cat := range categories {
			lines = append(lines,
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("155")).
					Render("• "+cat),
			)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (dp *DebugPanel) renderFooter(width int) string {
	helpStyle := lipgloss.NewStyle().
		Width(width).
		Background(lipgloss.Color("238")).
		Foreground(lipgloss.Color("246")).
		Padding(0, 1).
		Italic(true)

	var help []string
	switch dp.ActiveTab {
	case TabLogs:
		help = []string{"↑↓: scroll", "f: filter level", "c: clear logs"}
	case TabPerformance:
		help = []string{"r: refresh", "↑↓: scroll"}
	default:
		help = []string{"↑↓: scroll"}
	}

	commonHelp := []string{"tab: next tab", "shift+tab: prev tab", "F12: close"}
	help = append(help, commonHelp...)

	return helpStyle.Render(strings.Join(help, " | "))
}

func (dp *DebugPanel) getLogLevelColor(level debug.LogLevel) lipgloss.Color {
	switch level {
	case debug.LogDebug:
		return lipgloss.Color("246") // dim gray
	case debug.LogInfo:
		return lipgloss.Color("250") // light gray (instead of white)
	case debug.LogWarn:
		return lipgloss.Color("214") // orange
	case debug.LogError:
		return lipgloss.Color("196") // red
	default:
		return lipgloss.Color("250")
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
