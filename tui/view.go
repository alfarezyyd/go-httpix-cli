package tui

import (
	"go-httpix-cli/config"
	"go-httpix-cli/tui/component"

	"github.com/charmbracelet/lipgloss"
)

const sidebarWidth = 28

// View is the Elm-architecture view function.
// It composes sub-views into the final terminal string.
func (m Model) View() string {
	if m.Width == 0 {
		return "Loading…"
	}

	// render tampilan normal seperti biasa
	base := m.renderBase()

	// kalau ada modal aktif, timpa di atas base
	if m.ActiveModal != config.ModalNone {
		overlay := m.renderModal()
		return lipgloss.Place(
			m.Width, m.Height, // ukuran canvas
			lipgloss.Center, // posisi horizontal: tengah
			lipgloss.Center, // posisi vertikal: tengah
			overlay,         // konten yang ditimpa
			lipgloss.WithWhitespaceChars("░"),
			lipgloss.WithWhitespaceForeground(config.Surface0),
		)
	}

	if m.CollectionOpen {
		return m.viewWithCollections()
	}

	return base
}

func (m Model) renderBase() string {
	mainW := max(40, m.Width-sidebarWidth-3)

	main := lipgloss.JoinVertical(lipgloss.Left,
		component.URLRow(m.urlRowProps(mainW)),
		component.RequestPanel(m.requestPanelProps(mainW)),
		component.ResponsePanel(m.responsePanelProps(mainW)),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		component.TopBar(m.topBarProps()),
		lipgloss.JoinHorizontal(lipgloss.Top, main, " ", component.Sidebar(m.sidebarProps())),
		component.StatusBar(m.statusBarProps()),
	)
}

func (m Model) renderModal() string {
	switch m.ActiveModal {
	case config.ModalSaveAs:
		return component.SaveAsModal(component.SaveAsModalProps{
			Input:  m.ModalInput,
			ErrMsg: m.ModalErrMsg,
		})

	case config.ModalNewFolder:
		return component.NewFolderModal(component.NewFolderModalProps{
			Input:  m.ModalInput,
			ErrMsg: m.ModalErrMsg,
		})

	case config.ModalRename:
		// cari nama lama untuk ditampilkan sebagai konteks
		targetName := ""
		if len(m.CollectionTree) > 0 {
			for _, node := range m.CollectionTree {
				if node.ID == m.RenameTargetID {
					targetName = node.Name
					break
				}
			}
		}
		return component.RenameModal(component.RenameModalProps{
			Input:      m.ModalInput,
			TargetName: targetName,
			ErrMsg:     m.ModalErrMsg,
		})

	case config.ModalEnvPicker:
		return component.EnvPickerModal(component.EnvPickerModalProps{
			Envs:      m.EnvNames,
			ActiveIdx: m.ActiveEnvIdx,
			Cursor:    m.ModalCursor,
		})
	}
	return ""
}

func (m Model) viewWithCollections() string {
	collW := 30
	mainW := max(40, m.Width-collW-sidebarWidth-4)

	collPanel := component.CollectionPanel(m.collectionPanelProps(collW))

	main := lipgloss.JoinVertical(lipgloss.Left,
		component.URLRow(m.urlRowProps(mainW)),
		component.RequestPanel(m.requestPanelProps(mainW)),
		component.ResponsePanel(m.responsePanelProps(mainW)),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top,
		collPanel,
		" ",
		main,
		" ",
		component.Sidebar(m.sidebarProps()),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		component.TopBar(m.topBarProps()),
		body,
		component.StatusBar(m.statusBarProps()),
	)
}
