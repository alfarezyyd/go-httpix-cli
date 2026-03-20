package tui

import (
	"go-httpix-cli/config"
	"go-httpix-cli/core"
	"go-httpix-cli/tui/view"

	"github.com/charmbracelet/lipgloss"
)

type App struct {
	core.Model
}

func (coreModel App) View() string {
	if coreModel.Width == 0 {
		return "Loading…"
	}

	// render tampilan normal seperti biasa
	baseView := view.RenderBase(coreModel)

	switch {
	case coreModel.Modal.Active != config.ModalNone:
		overlay := view.RenderModal(coreModel)
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
			return view.RenderViewWithCollections(coreModel)
		}
		return baseView

	}
}
