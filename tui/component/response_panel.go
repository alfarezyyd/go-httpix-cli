package component

import (
	"fmt"
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// ResponsePanelProps is everything the response panel needs to render.
type ResponsePanelProps struct {
	Width    int
	Height   int
	Focused  bool
	Loading  bool
	ErrMsg   string
	Response *entity.Response
	VP       viewport.Model
	Spinner  spinner.Model
	SendKey  string // platform-specific send shortcut label
}

// ResponsePanel renders the response viewer (status, meta, body).
func ResponsePanel(p ResponsePanelProps) string {
	title := config.PanelTitleStyle.Render("◈ Response")
	body := responseBody(p)

	ps := config.PanelStyle
	if p.Focused {
		ps = config.PanelFocusedStyle
	}

	return ps.Width(p.Width-2).Margin(0, 1, 1, 1).Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, body))
}

func responseBody(p ResponsePanelProps) string {
	switch {
	case p.Loading:
		return lipgloss.NewStyle().Foreground(config.Mauve).Padding(2, 3).
			Render(p.Spinner.View() + "  Sending request…")

	case p.ErrMsg != "":
		return lipgloss.NewStyle().Foreground(config.Red).Padding(1, 2).
			Render("✗ Error\n\n" + p.ErrMsg)

	case p.Response == nil:
		return lipgloss.NewStyle().Foreground(config.Overlay0).Italic(true).Padding(3, 3).
			Render("◌  Send a request to see the response here\n\n" +
				"   Press  " + p.SendKey + "  to fire away")

	default:
		return responseContent(p)
	}
}

func responseContent(p ResponsePanelProps) string {
	r := p.Response

	// ── Status badge ─────────────────────────────────────────
	statusStr := fmt.Sprintf("  %d %s  ", r.StatusCode, r.Status)
	var badge string
	switch {
	case r.StatusCode >= 400:
		badge = config.StatusErr.Render(statusStr)
	case r.StatusCode >= 300:
		badge = config.StatusRedir.Render(statusStr)
	default:
		badge = config.StatusOK.Render(statusStr)
	}

	meta := lipgloss.JoinHorizontal(lipgloss.Center,
		badge, " ",
		config.Meta.Render(fmt.Sprintf(" %s", r.Duration.Round(time.Millisecond))),
		"  ",
		config.Meta.Render(humanSize(r.Size)),
		"  ",
		config.Meta.Render(r.Proto),
	)

	divider := config.Divider.Render(strings.Repeat("─", p.Width-4))

	// ── Viewport ─────────────────────────────────────────────
	p.VP.Width = p.Width - 6
	p.VP.Height = max(5, p.Height-36)

	scrollHint := ""
	if p.VP.TotalLineCount() > p.VP.Height {
		scrollHint = config.Meta.Render(fmt.Sprintf(" ↕ %d%%", int(p.VP.ScrollPercent()*100)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, meta, divider, p.VP.View(), scrollHint)
}

func humanSize(b int) string {
	switch {
	case b >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	case b >= 1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	default:
		return fmt.Sprintf("%d B", b)
	}
}
