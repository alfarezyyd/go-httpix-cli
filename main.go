package main

import (
	"fmt"
	"go-httpix-cli/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/imroc/req/v3"
)

func main() {
	req.DevMode()

	p := tea.NewProgram(
		tui.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
