package core

import (
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"go-httpix-cli/outbound"
	"go-httpix-cli/utils"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Update is the Elm-architecture update function.
func (coreModel Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return HandleResize(coreModel, msg), nil
	case TickMsg:
		return HandleTick(coreModel)
	case spinner.TickMsg:
		return HandleSpinnerTick(coreModel, msg)
	case entity.ResponseMsg:
		return HandleResponse(coreModel, msg)
	case entity.SaveResultMsg:
		return HandleSaveResult(coreModel, msg)
	case tea.KeyMsg:
		return HandleKey(coreModel, msg)
	}
	return coreModel, nil
}

// ── Focus management ─────────────────────────────────────────

func (coreModel Model) CycleFocus(dir int) Model {
	cur := 0
	for i, p := range config.PanelOrder {
		if p == coreModel.Focused {
			cur = i
			break
		}
	}
	coreModel.Focused = config.PanelOrder[(cur+dir+len(config.PanelOrder))%len(config.PanelOrder)]
	return coreModel.syncInputFocus()
}

func (coreModel Model) syncInputFocus() Model {
	coreModel.URLInput.Blur()
	coreModel.BodyInput.Blur()
	coreModel.HeadersInput.Blur()
	coreModel.ParamsInput.Blur()

	switch coreModel.Focused {
	case config.PanelURL:
		coreModel.URLInput.Focus()
	case config.PanelBody:
		switch coreModel.BodyTabI {
		case config.TabBody:
			coreModel.BodyInput.Focus()
		case config.TabHeaders:
			coreModel.HeadersInput.Focus()
		case config.TabParams:
			coreModel.ParamsInput.Focus()
		}
	}
	return coreModel
}

// ── History helpers ──────────────────────────────────────────

const maxHistory = 50

func PrependHistory(history []entity.HistoryEntry, e entity.HistoryEntry) []entity.HistoryEntry {
	history = append([]entity.HistoryEntry{e}, history...)
	if len(history) > maxHistory {
		history = history[:maxHistory]
	}
	return history
}

func (coreModel Model) LoadFromHistory() Model {
	if coreModel.HistoryIdx < 0 || coreModel.HistoryIdx >= len(coreModel.History) {
		return coreModel
	}
	entry := coreModel.History[coreModel.HistoryIdx]
	coreModel.URLInput.SetValue(entry.URL)
	for i, meth := range config.HTTPMethods {
		if meth == entry.Method {
			coreModel.MethodIdx = i
			break
		}
	}
	return coreModel
}

// ── HTTP command ─────────────────────────────────────────────

func SendRequest(m Model) tea.Cmd {
	req := entity.Request{
		Method:  config.HTTPMethods[m.MethodIdx],
		URL:     m.URLInput.Value(),
		Body:    m.BodyInput.Value(),
		Headers: m.HeadersInput.Value(),
		Params:  m.ParamsInput.Value(),
	}
	return func() tea.Msg {
		resp, err := outbound.GetHttpClient().Execute(&req)
		return entity.ResponseMsg{Response: resp, Err: err}
	}
}

func RenderResponseBody(m Model) string {
	if m.Response == nil {
		return ""
	}
	if pretty, ok := utils.TryPrettyJSON(m.Response.Body); ok {
		md := "```json\n" + pretty + "\n```\n"
		if rendered, err := m.Renderer.Render(md); err == nil {
			return rendered
		}
		return pretty
	}
	return m.Response.Body
}

func SaveRequestCmd(m Model, collectionName string) tea.Cmd {
	return func() tea.Msg {

		return entity.SaveResultMsg{Err: nil}
	}
}

// keep glamour import alive (only used via newRenderer in init.go)
var _ = glamour.NewTermRenderer

// update.go — handler khusus collection panel

func (coreModel Model) ToggleExpand(cursor int) Model {
	// cari collection di nested structure berdasarkan ID
	// toggle Expanded-nya
	// rebuild CollectionTree dengan Flatten()
	node := coreModel.CollectionTree[cursor]
	coreModel.Collections = toggleInTree(coreModel.Collections, node.ID)
	coreModel.CollectionTree = utils.Flatten(coreModel.Collections)
	return coreModel
}

