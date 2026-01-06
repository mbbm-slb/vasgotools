// Package main provides VasGoTools - a utility tool for managing Go projects.
// This application provides commands to simplify the creation and management of Go projects,
// including generating Go workspaces, applications, and libraries.
package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	openVSCodeBatchFile = "open_vscode.bat"
	openVSCodeShellFile = "open_vscode.sh"
	modulePrefixMbbVas  = "github.com/muellerbbm-vas/"
	modulePrefixMbbmSlb = "github.com/mbbm-slb/"
)

// Embed the template files

//go:embed build.bat
var buildBatTemplate string

//go:embed build.sh
var buildShTemplate string

//go:embed cross-build.bat
var crossBuildBatTemplate string

//go:embed cross-build.sh
var crossBuildShTemplate string

//go:embed main.go.template
var mainGoTemplate string

//go:embed golangci_win.yml
var golangciWinYmlTemplate string

//go:embed golangci.yml
var golangciYmlTemplate string

//go:embed LICENSE
var licenseTemplate string

func main() {
	// Ensure a subcommand is provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Determine the subcommand
	switch os.Args[1] {
	case "work":
		generateWorkCommand(os.Args[2:])
	case "app":
		generateModuleCommand(os.Args[2:], false)
	case "lib":
		generateModuleCommand(os.Args[2:], true)
	case "help", "--help", "-h":
		printUsage()
		os.Exit(0)
	case "version", "--version", "-v":
		fmt.Println("Version: ", getVersionString())
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func getVersionString() string {
	version := "unknown version"

	// Try to read build info from the binary
	if info, ok := debug.ReadBuildInfo(); ok {
		// Use the module version if available (from go install)
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}

		// Try to get version from VCS info (Git)
		var revision, modified string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value
			}
		}

		// Append VCS info if available
		if revision != "" {
			shortRev := revision
			if len(revision) > 7 {
				shortRev = revision[:7]
			}
			version = fmt.Sprintf("%s (commit: %s)", version, shortRev)
			if modified == "true" {
				version += " [modified]"
			}
		}
	}
	return version
}

func printUsage() {
	fmt.Println("VasGoTools - A utility tool for managing Go projects")
	fmt.Println()
	fmt.Println("Version: ", getVersionString())
	fmt.Println()
	fmt.Println("This application provides commands to simplify the creation and management of Go projects,")
	fmt.Println("including generating Go workspaces, applications, and libraries.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vasgotools.exe <command> [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  work    Generate a Go workspace (i.e., a go.work file)")
	fmt.Println("  app     Create a new Go application")
	fmt.Println("  lib     Create a new Go library")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --path <path>        Specify the folder path (defaults to the current working directory)")
	fmt.Println("  --module-prefix <prefix>")
	fmt.Println("                      Specify the module prefix (default: github.com/muellerbbm-vas/)")
	fmt.Println("  nogit                Skip Git repository initialization")
	fmt.Println("  nocode               Skip creation and execution of the open_vscode file")
	fmt.Println("  nomain               Skip creation of the main.go file (only for app)")
	fmt.Println()
	fmt.Println("Endorsed Folder Structure for Workspaces:")
	fmt.Println("  The recommended folder structure for a Go workspace is as follows:")
	fmt.Println()
	fmt.Println("  <workspace-root>/")
	fmt.Println("  ├── go.work         # The Go workspace file")
	fmt.Println("  ├── app1/           # Application 1 folder")
	fmt.Println("  ├── app2/           # Application 2 folder")
	fmt.Println("  └── ext/            # Folder for libraries")
	fmt.Println("      ├── lib1/       # Library 1 folder")
	fmt.Println("      └── lib2/       # Library 2 folder")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vasgotools.exe work --path \"C:\\projects\\myworkspace\"")
	fmt.Println("  vasgotools.exe app myapp --path \"C:\\projects\"")
	fmt.Println("  vasgotools.exe lib mylib nogit nocode")
	fmt.Println("  vasgotools.exe app myapp nomain nogit")
	fmt.Println("  vasgotools.exe app myapp --module-prefix \"github.com/custom-prefix/\"")
	fmt.Println()
	fmt.Println("For more information, use 'go run main.go <command>' to see command-specific options.")
}

