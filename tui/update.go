package tui

import (
	"fmt"
	"go-httpix-cli/config"
	"go-httpix-cli/core"
	"go-httpix-cli/entity"
	"go-httpix-cli/tui/collection"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/google/uuid"
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
	case saveResultMsg:
		if msg.err != nil {
			m.ModalErrMsg = "failed to save: " + msg.err.Error()
			m.ActiveModal = config.ModalSaveAs // buka lagi modal dengan error
		}
		return m, nil
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
	if m.CurrentPage == config.PageEnv {
		return m.handleEnvPageKey(msg)
	}
	if m.ActiveModal != config.ModalNone {
		return m.handleModalKey(msg)
	}
	if m.CollectionOpen {
		return m.handleCollectionKey(msg)
	}
	fmt.Printf("%q", msg.String())
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

	case key.Matches(msg, k.OpenPanelCollection):

		m.CollectionOpen = !m.CollectionOpen
		if m.CollectionOpen {
			m.Focused = config.PanelCollections
		} else {
			m.Focused = config.PanelURL
		}
		return m, nil

	case key.Matches(msg, k.SaveRequest):
		m.ActiveModal = config.ModalSaveAs
		m.ModalInput.SetValue("")
		m.ModalInput.Focus()
		core.Logger.Println("enter ditekan, name:", msg.String())

		return m, nil

	case key.Matches(msg, k.OpenEnvPicker):
		m.ActiveModal = config.ModalEnvPicker
		m.ModalCursor = 0
		return m, nil

	case key.Matches(msg, k.OpenEnvPage):
		m.CurrentPage = config.PageEnv
		m.EnvPageList = m.Envs
		m.EnvPageCursor = 0
		m.EnvPageFocus = EnvFocusList
		m.EnvPageRows = []entity.EnvTableRow{}
		m.EnvPageRowCursor = 0
		// load rows env pertama kalau ada
		if len(m.Envs) > 0 {
			m.EnvPageRows = envToTableRows(m.Envs[0])
		}
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

func (m Model) handleModalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.ActiveModal {
	case config.ModalSaveAs:
		return m.handleSaveAsKey(msg)
	case config.ModalEnvPicker:
		return m.handleEnvPickerKey(msg)
	case config.ModalNewFolder: // ← tambah
		return m.handleNewFolderKey(msg)
	case config.ModalRename: // ← tambah
		return m.handleRenameKey(msg)
	}
	return m, nil
}

func (m Model) handleSaveAsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.ActiveModal = config.ModalNone
		m.ModalInput.Blur()
		return m, nil

	case "enter":
		name := strings.TrimSpace(m.ModalInput.Value())
		if name == "" {
			m.ModalErrMsg = "collection name cannot be empty"
			return m, nil // tidak tutup modal, tampilkan error
		}

		if m.URLInput.Value() == "" {
			m.ModalErrMsg = "url input cannot be empty"
			return m, nil // tidak tutup modal, tampilkan error
		}
		m.ActiveModal = config.ModalNone
		m.ModalErrMsg = ""
		m.ModalInput.Blur()
		return m, saveRequestCmd(m, name)

	default:
		// clear error saat user mulai mengetik lagi
		m.ModalErrMsg = ""
		var cmd tea.Cmd
		m.ModalInput, cmd = m.ModalInput.Update(msg)
		return m, cmd
	}
}

func (m Model) handleEnvPickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.ActiveModal = config.ModalNone
		return m, nil

	case "enter":
		// set env aktif sesuai m.ModalCursor
		if len(m.Envs) == 0 {
			m.ModalErrMsg = "Environment variables is empty, can't pick"
			return m, nil
		}
		selected := m.Envs[m.ModalCursor] // ambil value dulu

		m.ActiveEnv = &selected // baru ambil alamatnya
		m.ActiveEnvIdx = m.ModalCursor
		m.ActiveModal = config.ModalNone

		return m, nil

	case "up", "k":
		if m.ModalCursor > 0 {
			m.ModalCursor--
		}
		return m, nil

	case "down", "j":
		if m.ModalCursor < len(m.ModalList)-1 {
			m.ModalCursor++
		}
		return m, nil
	}
	return m, nil
}

func saveRequestCmd(m Model, collectionName string) tea.Cmd {
	return func() tea.Msg {

		return saveResultMsg{err: nil}
	}
}

// ── Misc helpers ─────────────────────────────────────────────

// keep glamour import alive (only used via newRenderer in init.go)
var _ = glamour.NewTermRenderer

// update.go — handler khusus collection panel