func (coreModel Model) LoadRequest(req *entity.Request) Model {
	for i, meth := range config.HTTPMethods {
		if meth == req.Method {
			coreModel.MethodIdx = i
			break
		}
	}
	coreModel.URLInput.SetValue(req.URL)
	coreModel.HeadersInput.SetValue(req.Headers)
	coreModel.BodyInput.SetValue(req.Body)
	coreModel.ParamsInput.SetValue(req.Params)
	return coreModel
}

func (coreModel Model) DeleteNode(cursor int) Model {
	if len(coreModel.CollectionTree) == 0 {
		return coreModel
	}
	node := coreModel.CollectionTree[cursor]
	coreModel.Collections = deleteFromTree(coreModel.Collections, node.ID)
	coreModel.CollectionTree = utils.Flatten(coreModel.Collections)

	// jaga cursor tetap dalam batas
	if coreModel.CollectionCursor >= len(coreModel.CollectionTree) {
		coreModel.CollectionCursor = len(coreModel.CollectionTree) - 1
	}
	if coreModel.CollectionCursor < 0 {
		coreModel.CollectionCursor = 0
	}
	return coreModel
}

func deleteFromTree(collections []entity.Collection, id string) []entity.Collection {
	result := []entity.Collection{}
	for _, c := range collections {
		if c.ID == id {
			continue // skip — ini yang dihapus
		}
		// cek juga di children dan requests
		c.Children = deleteFromTree(c.Children, id)
		c.Requests = deleteRequestFromList(c.Requests, id)
		result = append(result, c)
	}
	return result
}

func deleteRequestFromList(requests []entity.Request, id string) []entity.Request {
	result := []entity.Request{}
	for _, r := range requests {
		if r.ID != id {
			result = append(result, r)
		}
	}
	return result
}

func toggleInTree(collections []entity.Collection, id string) []entity.Collection {
	for i, c := range collections {
		if c.ID == id {
			collections[i].Expanded = !collections[i].Expanded
			return collections
		}
		// cari di children secara rekursif
		collections[i].Children = toggleInTree(c.Children, id)
	}
	return collections
}

// renameInTree mencari node dengan ID tertentu lalu mengubah namanya
func RenameInTree(collections []entity.Collection, id, newName string) []entity.Collection {
	for i, c := range collections {
		if c.ID == id {
			collections[i].Name = newName
			return collections
		}
		// cari di request
		for j, r := range c.Requests {
			if r.ID == id {
				collections[i].Requests[j].Name = newName
				return collections
			}
		}
		// rekursif ke children
		collections[i].Children = RenameInTree(c.Children, id, newName)
	}
	return collections
}

func SaveCollectionsCmd(collectionEntity entity.Collection) tea.Cmd {
	return func() tea.Msg {
		err := outbound.SaveFile(collectionEntity.Name, collectionEntity)
		return entity.SaveResultMsg{
			Err: err,
		}
	}
}

// findRootCollection mencari root collection mana yang mengandung ID tersebut
func FindRootCollection(collections []entity.Collection, id string) entity.Collection {
	for i, c := range collections {
		if containsID(c, id) {
			return collections[i]
		}
	}
	return entity.Collection{}
}

// containsID cek apakah collection ini atau anaknya mengandung ID
func containsID(c entity.Collection, id string) bool {
	if c.ID == id {
		return true
	}
	for _, r := range c.Requests {
		if r.ID == id {
			return true
		}
	}
	for _, child := range c.Children {
		if containsID(child, id) {
			return true
		}
	}
	return false
}

// saveCurrentEnv menyimpan env yang sedang diedit
func (coreModel Model) SaveCurrentEnv() tea.Cmd {
	if coreModel.EnvPageState.Cursor < 0 || coreModel.EnvPageState.Cursor >= len(coreModel.EnvPageState.List) {
		return nil
	}
	name := coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor].Name
	updated := utils.TableRowsToEnv(name, coreModel.EnvPageState.Rows)
	coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor] = updated
	coreModel.Envs = coreModel.EnvPageState.List // sync ke model utama
	return SaveEnvsCmd(coreModel.EnvPageState.List)
}

func SaveEnvsCmd(envPageList []entity.Env) tea.Cmd {
	return func() tea.Msg {

		return entity.SaveResultMsg{Err: nil}
	}
}
