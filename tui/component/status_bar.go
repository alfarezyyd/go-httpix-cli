package component

import (
	"go-httpix-cli/config"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusBarProps is everything the status bar needs to render.
type StatusBarProps struct {
	Width      int
	FocusLabel string
	Keys       config.KeyMap
}

// StatusBar renders the bottom hint bar.
func StatusBar(p StatusBarProps) string {
	focusIndicator := lipgloss.NewStyle().
		Foreground(config.Crust).Background(config.Lavender).Bold(true).Padding(0, 1).
		Render("  " + p.FocusLabel + " ")

	hints := []struct{ key, desc string }{
		{p.Keys.LabelSend, "Send"},
		{p.Keys.LabelFocus, "Focus"},
		{p.Keys.LabelMethod, "Method"},
		{p.Keys.LabelTab, "Tab"},
		{p.Keys.LabelFormatJSON, "Format JSON"},
		{p.Keys.LabelQuit, "Quit"},
	}

	var parts []string
	for _, h := range hints {
		k := config.KeyBadge.Render(" " + h.key + " ")
		d := config.KeyDesc.Render(" " + h.desc + " ")
		parts = append(parts, k+d)
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Center,
		focusIndicator, " ", strings.Join(parts, " "),
	)

	return lipgloss.NewStyle().
		Background(config.Crust).Width(p.Width).Padding(0, 1).
		Render(bar)
}