func generateWorkCommand(args []string) {
	// Define a flag set for the "work" command
	fs := flag.NewFlagSet("work", flag.ExitOnError)
	folderPath := fs.String("path", "", "Path to the folder (defaults to current working directory)")
	if err := fs.Parse(args); err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	// Check for optional flags
	noGit, noCode := parseOptionalFlags(fs.Args())

	// Use the current working directory if no path is provided
	err := setDefaultFolderPath(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Slice to store relative paths of subfolders containing go.mod
	var goModFolders []string

	// Walk through the directory
	err = filepath.Walk(*folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current item is a file named "go.mod"
		if info.Name() == "go.mod" {
			// Get the relative path of the folder containing go.mod
			relativePath, err := filepath.Rel(*folderPath, filepath.Dir(path))
			if err != nil {
				return err
			}
			goModFolders = append(goModFolders, relativePath)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
		return
	}

	// Check, if the go.work file already exists. If so, delete it to allow recreation (we do currently not support updating)
	goWorkFilePath := filepath.Join(*folderPath, "go.work")
	if _, err := os.Stat(goWorkFilePath); err == nil {
		fmt.Printf("go.work file already exists at %s => deleting\n", goWorkFilePath)
		err := os.Remove(goWorkFilePath)
		if err != nil {
			fmt.Println("Error deleting go.work file:", err)
			return
		}
	}

	// Check, if the go.work.sum file already exists. If so, delete it to allow recreation (we do currently not support updating)
	goWorkSumFilePath := filepath.Join(*folderPath, "go.work.sum")
	if _, err := os.Stat(goWorkSumFilePath); err == nil {
		fmt.Printf("go.work.sum file already exists at %s => deleting\n", goWorkSumFilePath)
		err := os.Remove(goWorkSumFilePath)
		if err != nil {
			fmt.Println("Error deleting go.work.sum file:", err)
			return
		}
	}

	// Print the collected relative paths
	fmt.Println("Subfolders containing go.mod:")
	for _, folder := range goModFolders {
		fmt.Println(folder)
	}

	// Run the "go work init" command with the relative paths
	if len(goModFolders) > 0 {
		args := append([]string{"work", "init"}, goModFolders...)
		//nolint:gosec // G204: Safe usage - args are controlled by the application
		cmd := exec.Command("go", args...)
		cmd.Dir = *folderPath // Set the working directory to the root folder
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("Running command:", cmd.String())
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running 'go work init':", err)
			return
		}

		fmt.Println("go.work file created successfully.")
	} else {
		fmt.Println("No subfolders with go.mod found. No go.work file created.")
	}

	// Create the open_vscode.bat file (if not suppressed)
	if !noCode {
		err = createOpenVSCodeFile(*folderPath)
		if err != nil {
			fmt.Println("Error creating open_vscode file:", err)
			return
		}

		// Execute the open_vscode file
		err = executeOpenVSCodeFile(*folderPath)
		if err != nil {
			fmt.Println("Error executing open_vscode file:", err)
			return
		}
		fmt.Println("Visual Studio Code opened successfully.")
	} else {
		fmt.Println("Creation and execution of open_vscode.bat skipped.")
	}

	// Initialize a Git repository (if not suppressed)
	if !noGit {
		err = initializeGitRepository(*folderPath)
		if err != nil {
			fmt.Println("Error initializing Git repository:", err)
			return
		}
		fmt.Println("Git repository initialized successfully.")

		// Search for Git repositories in subfolders and add them as submodules
		err = addGitSubmodules(*folderPath)
		if err != nil {
			fmt.Println("Error adding Git submodules:", err)
			return
		}

		// Add all files and commit with message "init"
		cmdAdd := exec.Command("git", "add", ".")
		cmdAdd.Dir = *folderPath
		cmdAdd.Stdout = os.Stdout
		cmdAdd.Stderr = os.Stderr
		if err := cmdAdd.Run(); err != nil {
			fmt.Println("Error adding files to git:", err)
			return
		}

		// create initial commit mit message "init"
		cmdCommit := exec.Command("git", "commit", "-m", "init")
		cmdCommit.Dir = *folderPath
		cmdCommit.Stdout = os.Stdout
		cmdCommit.Stderr = os.Stderr
		if err := cmdCommit.Run(); err != nil {
			fmt.Println("Error committing files to git:", err)
			return
		}
		fmt.Println("All files and submodules added and initial commit created.")
	} else {
		fmt.Println("Git repository initialization skipped.")
	}
}

// addGitSubmodules searches for Git repositories in subfolders and adds them as submodules
func addGitSubmodules(rootPath string) error {
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == rootPath {
			return nil
		}

		// Check if the current folder is a Git repository
		if info.IsDir() && filepath.Base(path) == ".git" {
			submodulePath := filepath.Dir(path)
			relativePath, err := filepath.Rel(rootPath, submodulePath)
			if err != nil {
				return err
			}

			// Skip adding the root directory as a submodule
			if relativePath == "." {
				return nil
			}

			// Add the Git repository as a submodule
			//nolint:gosec // G204: Safe usage - submodulePath and relativePath are controlled by the application
			cmd := exec.Command("git", "submodule", "add", submodulePath, relativePath)
			cmd.Dir = rootPath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Printf("Adding submodule: %s\n", relativePath)
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("error adding submodule %s: %w", relativePath, err)
			}
		}
		return nil
	})

	return err
}

