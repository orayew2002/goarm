package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/orayew2002/goarm/domain"
	"github.com/orayew2002/goarm/manager"
	"github.com/orayew2002/goarm/utils"
)

//go:embed template/**/*
//go:embed template/.golangci.yml
//go:embed template/Makefile
var templateFS embed.FS

func main() {
	app := utils.OpenForm()

	// Create project files and initialize go.mod
	if err := createProjectFiles(app.Name); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating project files: %v\n", err)
		os.Exit(1)
	}

	if err := bindDependencies(app.Name, app.DbType); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing go.mod: %v\n", err)
		os.Exit(1)
	}

	if err := initializeGoMod(app.Name, app.Name); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing go.mod: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Project setup completed successfully.")
}

// bindDependencies writes embedded database files to disk under the app directory.
// It creates the necessary files: config.yaml, init.go, and domain.go.
func bindDependencies(appName string, database domain.DbType) error {
	mngr := manager.Manage(database.ToCoreDatabase())
	baseDir := filepath.Join(appName, "pkg", database.ToCoreDatabase())

	// Define files to write: filename -> content provider function
	files := map[string]func() []byte{
		"config.yaml": mngr.Database.GetConfig,
		"init.go":     mngr.Database.GetInit,
		"domain.go":   mngr.Database.GetDomain,
	}

	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", baseDir, err)
	}

	// Write all files
	for filename, contentFunc := range files {
		filePath := filepath.Join(baseDir, filename)

		if err := os.WriteFile(filePath, contentFunc(), os.ModePerm); err != nil {
			return fmt.Errorf("failed to write %q: %w", filePath, err)
		}
	}

	return nil
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
