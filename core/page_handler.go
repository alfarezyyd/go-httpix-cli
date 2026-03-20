package core

import (
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"go-httpix-cli/outbound"
	"go-httpix-cli/utils"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func HandleResize(coreModel Model, msg tea.WindowSizeMsg) Model {
	coreModel.Width = msg.Width
	coreModel.Height = msg.Height
	coreModel.ResponseVP.Width = max(20, coreModel.Width-6)
	coreModel.ResponseVP.Height = max(5, coreModel.Height-20)
	coreModel.Renderer = utils.NewRenderer(coreModel.ResponseVP.Width - 4)
	return coreModel
}

func HandleTick(coreModel Model) (Model, tea.Cmd) {
	coreModel.Tick++
	return coreModel, TickEvery()
}

func HandleSpinnerTick(coreModel Model, msg spinner.TickMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	coreModel.Spinner, cmd = coreModel.Spinner.Update(msg)
	return coreModel, cmd
}

func HandleResponse(coreModel Model, msg entity.ResponseMsg) (Model, tea.Cmd) {
	coreModel.Loading = false
	if msg.Err != nil {
		coreModel.ErrMsg = msg.Err.Error()
		coreModel.Response = nil
		return coreModel, nil
	}
	coreModel.ErrMsg = ""
	coreModel.Response = msg.Response
	coreModel.History = PrependHistory(coreModel.History, entity.HistoryEntry{
		Method:   config.HTTPMethods[coreModel.MethodIdx],
		URL:      coreModel.URLInput.Value(),
		Status:   msg.Response.StatusCode,
		Duration: msg.Response.Duration,
		At:       time.Now(),
	})
	coreModel.ResponseVP.SetContent(RenderResponseBody(coreModel))
	coreModel.ResponseVP.GotoTop()
	return coreModel, nil
}

func HandleKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	pressedKey := coreModel.Keys
	if coreModel.CurrentPage() == config.PageEnv {
		return handleEnvPageKey(coreModel, msg)
	}
	if coreModel.Modal.Active != config.ModalNone {
		return handleModalKey(coreModel, msg)
	}
	if coreModel.CollectionOpen {
		return handleCollectionKey(coreModel, msg)
	}
	switch {
	case key.Matches(msg, pressedKey.Quit):
		return coreModel, tea.Quit

	case key.Matches(msg, pressedKey.Send) && !coreModel.Loading:
		coreModel.Loading = true
		coreModel.ErrMsg = ""
		return coreModel, SendRequest(coreModel)

	case key.Matches(msg, pressedKey.NextMethod):
		coreModel.MethodIdx = (coreModel.MethodIdx + 1) % len(config.HTTPMethods)
		return coreModel, nil

	case key.Matches(msg, pressedKey.PrevMethod):
		coreModel.MethodIdx = (coreModel.MethodIdx - 1 + len(config.HTTPMethods)) % len(config.HTTPMethods)
		return coreModel, nil

	case key.Matches(msg, pressedKey.NextTab):
		coreModel.BodyTabI = (coreModel.BodyTabI + 1) % 3
		return coreModel, nil

	case key.Matches(msg, pressedKey.PrevTab):
		coreModel.BodyTabI = (coreModel.BodyTabI - 1 + 3) % 3
		return coreModel, nil

	case key.Matches(msg, pressedKey.FormatJSON) && coreModel.Focused == config.PanelBody:
		if f := utils.PrettyJSON(coreModel.BodyInput.Value()); f != "" {
			coreModel.BodyInput.SetValue(f)
		}
		return coreModel, nil

	case key.Matches(msg, pressedKey.NextPanel):
		return coreModel.CycleFocus(+1), nil

	case key.Matches(msg, pressedKey.PrevPanel):
		return coreModel.CycleFocus(-1), nil

	case key.Matches(msg, pressedKey.HistoryUp) && coreModel.Focused == config.PanelHistory:
		if coreModel.HistoryIdx > 0 {
			coreModel.HistoryIdx--
			coreModel = coreModel.LoadFromHistory()
		}
		return coreModel, nil

	case key.Matches(msg, pressedKey.HistoryDown) && coreModel.Focused == config.PanelHistory:
		if coreModel.HistoryIdx < len(coreModel.History)-1 {
			coreModel.HistoryIdx++
			coreModel = coreModel.LoadFromHistory()
		}
		return coreModel, nil

	case key.Matches(msg, pressedKey.ClearHistory):
		coreModel.History = nil
		coreModel.HistoryIdx = -1
		return coreModel, nil

	case key.Matches(msg, pressedKey.OpenPanelCollection):

		coreModel.CollectionOpen = !coreModel.CollectionOpen
		if coreModel.CollectionOpen {
			coreModel.Focused = config.PanelCollections
		} else {
			coreModel.Focused = config.PanelURL
		}
		return coreModel, nil

	case key.Matches(msg, pressedKey.SaveRequest):
		coreModel.Modal.Active = config.ModalSaveAs
		coreModel.Modal.Cursor = 0
		coreModel.Modal.ErrMsg = ""
		coreModel.Modal.SaveAsSelectedID = ""
		coreModel.Modal.SaveAsInputFocused = true

		// reset dan focus input
		input := coreModel.Modal.SaveAsNameInput
		input.SetValue("")
		input.Focus()
		coreModel.Modal.SaveAsNameInput = input

		return coreModel, nil

	case key.Matches(msg, pressedKey.OpenEnvPicker):
		coreModel.Modal.Active = config.ModalEnvPicker
		coreModel.Modal.Cursor = 0
		return coreModel, nil

		//case key.Matches(msg, pressedKey.OpenEnvPage):
		//	coreModel.CurrentPage = config.PageEnv
		//	coreModel.EnvPageState.List = coreModel.Envs
		//	coreModel.EnvPageCursor = 0
		//	coreModel.EnvPageFocus = EnvFocusList
		//	coreModel.EnvPageRows = []entity.EnvTableRow{}
		//	coreModel.EnvPageRowCursor = 0
		//	// load rows env pertama kalau ada
		//	if len(coreModel.Envs) > 0 {
		//		coreModel.EnvPageRows = core.envToTableRows(coreModel.Envs[0])
		//	}
		//	return coreModel, nil
	}

	return delegateKey(coreModel, msg)
}

func HandleSaveResult(m Model, msg entity.SaveResultMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.Modal.ErrMsg = "failed to save: " + msg.Err.Error()
		m.Modal.Active = config.ModalSaveAs
		return m, nil
	}

	// reload collections dari disk supaya tree terupdate
	return m, reloadCollectionsCmd()
}

func reloadCollectionsCmd() tea.Cmd {
	return func() tea.Msg {
		collections, err := outbound.LoadAllCollections()
		return entity.CollectionsLoadedMsg{
			Collections: collections,
			Err:         err,
		}
	}
}

// delegateKey forwards a keystroke to the focused widget.
func delegateKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch coreModel.Focused {
	case config.PanelURL:
		coreModel.URLInput, cmd = coreModel.URLInput.Update(msg)

	case config.PanelBody:
		switch coreModel.BodyTabI {
		case config.TabBody:
			coreModel.BodyInput, cmd = coreModel.BodyInput.Update(msg)
		case config.TabHeaders:
			coreModel.HeadersInput, cmd = coreModel.HeadersInput.Update(msg)
		case config.TabParams:
			coreModel.ParamsInput, cmd = coreModel.ParamsInput.Update(msg)
		}
	case config.PanelResponse:
		coreModel.ResponseVP, cmd = coreModel.ResponseVP.Update(msg)
	}
	return coreModel, cmd
}
