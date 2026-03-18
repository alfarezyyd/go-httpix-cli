package tui

import (
	"go-httpix-cli/tui/component"

	"github.com/charmbracelet/lipgloss"
)

const sidebarWidth = 28

// View is the Elm-architecture view function.
// It composes sub-views into the final terminal string.
func (m Model) View() string {
	if m.Width == 0 {
		return "Loading…"
	}

	mainW := max(40, m.Width-sidebarWidth-3)

	main := lipgloss.JoinVertical(lipgloss.Left,
		component.URLRow(m.urlRowProps(mainW)),
		component.RequestPanel(m.requestPanelProps(mainW)),
		component.ResponsePanel(m.responsePanelProps(mainW)),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		component.TopBar(m.topBarProps()),
		lipgloss.JoinHorizontal(lipgloss.Top, main, " ", component.Sidebar(m.sidebarProps())),
		component.StatusBar(m.statusBarProps()),
	)
}
