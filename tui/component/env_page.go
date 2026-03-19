package component

import (
	"go-httpix-cli/config"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type EnvRow struct {
	KeyView   string // hasil dari textinput.View()
	ValueView string
	KeyVal    string // nilai aktual untuk display mode
	ValueVal  string
	Editing   bool // apakah row ini sedang diedit
}
type EnvPageProps struct {
	Width       int
	Height      int
	Envs        []string // hanya nama, untuk list panel
	ListCursor  int
	ActiveIdx   int
	ListFocused bool

	EnvName      string // nama env yang sedang dibuka di table
	Rows         []EnvRow
	RowCursor    int
	TableFocused bool
	Editing      bool
}

func EnvPage(p EnvPageProps) string {
	listW := 28
	tableW := p.Width - listW - 3

	left := envListPanel(p, listW)
	right := envTablePanel(p, tableW)

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)

	return lipgloss.JoinVertical(lipgloss.Left,
		envPageTopBar(p.Width),
		body,
		envPageStatusBar(p),
	)
}

func envPageTopBar(width int) string {
	logo := config.Logo.Render(" ◈ BLINK ")
	section := config.Version.Render(" Environments ")
	hint := lipgloss.NewStyle().
		Foreground(config.Overlay0).Background(config.Mantle).
		Render(" esc back to main  ")

	gap := strings.Repeat(" ", max(0, width-lipgloss.Width(logo+section+hint)))
	return config.TopBar.Width(width).Render(logo + section + gap + hint)
}

func envListPanel(p EnvPageProps, w int) string {
	title := config.PanelTitleStyle.Render("◈ Environments")

	var rows []string
	for i, env := range p.Envs {
		// tanda aktif
		prefix := "  "
		if i == p.ActiveIdx {
			prefix = "● "
		}
		line := prefix + env

		if i == p.ListCursor {
			rows = append(rows, lipgloss.NewStyle().
				Foreground(config.Crust).Background(config.Mauve).
				Width(w-4).Render(line))
		} else {
			rows = append(rows, lipgloss.NewStyle().
				Foreground(config.Text).Width(w-4).Render(line))
		}
	}

	hint := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("\n  ↵ set active\n  n new  d delete\n  tab → table")

	content := lipgloss.JoinVertical(lipgloss.Left,
		append(rows, hint)...,
	)

	ps := config.PanelStyle
	if p.ListFocused {
		ps = config.PanelFocusedStyle
	}

	return ps.Width(w).Height(p.Height-4).Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}

func envTablePanel(p EnvPageProps, w int) string {
	envName := ""
	if len(p.Envs) > 0 && p.ListCursor < len(p.Envs) {
		envName = p.Envs[p.ListCursor]
	}

	title := config.PanelTitleStyle.Render("◈ " + envName)

	// header tabel
	keyW := 22
	valW := w - keyW - 8
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		lipgloss.NewStyle().Foreground(config.Mauve).Bold(true).Width(keyW).Render("  KEY"),
		lipgloss.NewStyle().Foreground(config.Mauve).Bold(true).Width(valW).Render("VALUE"),
	)
	divider := config.Divider.Render(strings.Repeat("┄", w-4))

	// rows
	var rows []string
	for i, row := range p.Rows {
		var line string
		if p.Editing && i == p.RowCursor {
			// mode edit — tampilkan textinput
			line = lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Width(keyW).Render("  "+row.KeyView),
				lipgloss.NewStyle().Width(valW).Render(row.ValueView),
			)
		} else {
			k := row.KeyView
			v := row.ValueView
			if k == "" {
				k = lipgloss.NewStyle().Foreground(config.Overlay0).Render("(empty)")
			}
			// truncate value kalau terlalu panjang
			if len(v) > valW-2 {
				v = v[:valW-5] + "..."
			}
			line = lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Foreground(config.Text).Width(keyW).Render("  "+k),
				lipgloss.NewStyle().Foreground(config.Subtext1).Width(valW).Render(v),
			)
		}

		if i == p.RowCursor && !p.Editing {
			rows = append(rows, lipgloss.NewStyle().
				Background(config.Surface0).Width(w-4).Render(line))
		} else {
			rows = append(rows, line)
		}
	}

	hint := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("\n  ↵ edit   n new row   d delete   tab ← list")

	content := lipgloss.JoinVertical(lipgloss.Left,
		header, divider,
		lipgloss.JoinVertical(lipgloss.Left, rows...),
		hint,
	)

	ps := config.PanelStyle
	if p.TableFocused {
		ps = config.PanelFocusedStyle
	}

	return ps.Width(w).Height(p.Height-4).Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}

func envPageStatusBar(p EnvPageProps) string {
	// status bar sederhana khusus env page
	focus := "Env List"
	if p.TableFocused {
		focus = "Variables"
		if p.Editing {
			focus = "Editing"
		}
	}

	indicator := lipgloss.NewStyle().
		Foreground(config.Crust).Background(config.Lavender).
		Bold(true).Padding(0, 1).
		Render("  " + focus + "  ")

	hints := lipgloss.NewStyle().Foreground(config.Overlay2).
		Render("  Tab switch panel   ↑↓ navigate   ↵ select/edit   n new   d delete   esc back  ")

	return lipgloss.NewStyle().
		Background(config.Crust).Width(p.Width).Padding(0, 1).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, indicator, hints))
}
