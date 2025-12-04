package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/samuelenocsson/devops-tui/internal/api"
	"github.com/samuelenocsson/devops-tui/internal/config"
	"github.com/samuelenocsson/devops-tui/internal/ui"
)

// Execute runs the application
func Execute() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// If config not found, try to create default
		if err := config.CreateDefaultConfig(); err == nil {
			fmt.Println("Created default config file at ~/.config/devops-tui/config.yaml")
			fmt.Println("Please edit the config file with your Azure DevOps settings.")
			os.Exit(0)
		}
		return fmt.Errorf("configuration error: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Create and run the TUI
	app := ui.NewApp(client)

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running application: %w", err)
	}

	return nil
}
