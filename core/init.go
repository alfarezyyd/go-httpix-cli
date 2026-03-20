package core

import (
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"go-httpix-cli/outbound"
	"go-httpix-cli/utils"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func New() tea.Model {

	return Model{
		// Request inputs
		MethodIdx:    0,
		URLInput:     utils.NewURLInput(),
		BodyInput:    utils.NewBodyInput(),
		HeadersInput: utils.NewHeadersInput(),
		ParamsInput:  utils.NewParamsInput(),

		// Focus
		Focused:  config.PanelURL,
		BodyTabI: config.TabBody,

		// Response
		ResponseVP: viewport.New(80, 20),
		Spinner:    utils.NewSpinner(),
		Renderer:   utils.NewRenderer(80),
		ErrMsg:     "",
		Loading:    false,
		Response:   nil,

		// History
		HistoryIdx: -1,

		// Cross-cutting
		Keys:  config.DefaultKeyMap(),
		IsMac: runtime.GOOS == "darwin",
		Tick:  0,

		// Modal

		Modal: ModalState{
			Active:          config.ModalNone,
			Input:           utils.NewModalInput(),
			SaveAsNameInput: utils.NewModalInput(),
			List:            []string{},
			Cursor:          0,
		},
		// Collections

		// Environment
		Envs:         []entity.Env{},
		ActiveEnvIdx: -1,

		CollectionOpen:   false,
		Collections:      []entity.Collection{},
		CollectionTree:   []entity.TreeNode{},
		CollectionCursor: 0,
	}
}

// Init fires the initial commands: cursor blink, spinner, first tick.
func (coreModel Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		coreModel.Spinner.Tick,
		TickEvery(),
		SetupCollection(),
	)
}

func TickEvery() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func SetupCollection() tea.Cmd {
	collections, err := outbound.LoadAllCollections()
	return func() tea.Msg {
		if err != nil {
			return err
		}
		return entity.CollectionsLoadedMsg{Collections: collections, Err: err}
	}
}

// ── Widget constructors ───────────────────────────────────────
