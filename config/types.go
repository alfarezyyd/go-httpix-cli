package config

// Panel identifies which UI region currently has keyboard focus.
type Panel int

const (
	PanelURL Panel = iota
	PanelBody
	PanelResponse
	PanelHistory
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
