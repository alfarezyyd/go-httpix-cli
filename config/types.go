package config

// Panel identifies which UI region currently has keyboard focus.
type Panel int
type ModalType int
type Page int
type EnvPageFocus int

const (
	PageMain Page = iota
	PageEnv       // ← full page env manager
)

const (
	EnvFocusList EnvPageFocus = iota
	EnvFocusTable
)
const (
	PanelURL Panel = iota
	PanelBody
	PanelResponse
	PanelHistory
	PanelCollections
)

// BodyTab is the active sub-tab inside the request editor.
type BodyTab int

const (
	TabBody BodyTab = iota
	TabHeaders
	TabParams
)

// HTTPMethods is the ordered list of selectable HTTP verbs.
var HTTPMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

const (
	ModalNone ModalType = iota
	ModalSaveAs
	ModalEnvPicker
	ModalNewFolder
	ModalRename
)

var PanelOrder = []Panel{
	PanelURL,
	PanelBody,
	PanelResponse,
	PanelHistory,
}
