package core

import (
	"fmt"
	"go-httpix-cli/config"
	"go-httpix-cli/entity"
	"go-httpix-cli/utils"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

func handleModalKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch coreModel.Modal.Active {
	case config.ModalSaveAs:
		return handleSaveAsKey(coreModel, msg)
	case config.ModalEnvPicker:
		return handleEnvPickerKey(coreModel, msg)
	case config.ModalNewFolder: // ← tambah
		return handleNewFolderKey(coreModel, msg)
	case config.ModalRename: // ← tambah
		return handleRenameKey(coreModel, msg)
	}
	return coreModel, nil
}

func handleSaveAsKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		coreModel.Modal.Active = config.ModalNone
		coreModel.Modal.Input.Blur()
		return coreModel, nil

	case "enter":
		name := strings.TrimSpace(coreModel.Modal.Input.Value())
		if name == "" {
			coreModel.Modal.ErrMsg = "collection name cannot be empty"
			return coreModel, nil // tidak tutup modal, tampilkan error
		}

		if coreModel.URLInput.Value() == "" {
			coreModel.Modal.ErrMsg = "url input cannot be empty"
			return coreModel, nil // tidak tutup modal, tampilkan error
		}
		coreModel.Modal.Active = config.ModalNone
		coreModel.Modal.ErrMsg = ""
		coreModel.Modal.Input.Blur()
		return coreModel, SaveRequestCmd(coreModel, name)

	default:
		// clear error saat user mulai mengetik lagi
		coreModel.Modal.ErrMsg = ""
		var cmd tea.Cmd
		coreModel.Modal.Input, cmd = coreModel.Modal.Input.Update(msg)
		return coreModel, cmd
	}
}

func handleEnvPickerKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		coreModel.Modal.Active = config.ModalNone
		return coreModel, nil

	case "enter":
		// set env aktif sesuai coreModel.Modal.Cursor
		if len(coreModel.Envs) == 0 {
			coreModel.Modal.ErrMsg = "Environment variables is empty, can't pick"
			return coreModel, nil
		}
		selected := coreModel.Envs[coreModel.Modal.Cursor] // ambil value dulu

		coreModel.ActiveEnv = &selected // baru ambil alamatnya
		coreModel.ActiveEnvIdx = coreModel.Modal.Cursor
		coreModel.Modal.Active = config.ModalNone

		return coreModel, nil

	case "up", "k":
		if coreModel.Modal.Cursor > 0 {
			coreModel.Modal.Cursor--
		}
		return coreModel, nil

	case "down", "j":
		if coreModel.Modal.Cursor < len(coreModel.Modal.List)-1 {
			coreModel.Modal.Cursor++
		}
		return coreModel, nil
	}
	return coreModel, nil
}

func handleCollectionKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := coreModel.Keys
	switch msg.String() {
	case "up", "k":
		if coreModel.CollectionCursor > 0 {
			coreModel.CollectionCursor--
		}

	case "down", "j":
		if coreModel.CollectionCursor < len(coreModel.CollectionTree)-1 {
			coreModel.CollectionCursor++
		}

	case "enter", " ":
		node := coreModel.CollectionTree[coreModel.CollectionCursor]
		if node.IsFolder {
			// toggle expand/collapse
			coreModel = coreModel.ToggleExpand(coreModel.CollectionCursor)
		} else {
			// load request ke form
			coreModel = coreModel.LoadRequest(node.Data)
			coreModel.CollectionOpen = false // tutup panel, fokus ke form
		}

	case "n":
		// buka modal → new request di dalam collection ini
		coreModel.Modal.Active = config.ModalSaveAs

	case "N":
		// buka modal new folder
		coreModel.Modal.Active = config.ModalNewFolder
		coreModel.Modal.Input = utils.NewModalInput()
		coreModel.Modal.Input.Placeholder = "e.g. Auth"
		coreModel.Modal.Input.Focus()
		return coreModel, nil

	case "d":
		// hapus node yang dipilih
		coreModel = coreModel.DeleteNode(coreModel.CollectionCursor)

	case "r":
		// buka modal rename — isi input dengan nama lama
		if len(coreModel.CollectionTree) == 0 {
			return coreModel, nil
		}
		node := coreModel.CollectionTree[coreModel.CollectionCursor]
		coreModel.Modal.Active = config.ModalRename
		coreModel.Modal.RenameID = node.ID
		coreModel.Modal.Input = utils.NewModalInput()
		coreModel.Modal.Input.SetValue(node.Name) // pre-fill nama lama
		coreModel.Modal.Input.Focus()
		return coreModel, nil
	}

	switch {
	case key.Matches(msg, k.OpenPanelCollection):
		coreModel.CollectionOpen = false
		coreModel.Focused = config.PanelURL
	}

	return coreModel, nil
}

func handleNewFolderKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	fmt.Printf("%q", msg.String())

	switch msg.String() {
	case "esc":
		coreModel.Modal.Active = config.ModalNone
		coreModel.Modal.ErrMsg = ""
		coreModel.Modal.Input.Blur()
		return coreModel, nil

	case "enter":
		name := strings.TrimSpace(coreModel.Modal.Input.Value())
		if name == "" {
			coreModel.Modal.ErrMsg = "folder name cannot be empty"
			return coreModel, nil
		}

		// buat collection baru sebagai root folder
		newFolder := entity.Collection{
			ID:       uuid.New().String(),
			Name:     name,
			Children: []entity.Collection{},
			Requests: []entity.Request{},
			Expanded: false,
		}

		coreModel.Collections = append(coreModel.Collections, newFolder)
		coreModel.CollectionTree = utils.Flatten(coreModel.Collections)
		coreModel.Modal.Active = config.ModalNone
		coreModel.Modal.ErrMsg = ""
		coreModel.Modal.Input.SetValue("")
		coreModel.Modal.Input.Blur()
		return coreModel, SaveCollectionsCmd(newFolder)

	default:
		tea.Println("enter ditekan, name:", msg.String())
		coreModel.Modal.ErrMsg = ""
		var cmd tea.Cmd
		coreModel.Modal.Input, cmd = coreModel.Modal.Input.Update(msg)
		return coreModel, cmd
	}

}

func handleRenameKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		coreModel.Modal.Active = config.ModalNone
		coreModel.Modal.ErrMsg = ""
		coreModel.Modal.RenameID = ""
		coreModel.Modal.Input.Blur()
		return coreModel, nil

	case "enter":
		name := strings.TrimSpace(coreModel.Modal.Input.Value())
		if name == "" {
			coreModel.Modal.ErrMsg = "name cannot be empty"
			return coreModel, nil
		}

		coreModel.Collections = RenameInTree(coreModel.Collections, coreModel.Modal.RenameID, name)
		coreModel.CollectionTree = utils.Flatten(coreModel.Collections)
		coreModel.Modal.Active = config.ModalNone
		root := FindRootCollection(coreModel.Collections, coreModel.Modal.RenameID)

		coreModel.Modal.ErrMsg = ""
		coreModel.Modal.RenameID = ""
		coreModel.Modal.Input.SetValue("")
		coreModel.Modal.Input.Blur()
		return coreModel, SaveCollectionsCmd(root)

	default:
		coreModel.Modal.ErrMsg = ""
		var cmd tea.Cmd
		coreModel.Modal.Input, cmd = coreModel.Modal.Input.Update(msg)
		return coreModel, cmd
	}
}

func handleEnvPageKey(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// kalau sedang edit cell, tangkap semua input
	if coreModel.EnvPageState.Editing {
		return handleEnvCellEdit(coreModel, msg)
	}

	switch coreModel.EnvPageState.Focus {
	case config.EnvFocusList:
		return handleEnvListPanel(coreModel, msg)
	case config.EnvFocusTable:
		return handleEnvTablePanel(coreModel, msg)
	}
	return coreModel, nil
}

