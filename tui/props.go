package tui

import (
	"go-httpix-cli/config"
	"go-httpix-cli/tui/component"
)

// props.go: builder methods that map Model → component Props structs.
// Components are pure functions — they know nothing about Model directly.

func (m Model) topBarProps() component.TopBarProps {
	return component.TopBarProps{
		Width:   m.Width,
		IsMac:   m.IsMac,
		Loading: m.Loading,
		Tick:    m.Tick,
		Keys:    m.Keys,
	}
}

func (m Model) urlRowProps(w int) component.URLRowProps {
	return component.URLRowProps{
		Width:    w,
		Method:   config.HTTPMethods[m.MethodIdx],
		URLInput: m.URLInput,
		Focused:  m.Focused == config.PanelURL,
	}
}

func (m Model) requestPanelProps(w int) component.RequestPanelProps {
	return component.RequestPanelProps{
		Width:        w,
		Focused:      m.Focused == config.PanelBody,
		ActiveTab:    m.BodyTabI, // both are config.BodyTab — no mismatch
		BodyInput:    m.BodyInput,
		HeadersInput: m.HeadersInput,
		ParamsInput:  m.ParamsInput,
		TabLabel:     m.Keys.LabelTab,
	}
}

func (m Model) responsePanelProps(w int) component.ResponsePanelProps {
	return component.ResponsePanelProps{
		Width:    w,
		Height:   m.Height,
		Focused:  m.Focused == config.PanelResponse,
		Loading:  m.Loading,
		ErrMsg:   m.ErrMsg,
		Response: m.Response,
		VP:       m.ResponseVP,
		Spinner:  m.Spinner,
		SendKey:  m.Keys.LabelSend,
	}
}

func (m Model) collectionPanelProps(w int) component.CollectionPanelProps {
	return component.CollectionPanelProps{
		Width:   w,
		Height:  m.Height,
		Focused: m.Focused == config.PanelCollections,
		Nodes:   m.CollectionTree,
		Cursor:  m.CollectionCursor,
	}
}

func (m Model) sidebarProps() component.SidebarProps {
	entries := make([]component.HistoryRow, len(m.History))
	for i, h := range m.History {
		entries[i] = component.HistoryRow{
			Method:   h.Method,
			URL:      h.URL,
			Status:   h.Status,
			Duration: h.Duration,
		}
	}
	return component.SidebarProps{
		Width:     sidebarWidth,
		Height:    m.Height,
		Focused:   m.Focused == config.PanelHistory,
		Entries:   entries,
		ActiveIdx: m.HistoryIdx,
	}
}

func (m Model) statusBarProps() component.StatusBarProps {
	labels := map[config.Panel]string{
		config.PanelURL:      "URL",
		config.PanelBody:     "Request Body",
		config.PanelResponse: "Response",
		config.PanelHistory:  "History",
	}
	return component.StatusBarProps{
		Width:      m.Width,
		FocusLabel: labels[m.Focused],
		Keys:       m.Keys,
	}
}

// tui/props.go

func (m Model) envPageProps() component.EnvPageProps {
	envNames := make([]string, len(m.EnvPageList))
	for i, e := range m.EnvPageList {
		envNames[i] = e.Name
	}

	rows := make([]component.EnvRow, len(m.EnvPageRows))
	for i, r := range m.EnvPageRows {
		rows[i] = component.EnvRow{
			KeyView:   r.Key.View(),
			ValueView: r.Value.View(),
			KeyVal:    r.Key.Value(),
			ValueVal:  r.Value.Value(),
			Editing:   m.EnvPageEditing && i == m.EnvPageRowCursor,
		}
	}

	envName := ""
	if m.EnvPageCursor >= 0 && m.EnvPageCursor < len(m.EnvPageList) {
		envName = m.EnvPageList[m.EnvPageCursor].Name
	}

	return component.EnvPageProps{
		Width:        m.Width,
		Height:       m.Height,
		Envs:         envNames,
		ListCursor:   m.EnvPageCursor,
		ActiveIdx:    m.ActiveEnvIdx,
		ListFocused:  m.EnvPageFocus == EnvFocusList,
		EnvName:      envName,
		Rows:         rows,
		RowCursor:    m.EnvPageRowCursor,
		TableFocused: m.EnvPageFocus == EnvFocusTable,
		Editing:      m.EnvPageEditing,
	}
}
