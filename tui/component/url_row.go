package component

import (
	"fmt"
	"go-httpix-cli/config"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// URLRowProps is everything the URL row needs to render.
type URLRowProps struct {
	Width    int
	Method   string
	URLInput textinput.Model
	Focused  bool
}

// URLRow renders the method badge + URL input + SEND button.
func URLRow(p URLRowProps) string {
	col, ok := config.MethodColor[p.Method]
	if !ok {
		col = config.Overlay2
	}
	methBadge := lipgloss.NewStyle().
		Foreground(config.Crust).Background(col).Bold(true).Padding(0, 1).
		Render(fmt.Sprintf("  %s  ", p.Method))

	urlW := max(10, p.Width-lipgloss.Width(methBadge)-12)
	p.URLInput.Width = urlW

	barStyle := config.URLBar
	if p.Focused {
		barStyle = config.URLBarFocused
	}
	urlBox := barStyle.Width(urlW + 2).Render(p.URLInput.View())

	sendBtn := lipgloss.NewStyle().
		Foreground(config.Crust).Background(config.Mauve).Bold(true).Padding(0, 2).
		Render("  SEND  ")

	row := lipgloss.JoinHorizontal(lipgloss.Center, methBadge, " ", urlBox, " ", sendBtn)

	return lipgloss.NewStyle().
		Background(config.Base).Padding(1, 1).Width(p.Width).
		Render(row)
}
