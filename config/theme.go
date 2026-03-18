// Package config holds app-wide constants: colour palette and lipgloss styles.
// Nothing in here has side-effects or depends on runtime state.
package config

import "github.com/charmbracelet/lipgloss"

// ─────────────────────────────────────────────────────────────
//  COLOUR PALETTE  (Catppuccin Mocha)
// ─────────────────────────────────────────────────────────────

var (
	Crust    = lipgloss.Color("#11111b")
	Mantle   = lipgloss.Color("#181825")
	Base     = lipgloss.Color("#1e1e2e")
	Surface0 = lipgloss.Color("#313244")
	Surface1 = lipgloss.Color("#45475a")
	Overlay0 = lipgloss.Color("#6c7086")
	Overlay2 = lipgloss.Color("#9399b2")
	Text     = lipgloss.Color("#cdd6f4")
	Subtext1 = lipgloss.Color("#bac2de")
	Lavender = lipgloss.Color("#b4befe")
	Blue     = lipgloss.Color("#89b4fa")
	Sky      = lipgloss.Color("#89dceb")
	Teal     = lipgloss.Color("#94e2d5")
	Green    = lipgloss.Color("#a6e3a1")
	Yellow   = lipgloss.Color("#f9e2af")
	Peach    = lipgloss.Color("#fab387")
	Red      = lipgloss.Color("#f38ba8")
	Mauve    = lipgloss.Color("#cba6f7")
	Flamingo = lipgloss.Color("#f2cdcd")
)

// MethodColor maps an HTTP method name to its accent colour.
var MethodColor = map[string]lipgloss.Color{
	"GET":     Green,
	"POST":    Blue,
	"PUT":     Yellow,
	"PATCH":   Peach,
	"DELETE":  Red,
	"HEAD":    Sky,
	"OPTIONS": Teal,
}

// ─────────────────────────────────────────────────────────────
//  STYLES
// ─────────────────────────────────────────────────────────────

var (
	// ── Top bar ──────────────────────────────────────────────
	TopBar  = lipgloss.NewStyle().Background(Mantle).Foreground(Text).Padding(0, 2)
	Logo    = lipgloss.NewStyle().Foreground(Mauve).Bold(true).Background(Mantle)
	Version = lipgloss.NewStyle().Foreground(Overlay0).Background(Mantle).Italic(true)

	// ── URL bar ──────────────────────────────────────────────
	URLBar = lipgloss.NewStyle().
		Background(Surface0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Surface1).
		Padding(0, 1)

	URLBarFocused = lipgloss.NewStyle().
			Background(Surface0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Lavender).
			Padding(0, 1)

	// ── Panels ───────────────────────────────────────────────
	PanelStyle = lipgloss.NewStyle().
			Background(Mantle).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Surface1)

	PanelFocusedStyle = lipgloss.NewStyle().
				Background(Mantle).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Mauve)

	PanelTitleStyle = lipgloss.NewStyle().
			Foreground(Mauve).Bold(true).Background(Mantle).Padding(0, 1)

	// ── Tabs ─────────────────────────────────────────────────
	Tab = lipgloss.NewStyle().Foreground(Overlay2).Background(Surface0).Padding(0, 2)

	ActiveTab = lipgloss.NewStyle().
			Foreground(Mauve).Background(Mantle).Bold(true).Padding(0, 2).
			Border(lipgloss.Border{
			Top: "─", Bottom: " ", Left: "│", Right: "│",
			TopLeft: "╭", TopRight: "╮", BottomLeft: "┘", BottomRight: "└",
		}).BorderForeground(Mauve)

	// ── HTTP Status badges ────────────────────────────────────
	StatusOK    = lipgloss.NewStyle().Foreground(Crust).Background(Green).Bold(true).Padding(0, 1)
	StatusRedir = lipgloss.NewStyle().Foreground(Crust).Background(Yellow).Bold(true).Padding(0, 1)
	StatusErr   = lipgloss.NewStyle().Foreground(Crust).Background(Red).Bold(true).Padding(0, 1)

	// ── Status bar ───────────────────────────────────────────
	KeyBadge = lipgloss.NewStyle().Foreground(Crust).Background(Surface1).Padding(0, 1)
	KeyDesc  = lipgloss.NewStyle().Foreground(Overlay2)

	// ── Misc ─────────────────────────────────────────────────
	Meta    = lipgloss.NewStyle().Foreground(Overlay0).Italic(true)
	Divider = lipgloss.NewStyle().Foreground(Surface1)

	HistoryItem   = lipgloss.NewStyle().Foreground(Subtext1).Background(Mantle).Padding(0, 1)
	HistoryActive = lipgloss.NewStyle().Foreground(Crust).Background(Mauve).Bold(true).Padding(0, 1)
)