func generateModuleCommand(args []string, isLibrary bool) {
	// Define a flag set for the "app" or "lib" command
	fs := flag.NewFlagSet("app", flag.ExitOnError)
	folderPath := fs.String("path", "", "Path to create the application or library folder (defaults to current working directory)")
	modulePrefixCmd := fs.String("module-prefix", "none", "Specify the module prefix (default: none, shortcuts: 'vas' for muellerbbm-vas, 'slb' for mbbm-slb)")
	if err := fs.Parse(args); err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	// Determine the module prefix
	var modulePrefix string
	if *modulePrefixCmd != "none" {
		switch *modulePrefixCmd {
		case "vas":
			modulePrefix = modulePrefixMbbVas
		case "slb":
			modulePrefix = modulePrefixMbbmSlb
		default:
			modulePrefix = *modulePrefixCmd
		}
	}

	// Ensure the application or library name is provided as the first positional argument
	if fs.NArg() < 1 {
		fmt.Println("Error: Name is required.")
		fmt.Println("Usage: vasgotools.exe app <name> [--path <path>] [--module-prefix <prefix>] [nogit] [nocode] [nomain]")
		os.Exit(1)
	}
	name := fs.Arg(0)

	// Check for optional flags
	noGit, noCode := parseOptionalFlags(fs.Args()[1:])
	noMain := isLibrary // Automatically skip main.go creation for libraries

	// Use the current working directory if no path is provided
	err := setDefaultFolderPath(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create the folder
	folder := filepath.Join(*folderPath, name)
	err = os.MkdirAll(folder, 0o750)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Run the "go mod init" command
	fullName := modulePrefix + name
	//nolint:gosec // G204: Safe usage - module name is controlled by the application
	cmd := exec.Command("go", "mod", "init", fullName)
	cmd.Dir = folder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running command:", cmd.String())
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running 'go mod init':", err)
		return
	}

	// Create analyze scripts and golangci-lint config files
	err = createScripts(folder)
	if err != nil {
		fmt.Println("Error creating analyze scripts:", err)
		return
	}
	fmt.Println("Analyze scripts and configuration files created successfully.")

	// Create LICENSE file
	err = createLicenseFile(folder)
	if err != nil {
		fmt.Println("Error creating LICENSE file:", err)
		return
	}
	fmt.Println("LICENSE file created successfully.")

	// Write main.go from the embedded template (if not suppressed)
	if !noMain {

		// Create the main.go file from the embedded template
		mainGoPath := filepath.Join(folder, "main.go")
		err = os.WriteFile(mainGoPath, []byte(mainGoTemplate), 0o600)
		if err != nil {
			fmt.Println("Error writing main.go:", err)
			return
		}
		fmt.Println("main.go created successfully.")
	} else {
		fmt.Println("Creation of main.go skipped.")
	}

	// Create the open_vscode.bat file (if not suppressed)
	if !noCode {
		err = createOpenVSCodeFile(folder)
		if err != nil {
			fmt.Println("Error creating open_vscode file:", err)
			return
		}

		// Execute the open_vscode file
		err = executeOpenVSCodeFile(folder)
		if err != nil {
			fmt.Println("Error executing open_vscode file:", err)
			return
		}
		fmt.Println("Visual Studio Code opened successfully.")
	} else {
		fmt.Println("Creation and execution of open_vscode.bat skipped.")
	}

	// Initialize a Git repository (if not suppressed)
	if !noGit {
		err = initializeGitRepository(folder)
		if err != nil {
			fmt.Println("Error initializing Git repository:", err)
			return
		}
		fmt.Println("Git repository initialized successfully.")

		// Add all files and commit with message "init"
		cmdAdd := exec.Command("git", "add", ".")
		cmdAdd.Dir = folder
		cmdAdd.Stdout = os.Stdout
		cmdAdd.Stderr = os.Stderr
		if err := cmdAdd.Run(); err != nil {
			fmt.Println("Error adding files to git:", err)
			return
		}

		// create initial commit mit message "init"
		cmdCommit := exec.Command("git", "commit", "-m", "init")
		cmdCommit.Dir = folder
		cmdCommit.Stdout = os.Stdout
		cmdCommit.Stderr = os.Stderr
		if err := cmdCommit.Run(); err != nil {
			fmt.Println("Error committing files to git:", err)
			return
		}
		fmt.Println("All files added and initial commit created.")
	} else {
		fmt.Println("Git repository initialization skipped.")
	}

	fmt.Printf("'%s' created successfully in folder '%s'.\n", fullName, folder)
}

