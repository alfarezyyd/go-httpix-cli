package entity

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Env struct {
	Name string
	Vars map[string]string
}

type EnvTableRow struct {
	Key   textinput.Model
	Value textinput.Model
}

type EnvRow struct {
	KeyView   string // hasil dari textinput.View()
	ValueView string
	KeyVal    string // nilai aktual untuk display mode
	ValueVal  string
	Editing   bool // apakah row ini sedang diedit
}
type EnvPageProps struct {
	Width       int
	Height      int
	Envs        []string // hanya nama, untuk list panel
	ListCursor  int
	ActiveIdx   int
	ListFocused bool

	EnvName      string // nama env yang sedang dibuka di table
	Rows         []EnvRow
	RowCursor    int
	TableFocused bool
	Editing      bool
}
