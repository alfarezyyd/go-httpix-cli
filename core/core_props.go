package core

import (
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"go-httpix-cli/tui/component"
)

// props.go: builder methods that map Model → component Props structs.
// Components are pure functions — they know nothing about Model directly.

func (coreModel Model) TopBarProps() component.TopBarProps {
	return component.TopBarProps{
		Width:   coreModel.Width,
		IsMac:   coreModel.IsMac,
		Loading: coreModel.Loading,
		Tick:    coreModel.Tick,
		Keys:    coreModel.Keys,
	}
}

func (coreModel Model) UrlRowProps(w int) component.URLRowProps {
	return component.URLRowProps{
		Width:    w,
		Method:   config.HTTPMethods[coreModel.MethodIdx],
		URLInput: coreModel.URLInput,
		Focused:  coreModel.Focused == config.PanelURL,
	}
}

func (coreModel Model) RequestPanelProps(w int) component.RequestPanelProps {
	return component.RequestPanelProps{
		Width:        w,
		Focused:      coreModel.Focused == config.PanelBody,
		ActiveTab:    coreModel.BodyTabI, // both are config.BodyTab — no mismatch
		BodyInput:    coreModel.BodyInput,
		HeadersInput: coreModel.HeadersInput,
		ParamsInput:  coreModel.ParamsInput,
		TabLabel:     coreModel.Keys.LabelTab,
	}
}

func (coreModel Model) ResponsePanelProps(w int) component.ResponsePanelProps {
	return component.ResponsePanelProps{
		Width:    w,
		Height:   coreModel.Height,
		Focused:  coreModel.Focused == config.PanelResponse,
		Loading:  coreModel.Loading,
		ErrMsg:   coreModel.ErrMsg,
		Response: coreModel.Response,
		VP:       coreModel.ResponseVP,
		Spinner:  coreModel.Spinner,
		SendKey:  coreModel.Keys.LabelSend,
	}
}

func (coreModel Model) CollectionPanelProps(w int) component.CollectionPanelProps {
	return component.CollectionPanelProps{
		Width:   w,
		Height:  coreModel.Height,
		Focused: coreModel.Focused == config.PanelCollections,
		Nodes:   coreModel.CollectionTree,
		Cursor:  coreModel.CollectionCursor,
	}
}

func (coreModel Model) SidebarProps() component.SidebarProps {
	entries := make([]component.HistoryRow, len(coreModel.History))
	for i, h := range coreModel.History {
		entries[i] = component.HistoryRow{
			Method:   h.Method,
			URL:      h.URL,
			Status:   h.Status,
			Duration: h.Duration,
		}
	}
	return component.SidebarProps{
		Width:     config.SidebarWidth,
		Height:    coreModel.Height,
		Focused:   coreModel.Focused == config.PanelHistory,
		Entries:   entries,
		ActiveIdx: coreModel.HistoryIdx,
	}
}

func (coreModel Model) StatusBarProps() component.StatusBarProps {
	labels := map[config.Panel]string{
		config.PanelURL:      "URL",
		config.PanelBody:     "Request Body",
		config.PanelResponse: "Response",
		config.PanelHistory:  "History",
	}
	return component.StatusBarProps{
		Width:      coreModel.Width,
		FocusLabel: labels[coreModel.Focused],
		Keys:       coreModel.Keys,
	}
}

// tui/props.go

func (coreModel Model) envPageProps() entity.EnvPageProps {
	envNames := make([]string, len(coreModel.EnvPageState.List))
	for i, e := range coreModel.EnvPageState.List {
		envNames[i] = e.Name
	}

	rows := make([]entity.EnvRow, len(coreModel.EnvPageState.Rows))
	for i, r := range coreModel.EnvPageState.Rows {
		rows[i] = entity.EnvRow{
			KeyView:   r.Key.View(),
			ValueView: r.Value.View(),
			KeyVal:    r.Key.Value(),
			ValueVal:  r.Value.Value(),
			Editing:   coreModel.EnvPageState.Editing && i == coreModel.EnvPageState.RowCursor,
		}
	}

	envName := ""
	if coreModel.EnvPageState.Cursor >= 0 && coreModel.EnvPageState.Cursor < len(coreModel.EnvPageState.List) {
		envName = coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor].Name
	}

	return entity.EnvPageProps{
		Width:        coreModel.Width,
		Height:       coreModel.Height,
		Envs:         envNames,
		ListCursor:   coreModel.EnvPageState.Cursor,
		ActiveIdx:    coreModel.ActiveEnvIdx,
		ListFocused:  coreModel.EnvPageState.Focus == config.EnvFocusList,
		EnvName:      envName,
		Rows:         rows,
		RowCursor:    coreModel.EnvPageState.RowCursor,
		TableFocused: coreModel.EnvPageState.Focus == config.EnvFocusTable,
		Editing:      coreModel.EnvPageState.Editing,
	}
}
