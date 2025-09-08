package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/MH-KodaCore/goarm/domain"
	"github.com/MH-KodaCore/goarm/manager"
	"github.com/MH-KodaCore/goarm/utils"
)

//go:embed templates/**/*
var templatesFS embed.FS

func main() {
	app := utils.OpenForm()

	if err := createProjectFiles(app.Name, app.Framework); err != nil {
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

	if err := utils.UpdatePackageNameOnGCI(app.Name); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating package name for .golangci.yml file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Project setup completed successfully.")
}

func bindDependencies(appName string, dbType domain.DbType) error {
	coreDB := dbType.ToCoreDatabase()
	dbValue := dbType.PackageVal()
	manager := manager.Manage(coreDB)

	baseDir := path.Join(appName, "pkg", coreDB)

	// ───── Step 1: Create base directory ─────
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", baseDir, err)
	}

	// ───── Step 2: Write Go source files ─────
	files := map[string][]byte{
		"init.go": manager.Database.GetInit(),
	}
	for name, content := range files {
		path := path.Join(baseDir, name)
		if err := os.WriteFile(path, content, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file %q: %w", path, err)
		}
	}

	// ───── Step 3: Append config to env files ─────
	envFiles := []string{"dev.yaml", "local.yaml", "prod.yaml"}
	for _, env := range envFiles {
		configPath := path.Join(appName, "etc", env)
		if err := utils.AppendToFile(configPath, manager.Database.GetConfig()); err != nil {
			return fmt.Errorf("failed to append config to %q: %w", configPath, err)
		}
	}

	// ───── Step 4: Add field to AppConfig struct ─────
	appStructPath := path.Join(appName, "internal", "domain", "app.go")
	appField := fmt.Sprintf(`DB %s.Config `+"`mapstructure:\"%s\" yaml:\"%s\"`", coreDB, coreDB, coreDB)
	if err := utils.AppendFieldStruct(appStructPath, "AppConfigs", appField); err != nil {
		return fmt.Errorf("failed to append field to AppConfigs struct: %w", err)
	}

	// ───── Step 5: Add import for DB package ─────
	if err := utils.AddImportToFile(appStructPath, path.Join(appName, "pkg", coreDB)); err != nil {
		return fmt.Errorf("failed to add import to app.go: %w", err)
	}

	// ───── Step 6: Update Repo struct ─────
	repoFile := path.Join(appName, "internal", "repo", "build.go")
	if err := utils.AddImportToFile(repoFile, dbType.PackagePath()); err != nil {
		return fmt.Errorf("failed to add DB import to repo/build.go: %w", err)
	}

	repoField := fmt.Sprintf("db *%s", dbValue)
	if err := utils.AppendFieldStruct(repoFile, "Repo", repoField); err != nil {
		return fmt.Errorf("failed to append field to Repo struct: %w", err)
	}

	// ───── Step 7: Update NewRepo constructor ─────
	if err := utils.AppendFuncArgument(repoFile, "NewRepo", "db", "*"+dbValue); err != nil {
		return fmt.Errorf("failed to append argument to NewRepo function: %w", err)
	}
	if err := utils.AddReturnFieldToConstructor(repoFile, "NewRepo", "db"); err != nil {
		return fmt.Errorf("failed to set constructor return value: %w", err)
	}

	// ───── Step 8: Update app run layer ─────
	appRunFile := path.Join(appName, "internal", "app", "build.go")
	if err := utils.AddImportToFile(appRunFile, path.Join(appName, "pkg", coreDB)); err != nil {
		return fmt.Errorf("failed to add DB import to app/build.go: %w", err)
	}

	callArg := fmt.Sprintf("%s.NewClient(appConfig.DB)", coreDB)
	if err := utils.AddArgumentToFunctionCall(appRunFile, "repo.NewRepo", callArg); err != nil {
		return fmt.Errorf("failed to inject database client into repo.NewRepo: %w", err)
	}

	dockerComposeFile := fmt.Sprintf("%s/docker-compose.yaml", appName)
	fileBody, err := os.ReadFile(dockerComposeFile)
	if err != nil {
		return fmt.Errorf("failed to open docker-compose file: %w", err)
	}

	// Replace the "@db" placeholder with the actual DB config
	fileContent := strings.Replace(string(fileBody), "@db", dbType.GetDockerConfig(), 1)
	fileContent = strings.Replace(fileContent, "@dn", dbType.GetDockerDependence(), 1)

	if err := os.WriteFile(dockerComposeFile, []byte(fileContent), 0o644); err != nil {
		return fmt.Errorf("failed to write to docker-compose file: %w", err)
	}

	return nil
}

// createProjectFiles copies template files to the new project directory,
// replacing the string "template" with the project name in the file contents and paths.
func createProjectFiles(projectName string, framework domain.FrameworkType) error {
	templatesDir := "templates/" + framework.ToDirectory()
	err := fs.WalkDir(templatesFS, templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file content
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace "template" with the project name inside file content
		updatedContent := strings.ReplaceAll(string(content), "templates", projectName)

		// Calculate relative path and target destination path
		relativePath, err := filepath.Rel(templatesDir, path)
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