func (m Model) handleCollectionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := m.Keys
	switch msg.String() {
	case "up", "k":
		if m.CollectionCursor > 0 {
			m.CollectionCursor--
		}

	case "down", "j":
		if m.CollectionCursor < len(m.CollectionTree)-1 {
			m.CollectionCursor++
		}

	case "enter", " ":
		node := m.CollectionTree[m.CollectionCursor]
		if node.IsFolder {
			// toggle expand/collapse
			m = m.toggleExpand(m.CollectionCursor)
		} else {
			// load request ke form
			m = m.loadRequest(node.Data)
			m.CollectionOpen = false // tutup panel, fokus ke form
		}

	case "n":
		// buka modal → new request di dalam collection ini
		m.ActiveModal = config.ModalSaveAs

	case "N":
		// buka modal new folder
		m.ActiveModal = config.ModalNewFolder
		m.ModalInput = newModalInput()
		m.ModalInput.Placeholder = "e.g. Auth"
		m.ModalInput.Focus()
		return m, nil

	case "d":
		// hapus node yang dipilih
		m = m.deleteNode(m.CollectionCursor)

	case "r":
		// buka modal rename — isi input dengan nama lama
		if len(m.CollectionTree) == 0 {
			return m, nil
		}
		node := m.CollectionTree[m.CollectionCursor]
		m.ActiveModal = config.ModalRename
		m.RenameTargetID = node.ID
		m.ModalInput = newModalInput()
		m.ModalInput.SetValue(node.Name) // pre-fill nama lama
		m.ModalInput.Focus()
		return m, nil
	}

	switch {
	case key.Matches(msg, k.OpenPanelCollection):
		m.CollectionOpen = false
		m.Focused = config.PanelURL
	}

	return m, nil
}

func (m Model) toggleExpand(cursor int) Model {
	// cari collection di nested structure berdasarkan ID
	// toggle Expanded-nya
	// rebuild CollectionTree dengan Flatten()
	node := m.CollectionTree[cursor]
	m.Collections = toggleInTree(m.Collections, node.ID)
	m.CollectionTree = collection.Flatten(m.Collections)
	return m
}

func (m Model) loadRequest(req *collection.Request) Model {
	for i, meth := range config.HTTPMethods {
		if meth == req.Method {
			m.MethodIdx = i
			break
		}
	}
	m.URLInput.SetValue(req.URL)
	m.HeadersInput.SetValue(req.Headers)
	m.BodyInput.SetValue(req.Body)
	m.ParamsInput.SetValue(req.Params)
	return m
}

func (m Model) deleteNode(cursor int) Model {
	if len(m.CollectionTree) == 0 {
		return m
	}
	node := m.CollectionTree[cursor]
	m.Collections = deleteFromTree(m.Collections, node.ID)
	m.CollectionTree = collection.Flatten(m.Collections)

	// jaga cursor tetap dalam batas
	if m.CollectionCursor >= len(m.CollectionTree) {
		m.CollectionCursor = len(m.CollectionTree) - 1
	}
	if m.CollectionCursor < 0 {
		m.CollectionCursor = 0
	}
	return m
}

func deleteFromTree(collections []collection.Collection, id string) []collection.Collection {
	result := []collection.Collection{}
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

func deleteRequestFromList(requests []collection.Request, id string) []collection.Request {
	result := []collection.Request{}
	for _, r := range requests {
		if r.ID != id {
			result = append(result, r)
		}
	}
	return result
}

func toggleInTree(collections []collection.Collection, id string) []collection.Collection {
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

func (m Model) handleNewFolderKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	fmt.Printf("%q", msg.String())

	switch msg.String() {
	case "esc":
		m.ActiveModal = config.ModalNone
		m.ModalErrMsg = ""
		m.ModalInput.Blur()
		return m, nil

	case "enter":
		name := strings.TrimSpace(m.ModalInput.Value())
		if name == "" {
			m.ModalErrMsg = "folder name cannot be empty"
			return m, nil
		}

		// buat collection baru sebagai root folder
		newFolder := collection.Collection{
			ID:       uuid.New().String(),
			Name:     name,
			Children: []collection.Collection{},
			Requests: []collection.Request{},
			Expanded: false,
		}

		m.Collections = append(m.Collections, newFolder)
		m.CollectionTree = collection.Flatten(m.Collections)
		m.ActiveModal = config.ModalNone
		m.ModalErrMsg = ""
		m.ModalInput.SetValue("")
		m.ModalInput.Blur()
		return m, saveCollectionsCmd(newFolder)

	default:
		tea.Println("enter ditekan, name:", msg.String())
		m.ModalErrMsg = ""
		var cmd tea.Cmd
		m.ModalInput, cmd = m.ModalInput.Update(msg)
		core.Logger.Println(msg)
		return m, cmd
	}

}

func (m Model) handleRenameKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.ActiveModal = config.ModalNone
		m.ModalErrMsg = ""
		m.RenameTargetID = ""
		m.ModalInput.Blur()
		return m, nil

	case "enter":
		name := strings.TrimSpace(m.ModalInput.Value())
		if name == "" {
			m.ModalErrMsg = "name cannot be empty"
			return m, nil
		}

		m.Collections = renameInTree(m.Collections, m.RenameTargetID, name)
		m.CollectionTree = collection.Flatten(m.Collections)
		m.ActiveModal = config.ModalNone
		root := findRootCollection(m.Collections, m.RenameTargetID)

		m.ModalErrMsg = ""
		m.RenameTargetID = ""
		m.ModalInput.SetValue("")
		m.ModalInput.Blur()
		return m, saveCollectionsCmd(root)

	default:
		m.ModalErrMsg = ""
		var cmd tea.Cmd
		m.ModalInput, cmd = m.ModalInput.Update(msg)
		return m, cmd
	}
}

