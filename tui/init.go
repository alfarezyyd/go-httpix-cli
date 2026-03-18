package tui

import (
	"go-httpix-cli/config"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func New() tea.Model {
	return Model{
		MethodIdx:    0,
		URLInput:     newURLInput(),
		BodyInput:    newBodyInput(),
		HeadersInput: newHeadersInput(),
		ParamsInput:  newParamsInput(),
		Focused:      config.PanelURL,
		BodyTabI:     config.TabBody,
		ResponseVP:   viewport.New(80, 20),
		Spinner:      newSpinner(),
		Renderer:     newRenderer(80),
		Keys:         config.DefaultKeyMap(),
		IsMac:        runtime.GOOS == "darwin",
		HistoryIdx:   -1,
	}
}

// Init fires the initial commands: cursor blink, spinner, first tick.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.Spinner.Tick,
		tickEvery(),
	)
}

func tickEvery() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ── Widget constructors ───────────────────────────────────────

func newURLInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "https://api.example.com/endpoint"
	ti.Focus()
	ti.CharLimit = 2048
	ti.Width = 60
	ti.TextStyle = lipgloss.NewStyle().Foreground(config.Text)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(config.Overlay0)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(config.Mauve)
	ti.Prompt = " "
	return ti
}

func newBodyInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "{\n  \"key\": \"value\"\n}"
	ta.SetWidth(60)
	ta.SetHeight(8)
	ta.ShowLineNumbers = true
	return ta
}

func newHeadersInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Content-Type: application/json\nAuthorization: Bearer <token>"
	ta.SetWidth(60)
	ta.SetHeight(8)
	ta.ShowLineNumbers = false
	return ta
}

func newParamsInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "limit=10\noffset=0\nq=search+term"
	ta.SetWidth(60)
	ta.SetHeight(8)
	ta.ShowLineNumbers = false
	return ta
}

func newSpinner() spinner.Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(config.Mauve)
	return sp
}

func newRenderer(wordWrap int) *glamour.TermRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(wordWrap),
	)
	return r
}
