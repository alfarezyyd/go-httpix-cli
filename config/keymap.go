package config

import (
	"runtime"

	"github.com/charmbracelet/bubbles/key"
)

// KeyMap holds every keybinding plus display labels for the status bar.
// Labels are already formatted for the current platform (⌃/⌥ vs ctrl).
type KeyMap struct {
	Send         key.Binding
	NextPanel    key.Binding
	PrevPanel    key.Binding
	NextMethod   key.Binding
	PrevMethod   key.Binding
	NextTab      key.Binding
	PrevTab      key.Binding
	HistoryUp    key.Binding
	HistoryDown  key.Binding
	Quit         key.Binding
	FormatJSON   key.Binding
	ClearHistory key.Binding

	// Human-readable labels shown in the status bar / top bar.
	LabelSend         string
	LabelFocus        string
	LabelMethod       string
	LabelTab          string
	LabelFormatJSON   string
	LabelQuit         string
	LabelClearHistory string
}

// DefaultKeyMap returns a KeyMap tuned for the current OS.
//
// macOS terminal caveats:
//   - ⌘ Cmd is captured by the OS — not available in terminal apps.
//   - ⌥ Option reaches the app as "alt+<key>" (ESC-prefixed byte sequence).
//   - ctrl+[ ≡ ESC — must be avoided on macOS.
//   - ctrl+enter is unreliable in iTerm2 / Terminal.app → use ctrl+s.
//   - ctrl+left/right conflict with Mission Control → use alt+arrow.
//   - ctrl+h ≡ Backspace in many terminals → avoid.
//   - ctrl+d sends EOF → avoid.
func DefaultKeyMap() KeyMap {
	if runtime.GOOS == "darwin" {
		return macKeyMap()
	}
	return linuxKeyMap()
}

func macKeyMap() KeyMap {
	return KeyMap{
		Send:         key.NewBinding(key.WithKeys("ctrl+s", "f5"), key.WithHelp("⌃S/F5", "send")),
		NextPanel:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("⇥", "next panel")),
		PrevPanel:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("⇧⇥", "prev panel")),
		NextMethod:   key.NewBinding(key.WithKeys("alt+f"), key.WithHelp("⌥→", "method +")),
		PrevMethod:   key.NewBinding(key.WithKeys("alt+b"), key.WithHelp("⌥←", "method -")),
		NextTab:      key.NewBinding(key.WithKeys("alt+]"), key.WithHelp("⌥]", "tab +")),
		PrevTab:      key.NewBinding(key.WithKeys("alt+["), key.WithHelp("⌥[", "tab -")),
		HistoryUp:    key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "history up")),
		HistoryDown:  key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "history down")),
		Quit:         key.NewBinding(key.WithKeys("ctrl+c", "ctrl+q"), key.WithHelp("⌃C", "quit")),
		FormatJSON:   key.NewBinding(key.WithKeys("ctrl+f"), key.WithHelp("⌃F", "format JSON")),
		ClearHistory: key.NewBinding(key.WithKeys("ctrl+k"), key.WithHelp("⌃K", "clear history")),

		LabelSend: "⌃S / F5", LabelFocus: "⇥ Tab", LabelMethod: "⌥ ←/→",
		LabelTab: "⌥ [/]", LabelFormatJSON: "⌃F", LabelQuit: "⌃C", LabelClearHistory: "⌃K",
	}
}

func linuxKeyMap() KeyMap {
	return KeyMap{
		Send:         key.NewBinding(key.WithKeys("ctrl+enter", "f5"), key.WithHelp("ctrl+↵/F5", "send")),
		NextPanel:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next panel")),
		PrevPanel:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("⇧tab", "prev panel")),
		NextMethod:   key.NewBinding(key.WithKeys("ctrl+right", "ctrl+l"), key.WithHelp("ctrl+→", "method +")),
		PrevMethod:   key.NewBinding(key.WithKeys("ctrl+left"), key.WithHelp("ctrl+←", "method -")),
		NextTab:      key.NewBinding(key.WithKeys("ctrl+]"), key.WithHelp("ctrl+]", "tab +")),
		PrevTab:      key.NewBinding(key.WithKeys("ctrl+["), key.WithHelp("ctrl+[", "tab -")),
		HistoryUp:    key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "history up")),
		HistoryDown:  key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "history down")),
		Quit:         key.NewBinding(key.WithKeys("ctrl+c", "ctrl+q"), key.WithHelp("ctrl+c", "quit")),
		FormatJSON:   key.NewBinding(key.WithKeys("ctrl+f"), key.WithHelp("ctrl+f", "format JSON")),
		ClearHistory: key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "clear history")),

		LabelSend: "ctrl+↵ / F5", LabelFocus: "Tab", LabelMethod: "ctrl+←/→",
		LabelTab: "ctrl+[/]", LabelFormatJSON: "ctrl+f", LabelQuit: "ctrl+c", LabelClearHistory: "ctrl+d",
	}
}