// renameInTree mencari node dengan ID tertentu lalu mengubah namanya
func renameInTree(collections []collection.Collection, id, newName string) []collection.Collection {
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
		collections[i].Children = renameInTree(c.Children, id, newName)
	}
	return collections
}

func saveCollectionsCmd(collectionEntity collection.Collection) tea.Cmd {
	return func() tea.Msg {
		err := core.SaveFile(collectionEntity.Name, collectionEntity)
		return saveResultMsg{
			err: err,
		}
	}
}

// findRootCollection mencari root collection mana yang mengandung ID tersebut
func findRootCollection(collections []collection.Collection, id string) collection.Collection {
	for i, c := range collections {
		if containsID(c, id) {
			return collections[i]
		}
	}
	return collection.Collection{}
}

// containsID cek apakah collection ini atau anaknya mengandung ID
func containsID(c collection.Collection, id string) bool {
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

func (m Model) handleEnvPageKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// kalau sedang edit cell, tangkap semua input
	if m.EnvPageEditing {
		return m.handleEnvCellEdit(msg)
	}

	switch m.EnvPageFocus {
	case EnvFocusList:
		return m.handleEnvListPanel(msg)
	case EnvFocusTable:
		return m.handleEnvTablePanel(msg)
	}
	return m, nil
}

func (m Model) handleEnvListPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CurrentPage = config.PageMain
		return m, nil

	case "tab":
		// pindah ke panel kanan
		m.EnvPageFocus = EnvFocusTable
		return m, nil

	case "up", "k":
		if m.EnvPageCursor > 0 {
			m.EnvPageCursor--
			m.EnvPageRows = envToTableRows(m.EnvPageList[m.EnvPageCursor])
			m.EnvPageRowCursor = 0
		}
		return m, nil

	case "down", "j":
		if m.EnvPageCursor < len(m.EnvPageList)-1 {
			m.EnvPageCursor++
			m.EnvPageRows = envToTableRows(m.EnvPageList[m.EnvPageCursor])
			m.EnvPageRowCursor = 0
		}
		return m, nil

	case "enter":
		// set sebagai env aktif
		selected := m.EnvPageList[m.EnvPageCursor]
		m.ActiveEnv = &selected
		m.ActiveEnvIdx = m.EnvPageCursor
		return m, nil

	case "n":
		// tambah env baru
		newEnv := entity.Env{
			Name: "new environment",
			Vars: map[string]string{},
		}
		m.EnvPageList = append(m.EnvPageList, newEnv)
		m.EnvPageCursor = len(m.EnvPageList) - 1
		m.EnvPageRows = []entity.EnvTableRow{newEnvTableRow()}
		m.EnvPageFocus = EnvFocusTable // langsung pindah ke tabel
		return m, nil

	case "d":
		// hapus env yang dipilih
		if len(m.EnvPageList) == 0 {
			return m, nil
		}
		m.EnvPageList = append(
			m.EnvPageList[:m.EnvPageCursor],
			m.EnvPageList[m.EnvPageCursor+1:]...,
		)
		if m.EnvPageCursor >= len(m.EnvPageList) {
			m.EnvPageCursor = len(m.EnvPageList) - 1
		}
		if m.EnvPageCursor >= 0 {
			m.EnvPageRows = envToTableRows(m.EnvPageList[m.EnvPageCursor])
		} else {
			m.EnvPageRows = []entity.EnvTableRow{}
		}
		return m, saveEnvsCmd(m.EnvPageList)
	}
	return m, nil
}

