// Package tui is the presentation layer. It wires the bubbletea
// Model / Update / View loop and delegates rendering to tui/component.
package tui

import (
	"go-httpix-cli/config"
	"go-httpix-cli/core"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
)

// HistoryEntry records one completed request shown in the sidebar.
type HistoryEntry struct {
	Method   string
	URL      string
	Status   int
	Duration time.Duration
	At       time.Time
}

// tickMsg drives the loading-spinner animation at 100 ms intervals.
type tickMsg time.Time

// responseMsg is delivered by the async HTTP goroutine when a request finishes.
type responseMsg struct {
	resp *core.Response
	err  error
}

// Model holds all application state. It is value-copied on every Update.
type Model struct {
	Width, Height int

	// Request inputs
	MethodIdx    int
	URLInput     textinput.Model
	BodyInput    textarea.Model
	HeadersInput textarea.Model
	ParamsInput  textarea.Model

	// Focus
	Focused  config.Panel
	BodyTabI config.BodyTab

	// Response
	ResponseVP viewport.Model
	Response   *core.Response
	Loading    bool
	Spinner    spinner.Model
	ErrMsg     string

	// History
	History    []HistoryEntry
	HistoryIdx int // -1 = nothing selected

	// Cross-cutting
	Renderer *glamour.TermRenderer
	Keys     config.KeyMap
	IsMac    bool
	Tick     int
}
