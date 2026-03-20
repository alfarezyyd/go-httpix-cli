package core

import (
	"go-httpix-cli/config"
	"go-httpix-cli/tui/component"
	"go-httpix-cli/tui/props"
	"go-httpix-cli/tui/view"

	"github.com/charmbracelet/lipgloss"
)

// View is the Elm-architecture view function.
// It composes sub-views into the final terminal string.
func (coreModel Model) View() string {
	if coreModel.Width == 0 {
		return "Loading…"
	}

	// render tampilan normal seperti biasa
	baseView := RenderBase(coreModel)

	switch {
	case coreModel.Modal.Active != config.ModalNone:
		overlay := RenderModal(coreModel)
		return lipgloss.Place(
			coreModel.Width, coreModel.Height, // ukuran canvas
			lipgloss.Center, // posisi horizontal: tengah
			lipgloss.Center, // posisi vertikal: tengah
			overlay,         // konten yang ditimpa
			lipgloss.WithWhitespaceChars("░"),
			lipgloss.WithWhitespaceForeground(config.Surface0),
		)
	case coreModel.CurrentPage() == config.PageEnv:
		return view.EnvPage(coreModel.envPageProps())
	default:
		if coreModel.CollectionOpen {
			return RenderViewWithCollections(coreModel)
		}
		return baseView

	}
}

func RenderBase(coreModel Model) string {
	mainW := max(40, coreModel.Width-config.SidebarWidth-3)

	main := lipgloss.JoinVertical(lipgloss.Left,
		component.URLRow(coreModel.UrlRowProps(mainW)),
		component.RequestPanel(coreModel.RequestPanelProps(mainW)),
		component.ResponsePanel(coreModel.ResponsePanelProps(mainW)),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		component.TopBar(coreModel.TopBarProps()),
		lipgloss.JoinHorizontal(lipgloss.Top, main, " ", component.Sidebar(coreModel.SidebarProps())),
		component.StatusBar(coreModel.StatusBarProps()),
	)
}

func RenderModal(coreModel Model) string {
	switch coreModel.Modal.Active {
	case config.ModalSaveAs:
		return component.SaveAsModal(props.SaveAsModalProps{
			Input:      coreModel.Modal.SaveAsNameInput,
			Tree:       coreModel.CollectionTree,
			Cursor:     coreModel.Modal.Cursor,
			SelectedID: coreModel.Modal.SaveAsSelectedID,
			ErrMsg:     coreModel.Modal.ErrMsg,
		})

	case config.ModalNewFolder:
		return component.NewFolderModal(props.NewFolderModalProps{
			Input:  coreModel.Modal.Input,
			ErrMsg: coreModel.Modal.ErrMsg,
		})

	case config.ModalRename:
		// cari nama lama untuk ditampilkan sebagai konteks
		targetName := ""
		if len(coreModel.CollectionTree) > 0 {
			for _, node := range coreModel.CollectionTree {
				if node.ID == coreModel.Modal.RenameID {
					targetName = node.Name
					break
				}
			}
		}
		return component.RenameModal(props.RenameModalProps{
			Input:      coreModel.Modal.Input,
			TargetName: targetName,
			ErrMsg:     coreModel.Modal.ErrMsg,
		})

	case config.ModalEnvPicker:
		return component.EnvPickerModal(props.EnvPickerModalProps{
			Envs:      coreModel.Envs,
			ActiveIdx: coreModel.ActiveEnvIdx,
			Cursor:    coreModel.Modal.Cursor,
		})
	}
	return ""
}

func RenderViewWithCollections(coreModel Model) string {
	collW := 30
	mainW := max(40, coreModel.Width-collW-config.SidebarWidth-4)

	collPanel := component.CollectionPanel(coreModel.CollectionPanelProps(collW))

	main := lipgloss.JoinVertical(lipgloss.Left,
		component.URLRow(coreModel.UrlRowProps(mainW)),
		component.RequestPanel(coreModel.RequestPanelProps(mainW)),
		component.ResponsePanel(coreModel.ResponsePanelProps(mainW)),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top,
		collPanel,
		" ",
		main,
		" ",
		component.Sidebar(coreModel.SidebarProps()),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		component.TopBar(coreModel.TopBarProps()),
		body,
		component.StatusBar(coreModel.StatusBarProps()),
	)
}