func (m Model) handleEnvTablePanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "tab":
		// kembali ke list
		m.EnvPageFocus = EnvFocusList
		return m, nil

	case "up", "k":
		if m.EnvPageRowCursor > 0 {
			m.EnvPageRowCursor--
		}
		return m, nil

	case "down", "j":
		if m.EnvPageRowCursor < len(m.EnvPageRows)-1 {
			m.EnvPageRowCursor++
		}
		return m, nil

	case "enter":
		// mulai edit row yang dipilih
		m.EnvPageEditing = true
		m.EnvPageRows[m.EnvPageRowCursor].Key.Focus()
		return m, nil

	case "n":
		// tambah row baru
		m.EnvPageRows = append(m.EnvPageRows, newEnvTableRow())
		m.EnvPageRowCursor = len(m.EnvPageRows) - 1
		m.EnvPageEditing = true
		m.EnvPageRows[m.EnvPageRowCursor].Key.Focus()
		return m, nil

	case "d":
		// hapus row
		if len(m.EnvPageRows) == 0 {
			return m, nil
		}
		m.EnvPageRows = append(
			m.EnvPageRows[:m.EnvPageRowCursor],
			m.EnvPageRows[m.EnvPageRowCursor+1:]...,
		)
		if m.EnvPageRowCursor >= len(m.EnvPageRows) {
			m.EnvPageRowCursor = len(m.EnvPageRows) - 1
		}
		return m, m.saveCurrentEnv()
	}
	return m, nil
}

func (m Model) handleEnvCellEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	row := &m.EnvPageRows[m.EnvPageRowCursor]

	switch msg.String() {
	case "tab":
		// pindah dari key ke value atau selesai edit
		if row.Key.Focused() {
			row.Key.Blur()
			row.Value.Focus()
		} else {
			row.Value.Blur()
			m.EnvPageEditing = false
			return m, m.saveCurrentEnv()
		}
		return m, nil

	case "esc":
		row.Key.Blur()
		row.Value.Blur()
		m.EnvPageEditing = false
		return m, m.saveCurrentEnv()

	case "enter":
		row.Key.Blur()
		row.Value.Blur()
		m.EnvPageEditing = false
		return m, m.saveCurrentEnv()
	}

	// teruskan ke input yang aktif
	var cmd tea.Cmd
	if row.Key.Focused() {
		m.EnvPageRows[m.EnvPageRowCursor].Key, cmd = row.Key.Update(msg)
	} else {
		m.EnvPageRows[m.EnvPageRowCursor].Value, cmd = row.Value.Update(msg)
	}
	return m, cmd
}

func envToTableRows(e entity.Env) []entity.EnvTableRow {
	rows := make([]entity.EnvTableRow, 0, len(e.Vars))
	for k, v := range e.Vars {
		row := newEnvTableRow()
		row.Key.SetValue(k)
		row.Value.SetValue(v)
		rows = append(rows, row)
	}
	return rows
}

func newEnvTableRow() entity.EnvTableRow {
	k := textinput.New()
	k.Placeholder = "KEY"
	k.Width = 20

	v := textinput.New()
	v.Placeholder = "value"
	v.Width = 30

	return entity.EnvTableRow{Key: k, Value: v}
}

func tableRowsToEnv(name string, rows []entity.EnvTableRow) entity.Env {
	vars := make(map[string]string, len(rows))
	for _, row := range rows {
		k := strings.TrimSpace(row.Key.Value())
		v := strings.TrimSpace(row.Value.Value())
		if k != "" {
			vars[k] = v
		}
	}
	return entity.Env{Name: name, Vars: vars}
}

// saveCurrentEnv menyimpan env yang sedang diedit
func (m Model) saveCurrentEnv() tea.Cmd {
	if m.EnvPageCursor < 0 || m.EnvPageCursor >= len(m.EnvPageList) {
		return nil
	}
	name := m.EnvPageList[m.EnvPageCursor].Name
	updated := tableRowsToEnv(name, m.EnvPageRows)
	m.EnvPageList[m.EnvPageCursor] = updated
	m.Envs = m.EnvPageList // sync ke model utama
	return saveEnvsCmd(m.EnvPageList)
}

func saveEnvsCmd(envPageList []entity.Env) tea.Cmd {
	return func() tea.Msg {

		return saveResultMsg{err: nil}
	}
}
