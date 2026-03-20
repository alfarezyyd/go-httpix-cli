package props

import (
	"go-httpix-cli/entity"

	"github.com/charmbracelet/bubbles/textinput"
)

// props/modal_props.go

type SaveAsModalProps struct {
	// input nama request
	Input textinput.Model

	// tree collection untuk dipilih
	Tree       []entity.TreeNode
	Cursor     int    // posisi cursor di tree
	SelectedID string // ID folder yang dipilih, "" = root

	ErrMsg string
}
type EnvPickerModalProps struct {
	Envs      []entity.Env
	ActiveIdx int // env yang sedang aktif (diberi tanda)
	Cursor    int // posisi highlight saat ini
	ErrMsg    string
}

type NewFolderModalProps struct {
	Input  textinput.Model
	ErrMsg string
}

type RenameModalProps struct {
	Input      textinput.Model
	TargetName string // nama lama — ditampilkan sebagai konteks
	ErrMsg     string
}
