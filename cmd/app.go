package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mboufous/gkns/app"
	"github.com/mboufous/gkns/app/tui"
	"os"
)

// Execute starts app using bubble tea
func Execute() int {
	config := app.NewConfig()
	client, err := app.NewClient(config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cannot create client: %v\n", err)
		return 1
	}

	if ok, err := client.CheckServerConnection(); !ok {
		_, _ = fmt.Fprintf(os.Stderr, "cannot connect to k8s server: %v\n", err)
		return 1
	}

	model, err := tui.New(client)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Run: %v\n", err)
		return 1
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return 1
	}
	return 0
}
