package component

import (
	"go-httpix-cli/config"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TopBarProps is everything the top bar needs to render.
type TopBarProps struct {
	Width   int
	IsMac   bool
	Loading bool
	Tick    int
	Keys    config.KeyMap
}

// TopBar renders the application title bar.
func TopBar(p TopBarProps) string {
	logo := config.Logo.Render("httpix_<")

	platLabel, platColor := " 🐧 Linux ", config.Sky
	if p.IsMac {
		platLabel, platColor = "  macOS ", config.Flamingo
	}

	platBadge := lipgloss.NewStyle().
		Foreground(config.Crust).
		MarginLeft(1).
		Background(platColor).Bold(true).
		Render(platLabel)

	version := config.Version.Render(" v1.0 · TUI HTTP Client ")

	spinner := ""
	if p.Loading {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spinner = " " + lipgloss.NewStyle().Foreground(config.Mauve).Render(frames[p.Tick%len(frames)])
	}

	right := lipgloss.NewStyle().
		Foreground(config.Overlay0).Background(config.Mantle).
		Render(" " + p.Keys.LabelSend + " Send · Tab Focus · " + p.Keys.LabelQuit + " Quit" + spinner + " ")

	left := logo + platBadge + version
	gap := strings.Repeat(" ", max(0, p.Width-lipgloss.Width(left+right)))

	return config.TopBar.Width(p.Width).Render(left + gap + right)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
