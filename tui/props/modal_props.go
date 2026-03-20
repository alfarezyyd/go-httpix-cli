package props

import (
	"go-httpix-cli/entity"

	"github.com/charmbracelet/bubbles/textinput"
)

type SaveAsModalProps struct {
	Input       textinput.Model
	Collections []string // existing collections untuk ditampilkan
	ErrMsg      string   // kosong = tidak ada error
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
