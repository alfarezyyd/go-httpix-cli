// Package tui is the presentation layer. It wires the bubbletea
// Model / Update / View loop and delegates rendering to tui/component.
package core

import (
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
)

type TickMsg time.Time

type Model struct {
	// ── Dimensions ───────────────────────────────────────────
	Width, Height int

	// ── Page stack ───────────────────────────────────────────
	// Stack memungkinkan navigasi dinamis: push untuk buka page baru,
	// pop untuk kembali. Index terakhir = page yang aktif.
	PageStack  []config.Page
	ActivePage config.Page

	// ── Request inputs ───────────────────────────────────────
	MethodIdx    int
	URLInput     textinput.Model
	BodyInput    textarea.Model
	HeadersInput textarea.Model
	ParamsInput  textarea.Model

	// ── Focus ────────────────────────────────────────────────
	Focused  config.Panel
	BodyTabI config.BodyTab

	// ── Response ─────────────────────────────────────────────
	ResponseVP viewport.Model
	Response   *entity.Response
	Loading    bool
	Spinner    spinner.Model
	ErrMsg     string

	// ── History ──────────────────────────────────────────────
	History    []entity.HistoryEntry
	HistoryIdx int // -1 = nothing selected

	// ── Collections ──────────────────────────────────────────
	CollectionOpen   bool
	Collections      []entity.Collection
	CollectionTree   []entity.TreeNode
	CollectionCursor int

	// ── Environment ──────────────────────────────────────────
	Envs         []entity.Env
	ActiveEnvIdx int // -1 = tidak ada yang aktif
	ActiveEnv    *entity.Env

	// ── Env page ─────────────────────────────────────────────
	EnvPageState EnvPageState

	// ── Modal ────────────────────────────────────────────────
	Modal ModalState

	// ── Cross-cutting ────────────────────────────────────────
	Renderer *glamour.TermRenderer
	Keys     config.KeyMap
	IsMac    bool
	Tick     int
}

// tui/model.go

// ModalState menyimpan semua state yang berkaitan dengan modal aktif.
type ModalState struct {
	Active   config.ModalType
	Input    textinput.Model
	List     []string
	Cursor   int
	ErrMsg   string
	RenameID string // ID node yang sedang di-rename

	SaveAsNameInput    textinput.Model // input nama request
	SaveAsSelectedID   string          // ID folder tujuan, "" = root
	SaveAsInputFocused bool
}

// EnvPageState menyimpan semua state halaman env manager.
type EnvPageState struct {
	List      []entity.Env
	Cursor    int                  // cursor di panel kiri (list)
	Focus     config.EnvPageFocus  // EnvFocusList atau EnvFocusTable
	Rows      []entity.EnvTableRow // rows di panel kanan
	RowCursor int
	Editing   bool
}

// tui/model.go

// CurrentPage mengembalikan page yang sedang aktif.
// Kalau stack kosong, default ke PageMain.
func (coreModel Model) CurrentPage() config.Page {
	if len(coreModel.PageStack) == 0 {
		return config.PageMain
	}
	return coreModel.PageStack[len(coreModel.PageStack)-1]
}

// PushPage membuka page baru di atas stack.
func (coreModel Model) PushPage(p config.Page) Model {
	coreModel.PageStack = append(coreModel.PageStack, p)
	return coreModel
}

// PopPage kembali ke page sebelumnya.
// Kalau stack sudah kosong atau tinggal satu, tidak melakukan apa-apa.
func (coreModel Model) PopPage() Model {
	if len(coreModel.PageStack) <= 1 {
		coreModel.PageStack = []config.Page{}
		return coreModel
	}
	coreModel.PageStack = coreModel.PageStack[:len(coreModel.PageStack)-1]
	return coreModel
}
