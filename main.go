package main

import (
	"fmt"
	"go-httpix-cli/core"
	"go-httpix-cli/outbound"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/imroc/req/v3"
)

func main() {
	req.DevMode()

	tuiProgram := tea.NewProgram(
		core.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	outbound.Setup()
	if _, err := tuiProgram.Run(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, "Error:", err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
