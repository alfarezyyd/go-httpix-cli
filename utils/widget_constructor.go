package utils

import (
	"go-httpix-cli/config"
	"go-httpix-cli/entity"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func NewURLInput() textinput.Model {
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

func NewBodyInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "{\n  \"key\": \"value\"\n}"
	ta.SetWidth(60)
	ta.SetHeight(8)
	ta.ShowLineNumbers = true
	return ta
}

func NewHeadersInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Content-Type: application/json\nAuthorization: Bearer <token>"
	ta.SetWidth(60)
	ta.SetHeight(8)
	ta.ShowLineNumbers = false
	return ta
}

func NewParamsInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "limit=10\noffset=0\nq=search+term"
	ta.SetWidth(60)
	ta.SetHeight(8)
	ta.ShowLineNumbers = false
	return ta
}

func NewSpinner() spinner.Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(config.Mauve)
	return sp
}

func NewRenderer(wordWrap int) *glamour.TermRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(wordWrap),
	)
	return r
}

func NewModalInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "e.g. my-collection"
	ti.CharLimit = 64
	ti.Width = 34
	return ti
}

func NewEnvTableRow() entity.EnvTableRow {
	k := textinput.New()
	k.Placeholder = "KEY"
	k.Width = 20

	v := textinput.New()
	v.Placeholder = "value"
	v.Width = 30

	return entity.EnvTableRow{Key: k, Value: v}
}
