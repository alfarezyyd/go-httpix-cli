// tui/component/modal.go
package component

import (
	"go-httpix-cli/config"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type SaveAsModalProps struct {
	Input       textinput.Model
	Collections []string // existing collections untuk ditampilkan
	ErrMsg      string   // kosong = tidak ada error

}

func SaveAsModal(p SaveAsModalProps) string {
	title := lipgloss.NewStyle().
		Foreground(config.Mauve).Bold(true).
		Render("  Save Request  ")

	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Lavender).
		Padding(0, 1).Width(36).
		Render(p.Input.View())

	// error message — hanya muncul kalau ErrMsg tidak kosong
	errLine := ""
	if p.ErrMsg != "" {
		errLine = lipgloss.NewStyle().
			Foreground(config.Red).
			Render("  ✗ " + p.ErrMsg)
	}

	hint := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("  ↵ confirm   esc cancel  ")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"\n",
		"  Collection name:\n",
		input,
		errLine, // muncul di antara input dan hint
		"\n",
		hint,
	)

	return lipgloss.NewStyle().
		Background(config.Mantle).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Mauve).
		Padding(1, 2).
		Width(44).
		Render(inner)
}

type EnvPickerModalProps struct {
	Envs      []string
	ActiveIdx int // env yang sedang aktif (diberi tanda)
	Cursor    int // posisi highlight saat ini
	ErrMsg    string
}

func EnvPickerModal(p EnvPickerModalProps) string {
	title := lipgloss.NewStyle().
		Foreground(config.Mauve).Bold(true).
		Render("  Switch Environment  ")

	// error message — hanya muncul kalau ErrMsg tidak kosong
	errLine := ""
	if p.ErrMsg != "" {
		errLine = lipgloss.NewStyle().
			Foreground(config.Red).
			Render("  ✗ " + p.ErrMsg)
	}

	var rows []string
	for i, name := range p.Envs {
		// tanda bahwa ini env aktif
		active := "  "
		if i == p.ActiveIdx {
			active = "● "
		}

		row := active + name

		if i == p.Cursor {
			rows = append(rows, lipgloss.NewStyle().
				Foreground(config.Crust).
				Background(config.Mauve).
				Width(34).
				Render(row),
			)
		} else {
			rows = append(rows, lipgloss.NewStyle().
				Foreground(config.Text).
				Width(34).
				Render(row),
			)
		}
	}

	hint := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("  ↑↓ navigate   ↵ select   esc cancel  ")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"\n",
		lipgloss.JoinVertical(lipgloss.Left, rows...),
		"\n", hint,
		errLine, // muncul di antara input dan hint
	)

	return lipgloss.NewStyle().
		Background(config.Mantle).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Mauve).
		Padding(1, 2).
		Width(40).
		Render(inner)
}

// ── New Folder Modal ─────────────────────────────────────────

type NewFolderModalProps struct {
	Input  textinput.Model
	ErrMsg string
}

func NewFolderModal(p NewFolderModalProps) string {
	title := lipgloss.NewStyle().
		Foreground(config.Mauve).Bold(true).
		Render("  New Collection Folder  ")

	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Lavender).
		Padding(0, 1).Width(36).
		Render(p.Input.View())

	errLine := ""
	if p.ErrMsg != "" {
		errLine = lipgloss.NewStyle().
			Foreground(config.Red).
			Render("  ✗ " + p.ErrMsg)
	}

	hint := lipgloss.NewStyle().
		Foreground(config.Overlay0).
		Render("  ↵ create   esc cancel  ")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"\n",
		"  Folder name:\n",
		input,
		errLine,
		"\n",
		hint,
	)

	return lipgloss.NewStyle().
		Background(config.Mantle).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Mauve).
		Padding(1, 2).
		Width(44).
		Render(inner)
}

// ── Rename Modal ─────────────────────────────────────────────

type RenameModalProps struct {
	Input      textinput.Model
	TargetName string // nama lama — ditampilkan sebagai konteks
	ErrMsg     string
}

func RenameModal(p RenameModalProps) string {
	title := lipgloss.NewStyle().
		Foreground(config.Mauve).Bold(true).
		Render("  Rename  ")

	// tampilkan nama lama sebagai konteks
	oldName := lipgloss.NewStyle().
		Foreground(config.Overlay0).
		Render("  « " + p.TargetName + " »")

	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Lavender).
		Padding(0, 1).Width(36).
		Render(p.Input.View())

	errLine := ""
	if p.ErrMsg != "" {
		errLine = lipgloss.NewStyle().
			Foreground(config.Red).
			Render("  ✗ " + p.ErrMsg)
	}

	hint := lipgloss.NewStyle().
		Foreground(config.Overlay0).
		Render("  ↵ confirm   esc cancel  ")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		title,
		oldName,
		"\n",
		"  New name:\n",
		input,
		errLine,
		"\n",
		hint,
	)

	return lipgloss.NewStyle().
		Background(config.Mantle).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Mauve).
		Padding(1, 2).
		Width(44).
		Render(inner)
}