func handleEnvListPanel(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		coreModel.ActivePage = config.PageMain
		return coreModel, nil

	case "tab":
		// pindah ke panel kanan
		coreModel.EnvPageState.Focus = config.EnvFocusTable
		return coreModel, nil

	case "up", "k":
		if coreModel.EnvPageState.Cursor > 0 {
			coreModel.EnvPageState.Cursor--
			coreModel.EnvPageState.Rows = utils.EnvToTableRows(coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor])
			coreModel.EnvPageState.Cursor = 0
		}
		return coreModel, nil

	case "down", "j":
		if coreModel.EnvPageState.Cursor < len(coreModel.EnvPageState.List)-1 {
			coreModel.EnvPageState.Cursor++
			coreModel.EnvPageState.Rows = utils.EnvToTableRows(coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor])
			coreModel.EnvPageState.Cursor = 0
		}
		return coreModel, nil

	case "enter":
		// set sebagai env aktif
		selected := coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor]
		coreModel.ActiveEnv = &selected
		coreModel.ActiveEnvIdx = coreModel.EnvPageState.Cursor
		return coreModel, nil

	case "n":
		// tambah env baru
		newEnv := entity.Env{
			Name: "new environment",
			Vars: map[string]string{},
		}
		coreModel.EnvPageState.List = append(coreModel.EnvPageState.List, newEnv)
		coreModel.EnvPageState.Cursor = len(coreModel.EnvPageState.List) - 1
		coreModel.EnvPageState.Rows = []entity.EnvTableRow{utils.NewEnvTableRow()}
		coreModel.EnvPageState.Focus = config.EnvFocusTable // langsung pindah ke tabel
		return coreModel, nil

	case "d":
		// hapus env yang dipilih
		if len(coreModel.EnvPageState.List) == 0 {
			return coreModel, nil
		}
		coreModel.EnvPageState.List = append(
			coreModel.EnvPageState.List[:coreModel.EnvPageState.Cursor],
			coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor+1:]...,
		)
		if coreModel.EnvPageState.Cursor >= len(coreModel.EnvPageState.List) {
			coreModel.EnvPageState.Cursor = len(coreModel.EnvPageState.List) - 1
		}
		if coreModel.EnvPageState.Cursor >= 0 {
			coreModel.EnvPageState.Rows = utils.EnvToTableRows(coreModel.EnvPageState.List[coreModel.EnvPageState.Cursor])
		} else {
			coreModel.EnvPageState.Rows = []entity.EnvTableRow{}
		}
		return coreModel, SaveEnvsCmd(coreModel.EnvPageState.List)
	}
	return coreModel, nil
}

func handleEnvTablePanel(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "tab":
		// kembali ke list
		coreModel.EnvPageState.Focus = config.EnvFocusList
		return coreModel, nil

	case "up", "k":
		if coreModel.EnvPageState.Cursor > 0 {
			coreModel.EnvPageState.Cursor--
		}
		return coreModel, nil

	case "down", "j":
		if coreModel.EnvPageState.Cursor < len(coreModel.EnvPageState.Rows)-1 {
			coreModel.EnvPageState.Cursor++
		}
		return coreModel, nil

	case "enter":
		// mulai edit row yang dipilih
		coreModel.EnvPageState.Editing = true
		coreModel.EnvPageState.Rows[coreModel.EnvPageState.Cursor].Key.Focus()
		return coreModel, nil

	case "n":
		// tambah row baru
		coreModel.EnvPageState.Rows = append(coreModel.EnvPageState.Rows, utils.NewEnvTableRow())
		coreModel.EnvPageState.Cursor = len(coreModel.EnvPageState.Rows) - 1
		coreModel.EnvPageState.Editing = true
		coreModel.EnvPageState.Rows[coreModel.EnvPageState.Cursor].Key.Focus()
		return coreModel, nil

	case "d":
		// hapus row
		if len(coreModel.EnvPageState.Rows) == 0 {
			return coreModel, nil
		}
		coreModel.EnvPageState.Rows = append(
			coreModel.EnvPageState.Rows[:coreModel.EnvPageState.Cursor],
			coreModel.EnvPageState.Rows[coreModel.EnvPageState.Cursor+1:]...,
		)
		if coreModel.EnvPageState.Cursor >= len(coreModel.EnvPageState.Rows) {
			coreModel.EnvPageState.Cursor = len(coreModel.EnvPageState.Rows) - 1
		}
		return coreModel, coreModel.SaveCurrentEnv()
	}
	return coreModel, nil
}

func handleEnvCellEdit(coreModel Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	row := &coreModel.EnvPageState.Rows[coreModel.EnvPageState.Cursor]

	switch msg.String() {
	case "tab":
		// pindah dari key ke value atau selesai edit
		if row.Key.Focused() {
			row.Key.Blur()
			row.Value.Focus()
		} else {
			row.Value.Blur()
			coreModel.EnvPageState.Editing = false
			return coreModel, coreModel.SaveCurrentEnv()
		}
		return coreModel, nil

	case "esc":
		row.Key.Blur()
		row.Value.Blur()
		coreModel.EnvPageState.Editing = false
		return coreModel, coreModel.SaveCurrentEnv()

	case "enter":
		row.Key.Blur()
		row.Value.Blur()
		coreModel.EnvPageState.Editing = false
		return coreModel, coreModel.SaveCurrentEnv()
	}

	// teruskan ke input yang aktif
	var cmd tea.Cmd
	if row.Key.Focused() {
		coreModel.EnvPageState.Rows[coreModel.EnvPageState.Cursor].Key, cmd = row.Key.Update(msg)
	} else {
		coreModel.EnvPageState.Rows[coreModel.EnvPageState.Cursor].Value, cmd = row.Value.Update(msg)
	}
	return coreModel, cmd
}
