package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/orayew2002/goarm/utils"
)

//go:embed template/**/*
//go:embed template/.golangci.yml
//go:embed template/Makefile
var templateFS embed.FS

func main() {
	// Initialize the project name form
	projectForm := utils.NewProjectNameForm()
	p := tea.NewProgram(projectForm)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run program: %v\n", err)
		os.Exit(1)
	}

	projectName := projectForm.GetAppName()
	if len(projectName) == 0 {
		fmt.Fprintf(os.Stderr, "Project can't be empty")
		os.Exit(1)
	}

	// Create project files and initialize go.mod
	if err := createProjectFiles(projectName); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating project files: %v\n", err)
		os.Exit(1)
	}

	if err := initializeGoMod(projectName, projectName); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing go.mod: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Project setup completed successfully.")
}

// createProjectFiles copies template files to the new project directory,
// replacing the string "template" with the project name in the file contents and paths.
func createProjectFiles(projectName string) error {
	err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file content
		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace "template" with the project name inside file content
		updatedContent := strings.ReplaceAll(string(content), "template", projectName)

		// Calculate relative path and target destination path
		relativePath, err := filepath.Rel("template", path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(projectName, relativePath)

		// Ensure the target directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		// Write the updated content to the target file
		return os.WriteFile(targetPath, []byte(updatedContent), os.ModePerm)
	})

	if err != nil {
		return err
	}

	fmt.Println("Project files created successfully.")
	return nil
}

// initializeGoMod runs `go mod init` and `go mod tidy` in the project directory.
func initializeGoMod(projectDir, moduleName string) error {
	projectDir = filepath.Clean(projectDir)

	if err := runCommand(projectDir, "go", "mod", "init", moduleName); err != nil {
		return fmt.Errorf("failed to run 'go mod init': %w", err)
	}

	if err := runCommand(projectDir, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %w", err)
	}

	fmt.Println("✅ go.mod initialized and tidied.")
	return nil
}

// runCommand executes a command in a specific directory and prints output.
func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	fmt.Printf("▶️ Running: %s %v\n%s\n", name, args, output)
	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