// parseOptionalFlags parses the optional "nogit" and "nocode" flags from the arguments.
func parseOptionalFlags(args []string) (noGit, noCode bool) {
	noGit = false
	noCode = false
	for _, arg := range args {
		switch arg {
		case "nogit":
			noGit = true
		case "nocode":
			noCode = true
		}
	}
	return noGit, noCode
}

// setDefaultFolderPath sets the folder path to the current working directory if it is empty.
func setDefaultFolderPath(folderPath *string) error {
	if *folderPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current working directory: %w", err)
		}
		*folderPath = cwd
	}
	return nil
}

func initializeGitRepository(folderPath string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = folderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Initializing Git repository...")
	return cmd.Run()
}

func createOpenVSCodeBatchFile(folderPath string) error {
	batchFilePath := filepath.Join(folderPath, openVSCodeBatchFile)
	batchFileContent := "code . | exit 0\n"
	return os.WriteFile(batchFilePath, []byte(batchFileContent), 0o600)
}

func executeOpenVSCodeBatchFile(folderPath string) error {
	batchFilePath := filepath.Join(folderPath, openVSCodeBatchFile)
	//nolint:gosec // G204: Safe usage - batchFilePath is controlled by the application
	cmd := exec.Command("cmd", "/C", batchFilePath)
	cmd.Dir = folderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Opening Visual Studio Code...")
	return cmd.Run()
}

func createOpenVSCodeShellScript(folderPath string) error {
	scriptFilePath := filepath.Join(folderPath, openVSCodeShellFile)
	scriptContent := "#!/bin/bash\ncode . || exit 0\n"
	//nolint:gosec // G306: Script needs to be executable
	err := os.WriteFile(scriptFilePath, []byte(scriptContent), 0o700) // Make the script executable
	if err != nil {
		return fmt.Errorf("error creating open_vscode.sh: %w", err)
	}
	return nil
}

