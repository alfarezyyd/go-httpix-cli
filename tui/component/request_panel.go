package component

import (
	"go-httpix-cli/config"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RequestPanelProps is everything the request panel needs to render.
type RequestPanelProps struct {
	Width        int
	Focused      bool
	ActiveTab    config.BodyTab // shared type from config — no duplication
	BodyInput    textarea.Model
	HeadersInput textarea.Model
	ParamsInput  textarea.Model
	TabLabel     string
}

// RequestPanel renders the tabbed request editor (body / headers / params).
func RequestPanel(p RequestPanelProps) string {
	const inputH = 10

	tabs := renderTabs(p.ActiveTab, p.TabLabel)

	var content string
	switch p.ActiveTab {
	case config.TabBody:
		p.BodyInput.SetWidth(p.Width - 4)
		p.BodyInput.SetHeight(inputH)
		content = p.BodyInput.View()
	case config.TabHeaders:
		p.HeadersInput.SetWidth(p.Width - 4)
		p.HeadersInput.SetHeight(inputH)
		content = p.HeadersInput.View()
	case config.TabParams:
		p.ParamsInput.SetWidth(p.Width - 4)
		p.ParamsInput.SetHeight(inputH)
		content = p.ParamsInput.View()
	}

	ps := config.PanelStyle
	if p.Focused {
		ps = config.PanelFocusedStyle
	}

	return ps.Width(p.Width-2).Margin(0, 1).Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, tabs, content))
}

func renderTabs(active config.BodyTab, hint string) string {
	names := []string{"  Body  ", " Headers ", " Params "}
	var parts []string
	for i, name := range names {
		if config.BodyTab(i) == active {
			parts = append(parts, config.ActiveTab.Render(name))
		} else {
			parts = append(parts, config.Tab.Render(name))
		}
	}
	hintStr := lipgloss.NewStyle().Foreground(config.Overlay0).Render("  " + hint + " switch tabs")
	row := lipgloss.JoinHorizontal(lipgloss.Bottom, parts...)
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, hintStr) + "\n"
}
