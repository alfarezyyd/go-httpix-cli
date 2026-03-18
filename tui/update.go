package tui

import (
	"fmt"
	"go-httpix-cli/config"
	"go-httpix-cli/core"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Update is the Elm-architecture update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg), nil
	case tickMsg:
		return m.handleTick()
	case spinner.TickMsg:
		return m.handleSpinnerTick(msg)
	case responseMsg:
		return m.handleResponse(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

// ── Message handlers ─────────────────────────────────────────

func (m Model) handleResize(msg tea.WindowSizeMsg) Model {
	m.Width = msg.Width
	m.Height = msg.Height
	m.ResponseVP.Width = max(20, m.Width-6)
	m.ResponseVP.Height = max(5, m.Height-20)
	m.Renderer = newRenderer(m.ResponseVP.Width - 4)
	return m
}

func (m Model) handleTick() (Model, tea.Cmd) {
	m.Tick++
	return m, tickEvery()
}

func (m Model) handleSpinnerTick(msg spinner.TickMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

func (m Model) handleResponse(msg responseMsg) (Model, tea.Cmd) {
	m.Loading = false
	if msg.err != nil {
		m.ErrMsg = msg.err.Error()
		m.Response = nil
		return m, nil
	}
	m.ErrMsg = ""
	m.Response = msg.resp
	m.History = prependHistory(m.History, HistoryEntry{
		Method:   config.HTTPMethods[m.MethodIdx],
		URL:      m.URLInput.Value(),
		Status:   msg.resp.Status,
		Duration: msg.resp.Duration,
		At:       time.Now(),
	})
	m.ResponseVP.SetContent(renderResponseBody(m))
	m.ResponseVP.GotoTop()
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := m.Keys
	fmt.Printf("KEY STRING: %q\n", msg.String())
	switch {
	case key.Matches(msg, k.Quit):
		return m, tea.Quit

	case key.Matches(msg, k.Send) && !m.Loading:
		m.Loading = true
		m.ErrMsg = ""
		return m, sendRequest(m)

	case key.Matches(msg, k.NextMethod):
		m.MethodIdx = (m.MethodIdx + 1) % len(config.HTTPMethods)
		return m, nil

	case key.Matches(msg, k.PrevMethod):
		m.MethodIdx = (m.MethodIdx - 1 + len(config.HTTPMethods)) % len(config.HTTPMethods)
		return m, nil

	case key.Matches(msg, k.NextTab):
		m.BodyTabI = (m.BodyTabI + 1) % 3
		return m, nil

	case key.Matches(msg, k.PrevTab):
		m.BodyTabI = (m.BodyTabI - 1 + 3) % 3
		return m, nil

	case key.Matches(msg, k.FormatJSON) && m.Focused == config.PanelBody:
		if f := prettyJSON(m.BodyInput.Value()); f != "" {
			m.BodyInput.SetValue(f)
		}
		return m, nil

	case key.Matches(msg, k.NextPanel):
		return m.cycleFocus(+1), nil

	case key.Matches(msg, k.PrevPanel):
		return m.cycleFocus(-1), nil

	case key.Matches(msg, k.HistoryUp) && m.Focused == config.PanelHistory:
		if m.HistoryIdx > 0 {
			m.HistoryIdx--
			m = m.loadFromHistory()
		}
		return m, nil

	case key.Matches(msg, k.HistoryDown) && m.Focused == config.PanelHistory:
		if m.HistoryIdx < len(m.History)-1 {
			m.HistoryIdx++
			m = m.loadFromHistory()
		}
		return m, nil

	case key.Matches(msg, k.ClearHistory):
		m.History = nil
		m.HistoryIdx = -1
		return m, nil
	}

	return m.delegateKey(msg)
}

// delegateKey forwards a keystroke to the focused widget.
func (m Model) delegateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.Focused {
	case config.PanelURL:
		m.URLInput, cmd = m.URLInput.Update(msg)
	case config.PanelBody:
		switch m.BodyTabI {
		case config.TabBody:
			m.BodyInput, cmd = m.BodyInput.Update(msg)
		case config.TabHeaders:
			m.HeadersInput, cmd = m.HeadersInput.Update(msg)
		case config.TabParams:
			m.ParamsInput, cmd = m.ParamsInput.Update(msg)
		}
	case config.PanelResponse:
		m.ResponseVP, cmd = m.ResponseVP.Update(msg)
	}
	return m, cmd
}

// ── Focus management ─────────────────────────────────────────

var panelOrder = []config.Panel{
	config.PanelURL,
	config.PanelBody,
	config.PanelResponse,
	config.PanelHistory,
}

func (m Model) cycleFocus(dir int) Model {
	cur := 0
	for i, p := range panelOrder {
		if p == m.Focused {
			cur = i
			break
		}
	}
	m.Focused = panelOrder[(cur+dir+len(panelOrder))%len(panelOrder)]
	return m.syncInputFocus()
}

func (m Model) syncInputFocus() Model {
	m.URLInput.Blur()
	m.BodyInput.Blur()
	m.HeadersInput.Blur()
	m.ParamsInput.Blur()

	switch m.Focused {
	case config.PanelURL:
		m.URLInput.Focus()
	case config.PanelBody:
		switch m.BodyTabI {
		case config.TabBody:
			m.BodyInput.Focus()
		case config.TabHeaders:
			m.HeadersInput.Focus()
		case config.TabParams:
			m.ParamsInput.Focus()
		}
	}
	return m
}

// ── History helpers ──────────────────────────────────────────

const maxHistory = 50

func prependHistory(history []HistoryEntry, e HistoryEntry) []HistoryEntry {
	history = append([]HistoryEntry{e}, history...)
	if len(history) > maxHistory {
		history = history[:maxHistory]
	}
	return history
}

func (m Model) loadFromHistory() Model {
	if m.HistoryIdx < 0 || m.HistoryIdx >= len(m.History) {
		return m
	}
	entry := m.History[m.HistoryIdx]
	m.URLInput.SetValue(entry.URL)
	for i, meth := range config.HTTPMethods {
		if meth == entry.Method {
			m.MethodIdx = i
			break
		}
	}
	return m
}

// ── HTTP command ─────────────────────────────────────────────

func sendRequest(m Model) tea.Cmd {
	req := core.Request{
		Method:  config.HTTPMethods[m.MethodIdx],
		URL:     m.URLInput.Value(),
		Body:    m.BodyInput.Value(),
		Headers: m.HeadersInput.Value(),
		Params:  m.ParamsInput.Value(),
	}
	return func() tea.Msg {
		resp, err := core.Execute(req)
		return responseMsg{resp: resp, err: err}
	}
}

// ── Response renderer ────────────────────────────────────────

func renderResponseBody(m Model) string {
	if m.Response == nil {
		return ""
	}
	if pretty, ok := tryPrettyJSON(m.Response.Body); ok {
		md := "```json\n" + pretty + "\n```\n"
		if rendered, err := m.Renderer.Render(md); err == nil {
			return rendered
		}
		return pretty
	}
	return m.Response.Body
}

// ── Misc helpers ─────────────────────────────────────────────

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// keep glamour import alive (only used via newRenderer in init.go)
var _ = glamour.NewTermRenderer
