package component

import (
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"go-httpix-cli/tui/props"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// tui/component/modal.go

func SaveAsModal(p props.SaveAsModalProps) string {
	title := lipgloss.NewStyle().
		Foreground(config.Mauve).Bold(true).
		Render("  Save Request  ")

	// ── input nama request ────────────────────────────────
	nameLabel := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("  Request name:")
	nameInput := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Lavender).
		Padding(0, 1).Width(36).
		Render(p.Input.View())

	// ── divider ───────────────────────────────────────────
	divider := lipgloss.NewStyle().Foreground(config.Surface1).
		Render("  " + strings.Repeat("─", 36))

	// ── tree collection ───────────────────────────────────
	destLabel := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("  Save to:")

	treeRows := renderSaveAsTree(p.Tree, p.Cursor, p.SelectedID)

	// ── error ─────────────────────────────────────────────
	errLine := ""
	if p.ErrMsg != "" {
		errLine = lipgloss.NewStyle().
			Foreground(config.Red).
			Render("  ✗ " + p.ErrMsg)
	}

	// ── hint ──────────────────────────────────────────────
	hint := lipgloss.NewStyle().Foreground(config.Overlay0).
		Render("  ↑↓ navigate   ↵ select folder   ⌃S save   esc cancel")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"\n",
		nameLabel,
		nameInput,
		"\n",
		divider,
		"\n",
		destLabel,
		treeRows,
		errLine,
		"\n",
		hint,
	)

	return lipgloss.NewStyle().
		Background(config.Mantle).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.Mauve).
		Padding(1, 2).
		Width(48).
		Render(inner)
}

func renderSaveAsTree(nodes []entity.TreeNode, cursor int, selectedID string) string {
	if len(nodes) == 0 {
		return lipgloss.NewStyle().
			Foreground(config.Overlay0).Italic(true).
			Render("  (no collections yet — press N to create one)")
	}

	// tambah opsi root di paling atas
	rootLine := renderSaveAsRow(
		"  ⊙ / (root)",
		cursor == 0,
		selectedID == "",
	)

	var rows []string
	rows = append(rows, rootLine)

	// offset 1 karena index 0 = root
	for i, node := range nodes {
		// hanya tampilkan folder (IsFolder == true)
		// request tidak bisa jadi tujuan simpan
		if !node.IsFolder {
			continue
		}

		indent := strings.Repeat("  ", node.Depth+1)

		icon := "▶ "
		if node.Expanded {
			icon = "▼ "
		}

		label := indent + icon + node.Name
		rows = append(rows, renderSaveAsRow(
			label,
			cursor == i+1, // +1 karena root di index 0
			selectedID == node.ID,
		))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderSaveAsRow(label string, isCursor bool, isSelected bool) string {
	// tanda bahwa ini yang sedang dipilih (confirmed)
	if isSelected {
		label = label + " ✓"
	}

	switch {
	case isCursor && isSelected:
		return lipgloss.NewStyle().
			Foreground(config.Crust).
			Background(config.Green).
			Bold(true).
			Width(40).
			Render(label)
	case isCursor:
		return lipgloss.NewStyle().
			Foreground(config.Crust).
			Background(config.Mauve).
			Width(40).
			Render(label)
	case isSelected:
		return lipgloss.NewStyle().
			Foreground(config.Green).
			Bold(true).
			Width(40).
			Render(label)
	default:
		return lipgloss.NewStyle().
			Foreground(config.Text).
			Width(40).
			Render(label)
	}
}

func EnvPickerModal(p props.EnvPickerModalProps) string {
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
	for i, env := range p.Envs {
		// tanda bahwa ini env aktif
		active := "  "
		if i == p.ActiveIdx {
			active = "● "
		}

		row := active + env.Name

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

func NewFolderModal(p props.NewFolderModalProps) string {
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

func RenameModal(p props.RenameModalProps) string {
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
