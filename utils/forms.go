package utils

import (
	"fmt"
	"os"

	"github.com/MH-KodaCore/goarm/domain"
	tea "github.com/charmbracelet/bubbletea"
)

func OpenForm() domain.App {
	// Clear before project name input
	clearScreen()
	projectForm := newProjectNameForm()

	if _, err := tea.NewProgram(projectForm).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run project name form: %v\n", err)
		os.Exit(1)
	}

	projectName := projectForm.GetAppName()
	if len(projectName) == 0 {
		fmt.Fprintln(os.Stderr, "Project name cannot be empty.")
		os.Exit(1)
	}

	// Clear before framework selection
	clearScreen()
	frameworkForm := newFrameworkSelectForm()

	if _, err := tea.NewProgram(&frameworkForm).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run framework select form: %v\n", err)
		os.Exit(1)
	}

	if len(frameworkForm.GetChoice()) == 0 {
		fmt.Fprintln(os.Stderr, "Project framework cannot be empty.")
		os.Exit(1)
	}

	// Clear before database selection
	clearScreen()
	databaseForm := newDatabaseSelectForm()

	if _, err := tea.NewProgram(&databaseForm).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run database select form: %v\n", err)
		os.Exit(1)
	}

	if len(databaseForm.GetChoice()) == 0 {
		fmt.Fprintln(os.Stderr, "Project database cannot be empty.")
		os.Exit(1)
	}

	// Return app struct (capture DB choice here if needed)
	clearScreen()
	return domain.App{
		Name:      projectForm.GetAppName(),
		Framework: domain.FrameworkType(frameworkForm.GetChoice()),
		DbType:    domain.DbType(databaseForm.GetChoice()),
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	showPreload()
}
