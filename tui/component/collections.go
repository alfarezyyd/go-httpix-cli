package component

import (
	"go-httpix-cli/config"
	"go-httpix-cli/tui/collection"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type CollectionPanelProps struct {
	Width   int
	Height  int
	Focused bool
	Nodes   []collection.TreeNode
	Cursor  int
}

func CollectionPanel(p CollectionPanelProps) string {
	title := config.PanelTitleStyle.Render("◈ Collections")

	var rows []string
	for i, node := range p.Nodes {
		indent := strings.Repeat("  ", node.Depth)

		// icon
		icon := "  " // request
		if node.IsFolder {
			if node.Expanded {
				icon = "▼ "
			} else {
				icon = "▶ "
			}
		}

		line := indent + icon + node.Name

		if i == p.Cursor {
			rows = append(rows, lipgloss.NewStyle().
				Foreground(config.Crust).
				Background(config.Mauve).
				Width(p.Width-4).
				Render(line),
			)
		} else if node.IsFolder {
			rows = append(rows, lipgloss.NewStyle().
				Foreground(config.Lavender).
				Width(p.Width-4).
				Render(line),
			)
		} else {
			rows = append(rows, lipgloss.NewStyle().
				Foreground(config.Text).
				Width(p.Width-4).
				Render(line),
			)
		}
	}

	hint := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("  ↵ open  n new  r rename  d delete  esc close")

	content := lipgloss.JoinVertical(lipgloss.Left,
		append(rows, "\n", hint)...,
	)

	ps := config.PanelStyle
	if p.Focused {
		ps = config.PanelFocusedStyle
	}

	return ps.
		Width(p.Width).
		Height(p.Height-4).
		Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}
