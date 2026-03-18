package component

import (
	"fmt"
	"go-httpix-cli/config"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// HistoryRow is a single entry shown in the history sidebar.
type HistoryRow struct {
	Method   string
	URL      string
	Status   int
	Duration time.Duration
}

// SidebarProps is everything the sidebar needs to render.
type SidebarProps struct {
	Width     int
	Height    int
	Focused   bool
	Entries   []HistoryRow
	ActiveIdx int // -1 = nothing selected
}

// Sidebar renders the request-history panel.
func Sidebar(p SidebarProps) string {
	title := config.PanelTitleStyle.Render("◈ History")

	var rows []string
	if len(p.Entries) == 0 {
		rows = append(rows, lipgloss.NewStyle().
			Foreground(config.Overlay0).Italic(true).Width(p.Width-4).
			Render("\n  No history yet\n"))
	} else {
		rows = historyRows(p)
	}

	ps := config.PanelStyle
	if p.Focused {
		ps = config.PanelFocusedStyle
	}

	return ps.Width(p.Width).Height(p.Height-6).Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			append([]string{title}, rows...)...,
		))
}

func historyRows(p SidebarProps) []string {
	maxVisible := max(1, p.Height-20)
	start := 0
	if p.ActiveIdx >= maxVisible {
		start = p.ActiveIdx - maxVisible + 1
	}
	end := min(len(p.Entries), start+maxVisible)

	var rows []string
	for i := start; i < end; i++ {
		rows = append(rows, historyRow(p.Entries[i], i == p.ActiveIdx, p.Width))
		rows = append(rows, config.Divider.Render(strings.Repeat("─", p.Width-4)))
	}
	return rows
}

func historyRow(e HistoryRow, active bool, w int) string {
	col, ok := config.MethodColor[e.Method]
	if !ok {
		col = config.Overlay2
	}
	methBadge := lipgloss.NewStyle().
		Foreground(config.Crust).Background(col).Bold(true).
		Render(" " + e.Method + " ")

	statusColor := config.Green
	if e.Status >= 400 {
		statusColor = config.Red
	} else if e.Status >= 300 {
		statusColor = config.Yellow
	}
	statusBadge := lipgloss.NewStyle().Foreground(statusColor).Bold(true).
		Render(fmt.Sprintf("%d", e.Status))

	urlShort := e.URL
	if maxW := w - 10; len(urlShort) > maxW {
		urlShort = "…" + urlShort[len(urlShort)-maxW+1:]
	}

	line1 := lipgloss.JoinHorizontal(lipgloss.Center, methBadge, " ", statusBadge)
	line2 := lipgloss.NewStyle().Foreground(config.Subtext1).Width(w - 4).Render(urlShort)
	line3 := lipgloss.NewStyle().Foreground(config.Overlay0).Italic(true).
		Render("  " + e.Duration.Round(time.Millisecond).String())

	entry := lipgloss.JoinVertical(lipgloss.Left, line1, line2, line3)

	if active {
		return config.HistoryActive.Width(w - 4).Render(entry)
	}
	return config.HistoryItem.Width(w - 4).Render(entry)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