func executeOpenVSCodeShellScript(folderPath string) error {
	scriptFilePath := filepath.Join(folderPath, openVSCodeShellFile)
	//nolint:gosec // G204: Safe usage - scriptFilePath is controlled by the application
	cmd := exec.Command("bash", scriptFilePath)
	cmd.Dir = folderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Opening Visual Studio Code...")
	return cmd.Run()
}

func createOpenVSCodeFile(folderPath string) error {
	err1 := createOpenVSCodeBatchFile(folderPath)
	err2 := createOpenVSCodeShellScript(folderPath)
	return errors.Join(err1, err2)
}

func executeOpenVSCodeFile(folderPath string) error {
	if runtime.GOOS == "windows" {
		return executeOpenVSCodeBatchFile(folderPath)
	}
	return executeOpenVSCodeShellScript(folderPath)
}

func createBuildBatchFile(folderPath string) error {
	batchFilePath := filepath.Join(folderPath, "build.bat")
	err := os.WriteFile(batchFilePath, []byte(buildBatTemplate), 0o600)
	if err != nil {
		return fmt.Errorf("error creating build.bat: %w", err)
	}
	return nil
}

func createBuildShellScript(folderPath string) error {
	scriptFilePath := filepath.Join(folderPath, "build.sh")
	//nolint:gosec // G306: Script needs to be executable
	err := os.WriteFile(scriptFilePath, []byte(buildShTemplate), 0o700) // Make the script executable
	if err != nil {
		return fmt.Errorf("error creating build.sh: %w", err)
	}
	return nil
}

func createBuildScript(folderPath string) error {
	err1 := createBuildBatchFile(folderPath)
	err2 := createBuildShellScript(folderPath)
	return errors.Join(err1, err2)
}

func createCrossBuildBatchFile(folderPath string) error {
	batchFilePath := filepath.Join(folderPath, "cross-build.bat")
	err := os.WriteFile(batchFilePath, []byte(crossBuildBatTemplate), 0o600)
	if err != nil {
		return fmt.Errorf("error creating cross-build.bat: %w", err)
	}
	return nil
}

func createCrossBuildShellScript(folderPath string) error {
	scriptFilePath := filepath.Join(folderPath, "cross-build.sh")
	//nolint:gosec // G306: Script needs to be executable
	err := os.WriteFile(scriptFilePath, []byte(crossBuildShTemplate), 0o700) // Make the script executable
	if err != nil {
		return fmt.Errorf("error creating cross-build.sh: %w", err)
	}
	return nil
}

func createCrossBuildScript(folderPath string) error {
	err1 := createCrossBuildBatchFile(folderPath)
	err2 := createCrossBuildShellScript(folderPath)
	return errors.Join(err1, err2)
}

func createGolangciWinYml(folderPath string) error {
	ymlFilePath := filepath.Join(folderPath, "golangci_win.yml")
	err := os.WriteFile(ymlFilePath, []byte(golangciWinYmlTemplate), 0o600)
	if err != nil {
		return fmt.Errorf("error creating golangci_win.yml: %w", err)
	}
	return nil
}

func createGolangciYml(folderPath string) error {
	ymlFilePath := filepath.Join(folderPath, "golangci.yml")
	err := os.WriteFile(ymlFilePath, []byte(golangciYmlTemplate), 0o600)
	if err != nil {
		return fmt.Errorf("error creating golangci.yml: %w", err)
	}
	return nil
}

func createScripts(folderPath string) error {
	err1 := createBuildScript(folderPath)
	err2 := createCrossBuildScript(folderPath)
	err3 := createGolangciWinYml(folderPath)
	err4 := createGolangciYml(folderPath)
	return errors.Join(err1, err2, err3, err4)
}

func createLicenseFile(folderPath string) error {
	licenseFilePath := filepath.Join(folderPath, "LICENSE")
	// Replace the year placeholder with the current year
	currentYear := time.Now().Year()
	licenseContent := strings.ReplaceAll(licenseTemplate, "2026", fmt.Sprintf("%d", currentYear))
	err := os.WriteFile(licenseFilePath, []byte(licenseContent), 0o600)
	if err != nil {
		return fmt.Errorf("error creating LICENSE: %w", err)
	}
	return nil
}
